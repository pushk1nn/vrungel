package git

import (
	"fmt"
	"time"

	gogit "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
)

type GitManager struct {
	Path string
	Repo *gogit.Repository
	Auth *http.BasicAuth
}

func (g *GitManager) Commit(filename string) {
	r, err := gogit.PlainOpen(g.Path)
	if err != nil {
		panic(err)
	}

	w, err := r.Worktree()
	if err != nil {
		panic(err)
	}

	_, err = w.Add("rolebinding-constraints/" + filename)
	if err != nil {
		panic(err)
	}

	commit, err := w.Commit("Automated constraint generation", &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  "Vrungel",
			Email: "vrungel@maxvk.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		panic(err)
	}

	obj, err := r.CommitObject(commit)
	if err != nil {
		panic(err)
	}

	fmt.Println(obj)
}

func (g *GitManager) Push() {
	r, err := gogit.PlainOpen(g.Path)
	if err != nil {
		panic(err)
	}

	err = r.Push(&gogit.PushOptions{
		Auth: g.Auth,
	})
	if err != nil {
		panic(err)
	}
}
