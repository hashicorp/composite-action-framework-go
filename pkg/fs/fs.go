// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fs

type FS struct {
	Settings
}

func New(opts ...Option) *FS {
	return &FS{
		Settings: newSettings(opts),
	}
}
