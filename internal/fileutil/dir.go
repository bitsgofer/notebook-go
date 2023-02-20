package fileutil

import (
	"errors"
	"fmt"
	"os"

	klog "k8s.io/klog/v2"
)

// EnsureDir ensures that a director is created on the local file system.
func EnsureDir(dir string) error {
	_, err := os.Lstat(dir)
	if err == nil {
		klog.V(4).InfoS("directory already exist", "dir", dir)
		return nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("cannot get info of directory %s; err= %w", dir, err)
	}

	// dir not found, create
	if err := os.MkdirAll(dir, 0775); err != nil {
		return fmt.Errorf("cannot create directory %s; err= %w", dir, err)
	}

	klog.V(4).InfoS("directory created", "dir", dir)
	return nil
}
