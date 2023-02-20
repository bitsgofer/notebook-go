package render

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/bitsgofer/notebook-go/internal/fileutil"
	klog "k8s.io/klog/v2"
)

// Minify combines source files embeded in a theme into a single minified file.
//
// NOTE: we "cheat" by copying files from the theme to a temp directory,
// then run minify on them.
//
// This approach is less effecient than simply concat the files in memory
// and call the API provided by the minify package.
//
// We can change to that later, but for now calling the binary is easier
// than figuring out the API of minify.
func Minify(dstPath string, themeFS fs.FS, srcPaths ...string) error {
	// make temp directory
	tmpDir, tmpSrcPaths, err := copyToTempDir(themeFS, srcPaths...)
	if err != nil {
		return fmt.Errorf("cannot copy assets files to temp dir; err= %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// redirect command stdout -> dstPath
	// minified, err := os.Create(dstPath)
	// if err != nil {
	// 	return fmt.Errorf("cannot create minified file; err= %w", err)
	// }
	// setup minify command
	args := []string{
		"--output", dstPath,
		"--bundle",
	}
	args = append(args, tmpSrcPaths...)
	cmd := exec.Command("minify", args...)
	// cmd.Stdout = minified
	// access command stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("cannot get stderr of minify command; err= %w", err)
	}
	// run command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cannot start minify command; err= %w", err)
	}
	// read stderr
	minifyStderr, _ := io.ReadAll(stderr)
	// wait til command finish, then close stdout, stderr
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("cannot wait for minify cmd to finish; stderr= %s; err= %w", minifyStderr, err)
	}

	klog.V(4).InfoS("minified files", "srcs", srcPaths, "dst", dstPath)
	return nil
}

func copyToTempDir(themeFS fs.FS, srcPaths ...string) (tempDir string, tmpSrcPaths []string, copyErr error) {
	var copiedFiles []string

	// make temp directory
	dir, err := os.MkdirTemp(os.TempDir(), "notebook-minify-")
	if err != nil {
		return "", nil, fmt.Errorf("cannot create temp dir to minify files; err= %w", err)
	}

	// copy files to temp dir
	for _, src := range srcPaths {
		f := filepath.Join(dir, src)

		if err := fileutil.EnsureDir(path.Dir(f)); err != nil {
			return "", nil, fmt.Errorf("cannot create dir for temp file %s; err= %w", f, err)
		}

		if err := fileutil.CopyFromFS(f, themeFS, src); err != nil {
			return "", nil, fmt.Errorf("cannot copy theme file %s to temp file %s; err= %w", src, f, err)
		}

		// append minify sources
		copiedFiles = append(copiedFiles, f)
	}

	return dir, copiedFiles, nil

}
