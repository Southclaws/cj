package heartbeat

import (
	"sort"
	"sync"
	"time"
)

type JobStatus struct {
	Name      string     `json:"name"`
	Provider  string     `json:"provider"`
	Schedule  string     `json:"schedule"`
	LastRun   *time.Time `json:"lastRun,omitempty"`
	LastError string     `json:"lastError,omitempty"`
	RunCount  int        `json:"runCount"`
}

type jobRegistry struct {
	mu     sync.Mutex
	byName map[string]*JobStatus
}

func newJobRegistry() *jobRegistry {
	return &jobRegistry{byName: make(map[string]*JobStatus)}
}

func (r *jobRegistry) register(name, provider, schedule string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byName[name] = &JobStatus{Name: name, Provider: provider, Schedule: schedule}
}

func (r *jobRegistry) recordRun(name string, runErr error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	status, ok := r.byName[name]
	if !ok {
		return
	}
	now := time.Now().UTC()
	status.LastRun = &now
	status.RunCount++
	status.LastError = ""
	if runErr != nil {
		status.LastError = runErr.Error()
	}
}

func (r *jobRegistry) snapshot() []JobStatus {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]JobStatus, 0, len(r.byName))
	for _, s := range r.byName {
		out = append(out, *s)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
