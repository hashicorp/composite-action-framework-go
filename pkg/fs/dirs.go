// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fs

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func DirExists(name string) (bool, error) {
	return existsAndPassesTest(name, func(info os.FileInfo) bool {
		return info.IsDir()
	})
}

// DirExistsJoin checks if the dir named by segments exists.
// Segments is a set of path segments that may or may not themselves
// contain path separators. The following three calls are all equivalent:
//
//   DirExistsJoin("some/long/path")
//   DirExistsJoin("some", "long/path")
//   DirExistsJoin("some", "long", "path")
func DirExistsJoin(name ...string) (bool, error) {
	return DirExists(filepath.Join(name...))
}

// Mkdir makes the directory at path, using default permissions.
func Mkdir(path string) error {
	return os.MkdirAll(path, fs.ModePerm)
}

// MkdirEmpty deletes any existing file or directory at path, and then creates
// a new empty directory at path, using default permissions.
func MkdirEmpty(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return Mkdir(path)
}

// Mkdirs calls Mkdir sequentially on paths and returns an error after the first failure.
func Mkdirs(paths ...string) error {
	for _, p := range paths {
		if err := Mkdir(p); err != nil {
			return err
		}
	}
	return nil
}

// SetMtime sets the mtime of all files inside dir to the provided time.
func SetMtimes(dir string, to time.Time) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		filePath := filepath.Join(dir, e.Name())
		if err := os.Chtimes(filePath, to, to); err != nil {
			return err
		}
	}
	return nil
}
