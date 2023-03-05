package astra

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	klog "k8s.io/klog/v2"

	"github.com/bitsgofer/notebook-go/internal/fileutil"
	"github.com/bitsgofer/notebook-go/internal/render"
)

// themeFS are files included with the theme.
//
//go:embed *
var themeFS embed.FS

// bitsgoferTheme is a wrapper struct for the theme.
// It should not have any data.
type bitsgoferTheme struct {
}

func New() *bitsgoferTheme {
	return &bitsgoferTheme{}
}

// CompileAssets implements the Theme interface.
func (bth *bitsgoferTheme) CompileAssets(outputDir string) error {
	// compile CSS and JS from components
	var cssFiles []string
	var jsFiles []string
	// -------------------------------------------------------------------------
	// mini.css
	cssFiles = append(cssFiles,
		"assets/minicss/mini-default.css",
	)
	// prism.js
	cssFiles = append(cssFiles,
		"assets/prism/prism.css",
	)
	jsFiles = append(jsFiles,
		"assets/prism/prism.js",
	)
	// astra.css
	cssFiles = append(cssFiles,
		"assets/astra/fonts.css",
	)

	// minify CSS
	dstCSS := filepath.Join(outputDir, "css", "astra.css")
	if err := render.Minify(dstCSS, themeFS, cssFiles...); err != nil {
		return fmt.Errorf("cannot minify CSS; err= %w", err)
	}
	klog.V(2).InfoS("minified CSS files", "dst", dstCSS)
	// minify JS
	dstJS := filepath.Join(outputDir, "js", "astra.js")
	if err := render.Minify(dstJS, themeFS, jsFiles...); err != nil {
		return fmt.Errorf("cannot minify JS; err= %w", err)
	}
	klog.V(2).InfoS("minified JS files", "dst", dstJS)
	// -------------------------------------------------------------------------

	if err := copyFonts(outputDir); err != nil {
		klog.ErrorS(err, "cannot copy fonts")
		return fmt.Errorf("cannot copy fonts; err= %w", err)
	}

	if err := copyImages(outputDir); err != nil {
		klog.ErrorS(err, "cannot copy logo")
		return fmt.Errorf("cannot copy logo; err= %w", err)
	}

	// copy prism.js's languages
	if err := copyPrismJSLanguages(outputDir); err != nil {
		return fmt.Errorf("cannot copy prism.js's languages; err= %w", err)
	}

	// copy favicon
	srcFavicon := "assets/favicon/favicon.ico"
	dstFavicon := filepath.Join(outputDir, "favicon.ico")
	if err := fileutil.CopyFromFS(dstFavicon, themeFS, srcFavicon); err != nil {
		klog.ErrorS(err, "cannot copy favicon", "src", srcFavicon, "dst", dstFavicon)
		return fmt.Errorf("cannot copy favicon (src= %s) from theme (dst= %s); err= %w", srcFavicon, dstFavicon, err)
	}
	klog.V(2).InfoS("copied favicon", "dst", dstFavicon)

	return nil
}

// copyFonts copies all the fonts file to a folder in the output directory.
func copyFonts(outputDir string) error {
	outputFontsDir := filepath.Join(outputDir, "fonts")
	if err := fileutil.EnsureDir(outputFontsDir); err != nil {
		return fmt.Errorf("cannot create fonts dir; err= %w", err)
	}

	fontFiles := []string{
		"Handlee-Regular.ttf",
	}

	for _, font := range fontFiles {
		srcFont := fmt.Sprintf("assets/fonts/%s", font)
		dstFont := filepath.Join(outputDir, "fonts", font)
		if err := fileutil.CopyFromFS(dstFont, themeFS, srcFont); err != nil {
			klog.ErrorS(err, "cannot copy font", "src", srcFont, "dst", dstFont)
			return fmt.Errorf("cannot copy font (src= %s) from theme (dst= %s); err= %w", srcFont, dstFont, err)
		}
		klog.V(3).InfoS("copied font", "font", font, "dst", dstFont)
	}
	klog.V(2).InfoS("copied fonts", "dst", outputFontsDir)

	return nil
}

// copyImages copies all image files to a folder in the output directory.
func copyImages(outputDir string) error {
	outputImgagesDir := filepath.Join(outputDir, "images")
	if err := fileutil.EnsureDir(outputImgagesDir); err != nil {
		return fmt.Errorf("cannot create images dir; err= %w", err)
	}

	imageFiles := []string{
		"logo.png",
	}

	for _, img := range imageFiles {
		srcImg := fmt.Sprintf("assets/images/%s", img)
		dstImg := filepath.Join(outputDir, "images", img)
		if err := fileutil.CopyFromFS(dstImg, themeFS, srcImg); err != nil {
			klog.ErrorS(err, "cannot copy image", "src", srcImg, "dst", dstImg)
			return fmt.Errorf("cannot copy image (src= %s) from theme (dst= %s); err= %w", srcImg, dstImg, err)
		}
		klog.V(3).InfoS("copied image", "image", img, "dst", dstImg)
	}
	klog.V(2).InfoS("copied images", "dst", outputImgagesDir)

	return nil
}

// copyPrismJSLanguages copies the languages for prism.js's Autoloader
// into a specific path in the output directory.
// The Autoloader plugin will load files in this path
// when it needs to highlight languages new didn' bundle in the core JS.
func copyPrismJSLanguages(outputDir string) error {
	// create output dir for prism.js languages
	outputLanguageDir := filepath.Join(outputDir, "js", "prism", "languages")
	if err := fileutil.EnsureDir(outputLanguageDir); err != nil {
		return fmt.Errorf("cannot create prism's Autoloader's languages dir; err= %w", err)
	}

	// copy [LANGUAGE].min.js from the theme
	languagesDir := "assets/prism/languages"
	walkErr := fs.WalkDir(themeFS, languagesDir, func(path string, d fs.DirEntry, pathErr error) error {
		fileName := d.Name()
		klog.V(4).InfoS("processing prism.js language file", "path", path, "fileName", fileName)

		// exit if we see a path error
		if pathErr != nil {
			return pathErr
		}
		// skip non *.min.js files
		if isMinJS := strings.HasSuffix(fileName, ".min.js"); !isMinJS {
			return nil
		}

		// copy to output
		src := path
		dst := filepath.Join(outputLanguageDir, fileName)
		if err := fileutil.CopyFromFS(dst, themeFS, src); err != nil {
			klog.ErrorS(err, "cannot copy language (.min.js) file from theme to output dir", "src", src, "dst", dst)
			return fmt.Errorf("cannot copy language file (src= %s) from theme (dst= %s); err= %w", src, dst, err)
		}
		klog.V(3).InfoS("copied prism.js language file", "src", path)

		return nil
	})
	if walkErr != nil {
		return fmt.Errorf("cannot process prism.js's languages (.min.js) file; err= %w", walkErr)
	}
	klog.V(2).InfoS("copied prism.js languages file", "dst", outputLanguageDir)

	return nil
}
