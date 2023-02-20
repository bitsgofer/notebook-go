package render

import (
	"io/fs"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinify(t *testing.T) {
	testCases := map[string]struct {
		themeFS  fs.FS
		srcPaths []string
		// output
		isError          bool
		expectedMinified []byte
	}{
		"minify-css": {
			themeFS: testThemeFS,
			srcPaths: []string{
				"one.css",
				"two.css",
			},
			expectedMinified: minifiedCSS,
		},
		"minify-js": {
			themeFS: testThemeFS,
			srcPaths: []string{
				"one.js",
				"two.js",
			},
			expectedMinified: minifiedJS,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			minifiedFile := mustCreateTempFile(t)

			err := Minify(minifiedFile, tc.themeFS, tc.srcPaths...)
			if tc.isError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// check minified content
			dstB, err := ioutil.ReadFile(minifiedFile)
			assert.NoErrorf(t, err, "cannot read (tmp) dst file")
			assert.Equalf(t, dstB, tc.expectedMinified, "minified content is wrong")
		})
	}
}
