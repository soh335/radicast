package main

import (
	"io"
	"os"
)

func RenameOrCopy(oldpath, newpath string) error {
	// failed if oldpath and newpath are not at same file system.
	if err := os.Rename(oldpath, newpath); err == nil {
		return nil
	}
	return _copy(oldpath, newpath)
}

func _copy(oldpath, newpath string) error {
	o, err := os.Open(oldpath)
	if err != nil {
		return err
	}

	defer o.Close()

	n, err := os.Create(newpath)
	if err != nil {
		return err
	}

	defer n.Close()

	if _, err := io.Copy(n, o); err != nil {
		return err
	}

	defer os.Remove(oldpath)

	return nil
}
