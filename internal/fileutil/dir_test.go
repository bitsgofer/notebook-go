package fileutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureDir(t *testing.T) {
	testCases := map[string]struct {
		setup   func() // run this before test
		dir     string
		isError bool
	}{
		"dir-not-exist": {
			dir: filepath.Join(os.TempDir(), "notebook-test-ensureDir"),
		},
		"dir-already-exist": {
			setup: func() {
				os.MkdirAll(filepath.Join(os.TempDir(), "notebook-test-ensureDir-already-exist"), 0775)
			},
			dir: filepath.Join(os.TempDir(), "notebook-test-ensureDir-already-exist"),
		},
		// "dir-not-exist-no-permission-to-create": {},
		// NOTE: to test this, we probably need to convert EnsureDir() to work with fs.FS
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			defer os.RemoveAll(tc.dir) // clean up after test

			err := EnsureDir(tc.dir)
			if tc.isError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// check that directory exist
			info, err := os.Lstat(tc.dir)
			assert.NoErrorf(t, err, "expected os.Lstat(%s) to work", tc.dir)
			// err != nil => exist
			assert.Truef(t, info.IsDir(), "ensured path is not a directory")
		})
	}
}
