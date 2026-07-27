package readmodel

import "github.com/Southclaws/cj/storage"

type ActionRunsProvider struct {
	storer storage.Storer
}

func NewActionRunsProvider(storer storage.Storer) *ActionRunsProvider {
	return &ActionRunsProvider{storer: storer}
}

func (p *ActionRunsProvider) List(limit int) ([]storage.ActionRun, error) {
	return p.storer.ListActionRuns(limit)
}

func (p *ActionRunsProvider) Get(id string) (storage.ActionRun, bool, error) {
	return p.storer.GetActionRun(id)
}
