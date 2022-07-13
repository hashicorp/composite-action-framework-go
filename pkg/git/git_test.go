package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/composite-action-framework-go/pkg/testhelpers/assert"
	tmp "github.com/hashicorp/composite-action-framework-go/pkg/testhelpers/tmptest"
)

func TestInit(t *testing.T) {
	dir := tmp.Dir(t)
	_, err := Init(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.DirExists(t, dir, ".git")
}

func TestOpen(t *testing.T) {
	dir := copyOfTestRepo(t)
	_, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetRemoteNamed(t *testing.T) {
	dir := copyOfTestRepo(t)

	c, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.GetRemoteNamed("origin")
	if err != nil {
		t.Fatal(err)
	}

	wantURLs := []string{"https://github.com/dadgarcorp/lockbox"}

	assert.Equal(t, got.URLs, wantURLs)
}

func copyOfTestRepo(t *testing.T) (dir string) {
	t.Helper()
	dir = tmp.Dir(t, tmp.WithContentsOf("testdata/repo1"))
	dotgit := filepath.Join(dir, "dotgit")
	dg := filepath.Join(dir, ".git")
	if err := os.Rename(dotgit, dg); err != nil {
		t.Fatal(err)
	}
	return dir
}
