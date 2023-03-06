// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import "errors"

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrNoArgsAllowed  = errors.New("no args allowed")
)
