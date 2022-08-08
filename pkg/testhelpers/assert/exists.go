package assert

import (
	"path/filepath"
	"testing"

	"github.com/hashicorp/composite-action-framework-go/pkg/fs"
)

func DirExists(t *testing.T, path ...string) {
	t.Helper()
	exists, err := fs.DirExistsJoin(path...)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatalf("dir %q does not exist", filepath.Join(path...))
	}
}
