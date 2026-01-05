// Copyright IBM Corp. 2022, 2025
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
