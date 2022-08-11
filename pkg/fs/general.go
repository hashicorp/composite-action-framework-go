package fs

import (
	"errors"
	"os"
)

// Move is like os.Rename except it first ensures there's nothing at dest by
// deleting anything there.
func Move(oldPath, newPath string) error {
	if err := os.RemoveAll(newPath); err != nil {
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
