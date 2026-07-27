package readmodel

import "time"

type Status struct {
	Version   string    `json:"version"`
	StartedAt time.Time `json:"startedAt"`
	Uptime    string    `json:"uptime"`
}

type StatusProvider struct {
	version   string
	startedAt time.Time
}

func NewStatusProvider(version string, startedAt time.Time) *StatusProvider {
	return &StatusProvider{version: version, startedAt: startedAt}
}

func (p *StatusProvider) Status() Status {
	return Status{
		Version:   p.version,
		StartedAt: p.startedAt,
		Uptime:    time.Since(p.startedAt).Round(time.Second).String(),
	}
}
