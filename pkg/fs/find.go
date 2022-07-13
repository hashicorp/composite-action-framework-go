package fs

import (
	"io/fs"
	"path/filepath"
)

// FindFilesNamed returns the paths of all files with the specified name
// in the file tree rooted at dir.
func FindFilesNamed(dir, name string) ([]string, error) {
	return FindFiles(dir, func(d fs.DirEntry, path string) bool {
		return d.Name() == name
	})
}

type FindPredicate func(d fs.DirEntry, path string) bool

// FindFiles looks for files in the repo, excluding the .git dir.
func FindFiles(dir string, predicate FindPredicate) ([]string, error) {
	var got []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return fs.SkipDir
			}
			return nil
		}
		if predicate(d, path) {
			got = append(got, path)
		}
		return nil
	})
	return got, err
}
