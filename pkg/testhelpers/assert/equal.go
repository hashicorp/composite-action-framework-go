// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package assert

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Equal(t *testing.T, got, want interface{}) {
	t.Helper()
	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Errorf("Mismatch (-got +want):\n%s", diff)
	}
}
