package readmodel

import "github.com/Southclaws/cj/bot/heartbeat"

type JobsProvider struct {
	heartbeat *heartbeat.Heartbeat
}

func NewJobsProvider(hb *heartbeat.Heartbeat) *JobsProvider {
	return &JobsProvider{heartbeat: hb}
}

func (p *JobsProvider) List() []heartbeat.JobStatus {
	if p.heartbeat == nil {
		return []heartbeat.JobStatus{}
	}
	statuses := p.heartbeat.Statuses()
	if statuses == nil {
		return []heartbeat.JobStatus{}
	}
	return statuses
}
