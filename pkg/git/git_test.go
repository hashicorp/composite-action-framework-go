package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/composite-action-framework-go/pkg/testhelpers/assert"
	tmp "github.com/hashicorp/composite-action-framework-go/pkg/testhelpers/tmptest"
)

func TestGetRemote(t *testing.T) {

	dir := tmp.Dir(t, tmp.WithContentsOf("testdata/repo1"))
	dotgit := filepath.Join(dir, "dotgit")
	dg := filepath.Join(dir, ".git")
	if err := os.Rename(dotgit, dg); err != nil {
		t.Fatal(err)
	}

	got, err := GetRemote(dir, "origin")
	if err != nil {
		t.Fatal(err)
	}

	wantURLs := []string{"https://github.com/dadgarcorp/lockbox"}

	assert.Equal(t, got.URLs, wantURLs)
}
