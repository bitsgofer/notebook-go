package fileutil

import (
	"embed"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata
var testDataFS embed.FS

func TestCopyFromFS(t *testing.T) {
	testCases := map[string]struct {
		setup   func() // call before test
		srcFS   fs.FS
		src     string
		dst     string
		isError bool
	}{
		"copy-to-new-location-succeed": {
			srcFS: testDataFS,
			src:   "testdata/hello",
			dst:   filepath.Join(os.TempDir(), "notebook-test-copyFromFS"),
		},
		"copy-to-existing-file-succeed": { // file is truncated
			setup: func() {
				f, _ := os.Create(filepath.Join(os.TempDir(), "notebook-test-copyFromFS-alreayd-exist"))
				f.Close()
			},
			srcFS: testDataFS,
			src:   "testdata/hello",
			dst:   filepath.Join(os.TempDir(), "notebook-test-copyFromFS-alreayd-exist"),
		},
		"copy-to-file-that-cannot-be-create": { // no upper level directory
			srcFS:   testDataFS,
			src:     "testdata/hello",
			dst:     filepath.Join("/a/path/that/donnot/exist", "notebook-test-copyFromFS"),
			isError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			defer os.RemoveAll(tc.dst)

			err := CopyFromFS(tc.dst, tc.srcFS, tc.src)
			if tc.isError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// compare content
			srcB, err := fs.ReadFile(tc.srcFS, tc.src)
			assert.NoErrorf(t, err, "cannot read src file (in srcFS)")
			dstB, err := ioutil.ReadFile(tc.dst)
			assert.NoErrorf(t, err, "cannot read dst file")
			assert.Equalf(t, dstB, srcB, "copied file is not the same as src")
		})
	}
}
