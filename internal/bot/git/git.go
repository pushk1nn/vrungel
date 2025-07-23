package git

import (
	"fmt"

	gogit "github.com/go-git/go-git/v6"
)

type GitManager struct {
	Path string
	Repo *gogit.Repository
}

func (g *GitManager) Commit() {
	fmt.Print("Reached the commit!")
}
