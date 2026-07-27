package readmodel

import (
	"fmt"

	"github.com/Southclaws/cj/discord"
	"github.com/Southclaws/cj/storage"
	"github.com/Southclaws/cj/types"
)

type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

type ServicesProvider struct {
	cfg     *types.Config
	storer  storage.Storer
	session *discord.Session
}

func NewServicesProvider(cfg *types.Config, storer storage.Storer, session *discord.Session) *ServicesProvider {
	return &ServicesProvider{cfg: cfg, storer: storer, session: session}
}

func (p *ServicesProvider) List() []ServiceStatus {
	return []ServiceStatus{
		p.mongoStatus(),
		p.discordStatus(),
		p.githubStatus(),
		p.algoliaStatus(),
	}
}

func (p *ServicesProvider) Check(name string) (ServiceStatus, error) {
	for _, status := range p.List() {
		if status.Name == name {
			return status, nil
		}
	}
	return ServiceStatus{}, fmt.Errorf("no such service: %s", name)
}

func (p *ServicesProvider) mongoStatus() ServiceStatus {
	if p.cfg.NoDatabase {
		return ServiceStatus{Name: "mongodb", Status: "disabled", Detail: "CJ_NO_DATABASE is set"}
	}
	if err := p.storer.Ping(); err != nil {
		return ServiceStatus{Name: "mongodb", Status: "unhealthy", Detail: err.Error()}
	}
	return ServiceStatus{Name: "mongodb", Status: "healthy"}
}

func (p *ServicesProvider) discordStatus() ServiceStatus {
	if p.session == nil || p.session.S == nil || p.session.S.State == nil || p.session.S.State.User == nil {
		return ServiceStatus{Name: "discord", Status: "unhealthy", Detail: "gateway session not established"}
	}
	return ServiceStatus{Name: "discord", Status: "healthy"}
}

func (p *ServicesProvider) githubStatus() ServiceStatus {
	if p.cfg.ReadmeGithubOwner == "" || p.cfg.ReadmeGithubRepository == "" || p.cfg.ReadmeFileName == "" {
		return ServiceStatus{Name: "github", Status: "not_configured"}
	}
	return ServiceStatus{
		Name:   "github",
		Status: "configured",
		Detail: p.cfg.ReadmeGithubOwner + "/" + p.cfg.ReadmeGithubRepository,
	}
}

func (p *ServicesProvider) algoliaStatus() ServiceStatus {
	return ServiceStatus{
		Name:   "algolia",
		Status: "configured",
		Detail: "search-only key embedded in the /wiki command",
	}
}
