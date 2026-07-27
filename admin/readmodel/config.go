package readmodel

import "github.com/Southclaws/cj/types"

type ConfigField struct {
	Name string `json:"name"`
	Set  bool   `json:"set"`
}

type ConfigStatus struct {
	DatabaseEnabled bool          `json:"databaseEnabled"`
	Fields          []ConfigField `json:"fields"`
}

type ConfigStatusProvider struct {
	cfg *types.Config
}

func NewConfigStatusProvider(cfg *types.Config) *ConfigStatusProvider {
	return &ConfigStatusProvider{cfg: cfg}
}

func (p *ConfigStatusProvider) Status() ConfigStatus {
	cfg := p.cfg
	return ConfigStatus{
		DatabaseEnabled: !cfg.NoDatabase,
		Fields: []ConfigField{
			{Name: "discord_token", Set: cfg.DiscordToken != ""},
			{Name: "guild_id", Set: cfg.GuildID != ""},
			{Name: "bot_id", Set: cfg.BotID != ""},
			{Name: "default_channel", Set: cfg.DefaultChannel != ""},
			{Name: "default_role", Set: cfg.DefaultRole != ""},
			{Name: "readme_channel", Set: cfg.ReadmeChannel != ""},
		},
	}
}
