package git

import (
	"fmt"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

type Client struct {
	opts    ClientOptions
	dir     string
	rootDir string
	repo    *git.Repository
}

type ClientOptions struct {
	authorName  string
	authorEmail string
}

type Option func(*ClientOptions)

func WithAuthor(name, email string) Option {
	return func(o *ClientOptions) {
		o.authorName = name
		o.authorEmail = email
	}
}

func Init(dir string, options ...Option) (*Client, error) {
	return newClient(dir, options, func() (*git.Repository, error) {
		return git.PlainInit(dir, false)
	})
}

func Open(dir string, options ...Option) (*Client, error) {
	return newClient(dir, options, func() (*git.Repository, error) {
		return git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{DetectDotGit: true})
	})
}

func (c *Client) RootDir() string {
	return c.rootDir
}

func (c *Client) WorkDir() string {
	return c.dir
}

func (c *Client) RepoRelativeDir() (string, error) {
	return filepath.Rel(c.rootDir, c.dir)
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

func newClient(dir string, options []Option, repoFunc func() (*git.Repository, error)) (*Client, error) {
	opts := ClientOptions{}
	for _, o := range options {
		o(&opts)
	}
	repo, err := repoFunc()
	if err != nil {
		return nil, err
	}
	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	rootDir := wt.Filesystem.Root()
	c, err := repo.Config()
	if err != nil {
		return nil, err
	}
	if opts.authorEmail == "" && c.User.Email == "" {
		opts.authorEmail = "git@example.com"
	}
	if opts.authorName == "" && c.User.Name == "" {
		opts.authorName = "Git User"
	}
	return &Client{
		opts:    opts,
		dir:     dir,
		rootDir: rootDir,
		repo:    repo,
	}, nil
}
