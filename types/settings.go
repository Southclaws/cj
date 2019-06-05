package types

import "time"

// CommandSettings represents command configuration
type CommandSettings struct {
	Cooldown time.Duration          `bson:"cooldown,omitempty"`
	Private  bool                   `bson:"private,omitempty"`
	Channels []string               `bson:"channels,omitempty"`
	Roles    []string               `bson:"roles,omitempty"`
	Misc     map[string]interface{} `bson:"misc,omitempty"`

	Command string `bson:"command"` // internal for DB only
}
