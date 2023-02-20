package fileutil

import (
	"fmt"
	"io"
	"io/fs"
	"os"

	klog "k8s.io/klog/v2"
)

// CopyFromFS copies a file from a fs.FS to a location on the local file system.
func CopyFromFS(dst string, srcFS fs.FS, src string) error {
	b, err := fs.ReadFile(srcFS, src)
	if err != nil {
		return fmt.Errorf("cannot read file from FS; err= %w", err)
	}

	toFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("cannot create destination file %s; err= %w", dst, err)
	}
	defer toFile.Close()

	if _, err := toFile.Write(b); err != nil {
		return fmt.Errorf("cannot write to %s; err= %w", dst, err)
	}

	klog.V(4).InfoS("file copied from FS", "src", src, "dst", dst)
	return nil
}

// CopyFrom copies a file on the local file system to another location.
func CopyFrom(dst, src string) error {
	fromFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("cannot open source file %s; err= %w", src, err)
	}
	defer fromFile.Close()

	toFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("cannot create destination file %s; err= %w", dst, err)
	}
	defer toFile.Close()

	if _, err := io.Copy(toFile, fromFile); err != nil {
		return fmt.Errorf("cannot copy %s to %s; err= %w", src, dst, err)
	}

	klog.V(4).InfoS("file copied", "src", src, "dst", dst)
	return nil
}
