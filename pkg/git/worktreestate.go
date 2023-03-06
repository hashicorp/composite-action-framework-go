// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package git

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"

	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/composite-action-framework-go/pkg/fs"
)

type WorktreeState struct {
	Head Commit
	// SourceHash is the SHA1 Git Commit ID of Head if the worktree is not dirty.
	// Otherwise it's a SHA1 hash of the head commit along with the full
	// contents of all changed files.
	// SHA1 is used for both so that you get a consistent 40 char hex string.
	SourceHash string
	// DirtyFiles is a list of all files flagged as dirty in the worktree.
	DirtyFiles []string
}

// IsDirty returns true if the source hash isn't exactly the same as the head
// commit ID.
func (ws WorktreeState) IsDirty() bool {
	return ws.Head.ID == ws.SourceHash
}

type WorktreeStateConfig struct {
	ignorePatterns []string
}

type WorktreeStateOption func(*WorktreeStateConfig)

func WorktreeStateIgnorePatterns(patterns ...string) WorktreeStateOption {
	return func(c *WorktreeStateConfig) { c.ignorePatterns = patterns }
}

// SourceInfo returns the SourceInfo of the repo.
func (c *Client) WorktreeState(opts ...WorktreeStateOption) (WorktreeState, error) {
	cfg := WorktreeStateConfig{}
	for _, o := range opts {
		o(&cfg)
	}
	ws := WorktreeState{}
	dirtyFiles, err := c.DirtyFiles(cfg.ignorePatterns...)
	if err != nil {
		return ws, err
	}
	ws.DirtyFiles = dirtyFiles
	commit, err := c.HeadCommit()
	if err != nil {
		return ws, err
	}
	ws.Head = commit
	if len(dirtyFiles) == 0 {
		ws.SourceHash = ws.Head.ID
		return ws, nil
	}
	summer := sha1.New()
	fmt.Fprintf(summer, "head: %s\n", commit.ID)
	for _, path := range dirtyFiles {
		if err := writeFileEntry(summer, path); err != nil {
			return ws, err
		}
	}
	ws.SourceHash = fmt.Sprintf("%x", summer.Sum(nil))
	return ws, nil
}

// writeFileEntry writes the name of the file, and its contents to w.
func writeFileEntry(w io.Writer, path string) error {
	exists, err := fs.FileExists(path)
	if err != nil {
		return err
	}
	if !exists {
		_, err := fmt.Fprintf(w, "deleted: %q\n", path)
		return err // Always return this error, we're done for deleted files.
	}
	if _, err := fmt.Fprintf(w, "changed: %q\n", path); err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	var closeErr error
	defer func() { closeErr = f.Close() }()
	if _, err := io.Copy(w, f); err != nil {
		return err
	}
	return closeErr
}

// DirtyFiles returns a list of files that have been modified in any way, including:
// new (untracked), contents changed, deleted, renamed, copied, unmerged, whether the
// changes have been staged or not. It also includes files that are ignored by the
// standard ignore files. The list is sorted using strings.Sort.
func (c *Client) DirtyFiles(excludePatterns ...string) ([]string, error) {
	wt, err := c.repo.Worktree()
	if err != nil {
		return nil, err
	}
	status, err := wt.Status()
	if err != nil {
		return nil, err
	}
	if status.IsClean() {
		return nil, nil
	}
	ignore := make([]*regexp.Regexp, len(excludePatterns))
	for i, p := range excludePatterns {
		var err error
		if ignore[i], err = regexp.Compile(p); err != nil {
			return nil, err
		}
	}
	var out []string
	for name, s := range status {
		if matchesAnyPattern(name, ignore) {
			continue
		}
		if s.Worktree != git.Unmodified || s.Staging != git.Unmodified {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out, nil
}

func matchesAnyPattern(name string, ignore []*regexp.Regexp) bool {
	for _, p := range ignore {
		if p.MatchString(name) {
			return true
		}
	}
	return false
}
