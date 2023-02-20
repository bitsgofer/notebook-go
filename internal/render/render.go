package render

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitsgofer/notebook-go/internal/fileutil"
	klog "k8s.io/klog/v2"
)

func RenderSinglePage(dstPath, srcPath string, themeFS fs.FS, templatePath string, luaFilterPaths ...string) error {
	// copy theme files to tmp dir.
	// NOTE: not too efficient, since we copy once per function call.
	// => try to call render() once per template/type of content,
	// rather than once per content file.
	// -------------------------------------------------------------------------
	// copy template
	templateTmpDir, tmpTemplatePaths, err := copyToTempDir(themeFS, templatePath)
	if err != nil {
		return fmt.Errorf("cannot copy template file to temp dir; err= %w", err)
	}
	defer os.RemoveAll(templateTmpDir)
	tmpTemplatePath := tmpTemplatePaths[0]
	// -------------------------------------------------------------------------
	// copy lua filters
	luaTmpDir, tmpLuaFilterPaths, err := copyToTempDir(themeFS, luaFilterPaths...)
	if err != nil {
		return fmt.Errorf("cannot copy lua filters to temp dir; err= %w", err)
	}
	defer os.RemoveAll(luaTmpDir)
	var luaFilters []string
	for _, f := range tmpLuaFilterPaths {
		luaFilters = append(luaFilters, "--lua-filter", f)
	}
	// -------------------------------------------------------------------------

	// get metadata
	const timeFormat = "2006-January-02"
	stat, err := os.Lstat(srcPath)
	if err != nil {
		return fmt.Errorf("cannot get stat of %s; err= %w", srcPath, err)
	}
	lastModifiedDate := stat.ModTime().Format(timeFormat)
	currentDate := time.Now().Format(timeFormat)
	dateMetadata := []string{
		"--metadata", fmt.Sprintf("current-date=%s", currentDate),
		"--metadata", fmt.Sprintf("last-modified-date=%s", lastModifiedDate),
	}

	// run pandoc command
	// -------------------------------------------------------------------------
	// args
	args := []string{
		"--standalone",
		"--template", tmpTemplatePath,
		"--no-highlight",
	}
	args = append(args, dateMetadata...)
	args = append(args, luaFilters...)
	args = append(args, fmt.Sprintf("--output=%s", dstPath))
	args = append(args, srcPath)
	klog.V(4).InfoS("pandoc", "args", args)
	// -------------------------------------------------------------------------
	cmd := exec.Command("pandoc", args...)
	// redirect command stdout -> dstPath
	html, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("cannot create html file; err= %w", err)
	}
	cmd.Stdout = html
	// access command stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("cannot get stderr of pandoc command; err= %w", err)
	}
	// run command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cannot start pandoc command; err= %w", err)
	}
	// read stderr
	pandocStderr, _ := io.ReadAll(stderr)
	// wait til command finish, then close stdout, stderr
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("cannot wait for pandoc cmd to finish; stderr= %s; err= %w", pandocStderr, err)
	}

	klog.V(3).InfoS("rendered with pandoc", "src", srcPath, "dst", dstPath, "template", templatePath, "lua-filters", luaFilterPaths)
	return nil
}

func RenderMultiplePages(dstDir, srcDir string, themeFS fs.FS, templatePath string, luaFilterPaths ...string) error {
	srcFS := os.DirFS(srcDir)
	renderErr := fs.WalkDir(srcFS, ".", func(filePath string, d fs.DirEntry, pathErr error) error {

		// exit if we see a path error
		if pathErr != nil {
			return pathErr
		}
		// skip non-README.md files
		if isContent := strings.HasSuffix(filePath, "README.md"); !isContent {
			return nil
		}

		dirName := path.Base(path.Dir(filePath))
		dst := filepath.Join(dstDir, dirName, "index.html")
		src := filepath.Join(srcDir, filePath)
		// ensure dstDir
		if err := fileutil.EnsureDir(path.Dir(dst)); err != nil {
			return fmt.Errorf("cannot create directory for rendered file; err= %w", err)
		}
		// render
		if err := RenderSinglePage(
			dst,
			src,
			themeFS,
			templatePath,
			luaFilterPaths...,
		); err != nil {
			return fmt.Errorf("cannot render src= %s to dst= %s; err= %w", src, dst, err)
		}

		return nil
	})
	if renderErr != nil {
		return fmt.Errorf("cannot render all content; err= %w", renderErr)
	}

	klog.V(3).InfoS("rendered with pandoc", "srcDir", srcDir, "dstDir", dstDir, "template", templatePath, "lua-filters", luaFilterPaths)
	return nil
}
