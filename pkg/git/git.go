package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

type GitClient struct {
	dir string
}

func NewGitClient(dir string) *GitClient {
	return &GitClient{dir: dir}
}

func Init(dir string) (*git.Repository, error) {
	return git.PlainInit(dir, false)
}

func Open(dir string) (*git.Repository, error) {
	return git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{DetectDotGit: true})
}

func GetConfig(dir string) (*config.Config, error) {
	configFile, err := os.Open(filepath.Join(dir, ".git", "config"))
	if err != nil {
		return nil, err
	}
	var closeErr error
	defer func() { closeErr = configFile.Close() }()

	c, err := config.ReadConfig(configFile)
	if err != nil {
		return nil, err
	}

	return c, closeErr
}

func GetRemote(dir, name string) (*config.RemoteConfig, error) {
	c, err := GetConfig(dir)
	if err != nil {
		return nil, err
	}
	r, ok := c.Remotes[name]
	if !ok {
		return nil, fmt.Errorf("no remote named %q", name)
	}
	return r, nil
}
