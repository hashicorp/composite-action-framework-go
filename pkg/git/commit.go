// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package git

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Commit struct {
	ID          string
	AuthorName  string
	AuthorEmail string
	AuthorTime  time.Time
	Message     string
	Date        time.Time
}

func newCommit(c *object.Commit) Commit {
	return Commit{
		ID:          c.ID().String(),
		AuthorName:  c.Author.Name,
		AuthorEmail: c.Author.Email,
		AuthorTime:  c.Author.When,
		Message:     c.Message,
	}
}

func (c *Client) HeadCommit() (Commit, error) {
	h, err := c.repo.Head()
	if err != nil {
		return Commit{}, err
	}
	commit, err := c.repo.CommitObject(h.Hash())
	if err != nil {
		return Commit{}, err
	}
	return newCommit(commit), nil
}

func (c *Client) Add(paths ...string) error {
	wt, err := c.repo.Worktree()
	if err != nil {
		return err
	}
	for _, p := range paths {
		if _, err := wt.Add(p); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Commit(message string) error {
	wt, err := c.repo.Worktree()
	if err != nil {
		return err
	}
	opts := &git.CommitOptions{}
	if c.opts.authorName != "" {
		if opts.Author == nil {
			opts.Author = &object.Signature{}
		}
		opts.Author.Email = c.opts.authorEmail
	}
	if c.opts.authorName != "" {
		if opts.Author == nil {
			opts.Author = &object.Signature{}
		}
		opts.Author.Name = c.opts.authorName
	}
	_, err = wt.Commit(message, opts)
	return err
}
