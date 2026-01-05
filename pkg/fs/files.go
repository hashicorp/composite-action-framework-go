// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Create(name string, opts ...Option) (*os.File, error) {
	fs := New(opts...)
	return fs.Create(name)
}

func Append(name string, opts ...Option) (*os.File, error) {
	fs := New(opts...)
	return fs.Append(name)
}

func (fs *FS) Append(name string) (*os.File, error) {
	exists, err := FileExists(name)
	if err != nil {
		return nil, err
	}
	if exists {
		return os.OpenFile(name, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	}
	return fs.Create(name)
}

func (fs *FS) Create(name string) (*os.File, error) {
	if err := fs.prepareContainingDir(name); err != nil {
		return nil, err
	}
	dirExists, err := DirExists(name)
	if err != nil {
		return nil, err
	}
	if dirExists {
		return nil, fmt.Errorf("%s exits and is a directory", name)
	}
	fileExists, err := FileExists(name)
	if err != nil {
		return nil, err
	}
	if fileExists && !fs.overwrite {
		return nil, fmt.Errorf("%s exists and is a file", name)
	}
	return os.Create(name)
}

func CreateOverwrite(name string) (*os.File, error) {
	if err := Mkdir(filepath.Dir(name)); err != nil {
		return nil, err
	}
	return os.Create(name)
}

// FileExists returns a boolean indicating that name is a real path
// and is not a directory.
func FileExists(name string) (bool, error) {
	return existsAndPassesTest(name, func(info os.FileInfo) bool {
		return !info.IsDir()
	})
}

// WriteFile writes a file to the specified path, with default permissions, and
// creates any needed directories.
func WriteFile[T Bytes](path string, contents T) error {
	if err := Mkdir(filepath.Dir(path)); err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(contents), os.ModePerm)
}

// WriteTempFile writes contents to a unique temporary file and returns its path.
func WriteTempFile[T Bytes](name string, contents T) (string, error) {
	return WithTempFile(name, func(w io.Writer) error {
		_, err := w.Write([]byte(contents))
		return err
	})
}

// WithTempFile creates a unique temporary file, runs the 'do' function you provide,
// passing the file as an io.Writer, and then closes the file. It returns the
// file's path, or an error from 'do', or a general write or close error.
func WithTempFile(name string, do func(io.Writer) error) (string, error) {
	name = fmt.Sprintf("%s.*", name)
	tempFile, err := os.CreateTemp("", name)
	if err != nil {
		return "", err
	}
	if err := do(tempFile); err != nil {
		tempFile.Close()
		return "", err
	}
	return tempFile.Name(), tempFile.Close()
}
