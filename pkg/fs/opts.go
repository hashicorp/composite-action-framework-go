// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fs

import (
	"fmt"
	"path/filepath"
)

type Settings struct {
	overwrite  bool
	createDirs bool
}

func newSettings(opts []Option) Settings {
	s := &Settings{
		overwrite:  true,
		createDirs: true,
	}
	for _, o := range opts {
		o(s)
	}
	return *s
}

func (s *Settings) prepareContainingDir(name string) error {
	dir := filepath.Dir(name)
	exists, err := DirExists(name)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	if !s.createDirs {
		return fmt.Errorf("directory %s does not exist", dir)
	}
	return Mkdir(dir)
}

type Option func(*Settings)

func WithOverwrite(t bool) Option  { return func(s *Settings) { s.overwrite = t } }
func WithCreateDirs(t bool) Option { return func(s *Settings) { s.createDirs = t } }
