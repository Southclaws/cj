package stats

func (a *Aggregator) gatherTopMessages(top int) (err error) {
	a.topMessages, err = a.Storage.GetTopMessages(top)
	return
}
