package storage

import (
	"context"
	"time"

	"gopkg.in/src-d/go-git.v4"
)

// EnsureWiki ensures that wiki is cloned and pulls it in every hour
func EnsureWiki(wikiURL string) (err error) {
	if WikiExists() {
		ctx, cancelFunc := context.WithCancel(context.Background())
		defer cancelFunc()
		err = pullHourly(ctx)
	} else {
		_, err = git.PlainClone("./wiki", false, &git.CloneOptions{
			URL: wikiURL,
		})
	}
	return err
}

func pullHourly(ctx context.Context) (err error) {
	var (
		repo  *git.Repository
		wtree *git.Worktree
	)

	repo, err = git.PlainOpen("./wiki")
	if err != nil {
		return
	}
	wtree, err = repo.Worktree()
	if err != nil {
		return
	}
	go func() {
		tick := time.NewTicker(time.Hour)

		defer tick.Stop()

		for {
			select {
			case <-tick.C:
				wtree.Pull(&git.PullOptions{RemoteName: "origin"})
			case <-ctx.Done():
				return
			}
		}
	}()
	return
}

// WikiExists checks if wiki directory exists on root path & is a git repository
func WikiExists() bool {
	_, err := git.PlainOpen("./wiki")
	if err == git.ErrRepositoryNotExists {
		return false
	}
	return true
}
