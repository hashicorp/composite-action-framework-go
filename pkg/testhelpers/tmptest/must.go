// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tmptest

import "testing"

// must is a quick way to fail a test depending on if an error is nil or not.
func must(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
