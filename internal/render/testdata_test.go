package render

import (
	_ "embed"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testThemeFS = os.DirFS("testdata/testTheme")

	//go:embed testdata/expectedRendered/css/minified.css
	minifiedCSS []byte

	//go:embed testdata/expectedRendered/js/minified.js
	minifiedJS []byte

	//go:embed testdata/expectedRendered/hello/index.html
	renderedHelloHTML []byte

	//go:embed testdata/expectedRendered/about-me/index.html
	renderedAboutMeHTML []byte

	renderedPublicHTMLFS = os.DirFS("testdata/expectedRendered")
)

// mustCreateTempFile returns a unique file path that could be used for a temp file.
// While the path is usable, a temp file at that path won't be created.
func mustCreateTempFile(t *testing.T) string {
	f, err := os.CreateTemp(os.TempDir(), "notebook.file.*")
	require.NoErrorf(t, err, "cannot create temp file")

	err = f.Close()
	require.NoErrorf(t, err, "cannot close temp file")

	return f.Name()
}

// mustCreateTempDir creates a dir under os.TempDir()
func mustCreateTempDir(t *testing.T) string {
	// create a file, then remove it. Use the name for the tmp placeholder
	dir, err := os.MkdirTemp(os.TempDir(), "notebook.dir.*")
	require.NoErrorf(t, err, "cannot create temp file (placeholder for temp dir)")

	return dir
}
