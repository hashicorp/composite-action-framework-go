// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tmptest

import (
	"os"
	"strings"
	"testing"

	cp "github.com/otiai10/copy"
)

type Settings struct {
	contentsFrom    string
	topLevelTempDir string
}

type Option func(*Settings)

func WithContentsOf(dir string) Option {
	return func(s *Settings) { s.contentsFrom = dir }
}

func WithTopLevelTempDir(dir string) Option {
	return func(s *Settings) { s.topLevelTempDir = dir }
}

func Dir(t *testing.T, opts ...Option) string {
	t.Helper()
	s := &Settings{}
	for _, o := range opts {
		o(s)
	}
	name := strings.ReplaceAll(t.Name(), "/", "_")
	dir, err := os.MkdirTemp(s.topLevelTempDir, name+".*")
	must(t, err)
	must(t, os.Chmod(dir, os.ModePerm))

	if s.contentsFrom != "" {
		must(t, cp.Copy(s.contentsFrom, dir))
	}

	return dir
}
