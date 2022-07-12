package tmptest

import (
	"os"
	"strings"
	"testing"
)

func Dir(t *testing.T) string {
	t.Helper()
	name := strings.ReplaceAll(t.Name(), "/", "_")
	f, err := os.MkdirTemp("", name+".*")
	must(t, err)
	must(t, os.Chmod(f, os.ModePerm))
	return f
}
