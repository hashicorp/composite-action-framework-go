package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

type Client struct {
	dir  string
	repo *git.Repository
}

func Init(dir string) (*Client, error) {
	return newClient(dir, func() (*git.Repository, error) {
		return git.PlainInit(dir, false)
	})
}

func Open(dir string) (*Client, error) {
	return newClient(dir, func() (*git.Repository, error) {
		return git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{DetectDotGit: true})
	})
}

func (c *Client) Log(n int) ([]Commit, error) {
	out := make([]Commit, 0, n)
	logIter, err := c.repo.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}
	defer logIter.Close()
	for i := 0; i < n; i++ {
		commit, err := logIter.Next()
		if err != nil {
			return nil, err
		}
		if commit == nil {
			break
		}
		out = append(out, newCommit(commit))
	}
	return out, nil
}

type Config = config.Config

func (c *Client) Config() (*Config, error) {
	return c.repo.Config()
	//configFile, err := os.Open(filepath.Join(c.dir, ".git", "config"))
	//if err != nil {
	//	return nil, err
	//}
	//var closeErr error
	//defer func() { closeErr = configFile.Close() }()

	//c, err := config.ReadConfig(configFile)
	//if err != nil {
	//	return nil, err
	//}

	//return c, closeErr
}

func (c *Client) GetRemoteNamed(name string) (*config.RemoteConfig, error) {
	cfg, err := c.Config()
	if err != nil {
		return nil, err
	}
	r, ok := cfg.Remotes[name]
	if !ok {
		return nil, fmt.Errorf("no remote named %q", name)
	}
	return r, nil
}

func newClient(dir string, repoFunc func() (*git.Repository, error)) (*Client, error) {
	repo, err := repoFunc()
	if err != nil {
		return nil, err
	}
	return &Client{
		dir:  dir,
		repo: repo,
	}, nil

}
