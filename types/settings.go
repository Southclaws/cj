package types

import "time"

// CommandSettings represents command configuration
type CommandSettings struct {
	Cooldown *time.Duration         `bson:"cooldown,omitempty"`
	Channels []string               `bson:"channels,omitempty"`
	Misc     map[string]interface{} `bson:"misc,omitempty"`

	Command string `bson:"command"` // internal for DB only
}
