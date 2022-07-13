package git

import (
	"time"

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
