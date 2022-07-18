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

// TestOpen tests opening an existing repo and checks that the client
// is functional, and returns the correct directories.
func TestOpen(t *testing.T) {

	cases := []struct {
		desc, reldir string
	}{
		{"root", "."},
		{"subdir", "subdir"},
	}

	for _, c := range cases {
		desc, reldir := c.desc, c.reldir
		t.Run(desc, func(t *testing.T) {
			dir := copyOfTestRepo(t)
			subdir := filepath.Join(dir, reldir)
			r, err := Open(subdir)
			if err != nil {
				t.Fatal(err)
			}
			l, err := r.Log(1)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(l[0])
			if rootDir := r.RootDir(); rootDir != dir {
				t.Errorf("RootDir() = %q; want %q", rootDir, dir)
			}
			if workDir := r.WorkDir(); workDir != subdir {
				t.Errorf("RootDir() = %q; want %q", workDir, subdir)
			}
			repoRelativeDir, err := r.RepoRelativeDir()
			if err != nil {
				t.Fatal(err)
			}
			if repoRelativeDir != reldir {
				t.Errorf("RepoRelativeDir() = %q; want %q", repoRelativeDir, reldir)
			}
		})
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
