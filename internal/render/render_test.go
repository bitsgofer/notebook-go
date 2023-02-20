package render

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderSinglePage(t *testing.T) {
	testCases := map[string]struct {
		themeFS        fs.FS
		srcPath        string
		templatePath   string
		luaFilterPaths []string
		// output
		isError          bool
		expectedRendered []byte
	}{
		"about-me": { // simple markdown
			themeFS:          testThemeFS,
			srcPath:          "testdata/content/about-me/README.md",
			templatePath:     "simple.template.html",
			luaFilterPaths:   []string{"pandoc/lua-filters/standard-code.lua"},
			expectedRendered: renderedAboutMeHTML,
		},
		"hello": { // more complex markdown
			themeFS:          testThemeFS,
			srcPath:          "testdata/content/about-me/README.md",
			templatePath:     "simple.template.html",
			luaFilterPaths:   []string{"pandoc/lua-filters/standard-code.lua"},
			expectedRendered: renderedAboutMeHTML,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			renderedFile := mustCreateTempFile(t)

			err := RenderSinglePage(renderedFile, tc.srcPath, tc.themeFS, tc.templatePath, tc.luaFilterPaths...)
			if tc.isError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// check rendered content
			dstB, err := ioutil.ReadFile(renderedFile)
			assert.NoErrorf(t, err, "cannot read (tmp) dst file")
			assert.Equalf(t, dstB, tc.expectedRendered, "rendered content is wrong")
		})
	}
}

func TestRenderMultiplePages(t *testing.T) {
	testCases := map[string]struct {
		themeFS        fs.FS
		srcDir         string
		templatePath   string
		luaFilterPaths []string
		// output
		isError          bool
		expectedRendered fs.FS
	}{
		"simple": {
			themeFS:          testThemeFS,
			srcDir:           "testdata/content",
			templatePath:     "simple.template.html",
			luaFilterPaths:   []string{"pandoc/lua-filters/standard-code.lua"},
			expectedRendered: renderedPublicHTMLFS,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			renderedDir := mustCreateTempDir(t)

			err := RenderMultiplePages(renderedDir, tc.srcDir, tc.themeFS, tc.templatePath, tc.luaFilterPaths...)
			if tc.isError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// walk rendered FS and compares *.html pages
			renderedFS := os.DirFS(renderedDir)
			fs.WalkDir(tc.expectedRendered, ".", func(filePath string, d fs.DirEntry, pathErr error) error {
				ext := filepath.Ext(filePath)
				if ext != ".html" {
					return nil
				}

				dstB, err := fs.ReadFile(renderedFS, filePath)
				assert.NoErrorf(t, err, "cannot read (tmp) dst file")

				srcB, err := fs.ReadFile(tc.expectedRendered, filePath)
				assert.NoErrorf(t, err, "cannot read src file")

				assert.Equalf(t, dstB, srcB, "rendered %s content is wrong", filePath)
				return nil
			})
		})
	}
}
