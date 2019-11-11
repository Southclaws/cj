package storage

import (
	"gopkg.in/src-d/go-git.v4"
)

const wikiDir = "./wiki"

func (m *MongoStorer) PullWiki(wikiURL string) (err error) {
	if m.WikiExists() {
		var (
			repo  *git.Repository
			wtree *git.Worktree
		)

		repo, err = git.PlainOpen(wikiDir)
		if err != nil {
			return
		}

		wtree, err = repo.Worktree()
		if err != nil {
			return
		}

		wtree.Pull(&git.PullOptions{RemoteName: "origin"})
	} else {
		_, err = git.PlainClone(wikiDir, false, &git.CloneOptions{
			URL: wikiURL,
		})
	}
	return
}

// WikiExists checks if wiki directory exists on root path & is a git repository
func (m *MongoStorer) WikiExists() (exists bool) {
	_, err := git.PlainOpen(wikiDir)
	if err == git.ErrRepositoryNotExists {
		exists = false
	} else {
		exists = true
	}
	return
}
