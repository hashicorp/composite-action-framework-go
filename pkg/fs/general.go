// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package fs

import (
	"errors"
	"fmt"
	"os"
)

// Move is like os.Rename except it first ensures there's nothing at dest by
// deleting anything there.
func Move(oldPath, newPath string, opts ...Option) error {
	return New(opts...).Move(oldPath, newPath)
}

func (fs *FS) Move(oldPath, newPath string) error {
	exists, err := Exists(oldPath)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%s does not exist", oldPath)
	}
	if err := os.RemoveAll(newPath); err != nil {
		return err
	}
	if err := fs.prepareContainingDir(newPath); err != nil {
		return err
	}
	return os.Rename(oldPath, newPath)
}

func Exists(name string) (bool, error) {
	return existsAndPassesTest(name, func(os.FileInfo) bool { return true })
}

func existsAndPassesTest(name string, test func(os.FileInfo) bool) (bool, error) {
	info, exists, err := stat(name)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	return test(info), nil
}

func stat(name string) (os.FileInfo, bool, error) {
	info, err := os.Stat(name)
	if err == nil {
		return info, true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	return nil, false, err
}
