package types

import "time"

// CommandSettings represents command configuration
type CommandSettings struct {
	Cooldown time.Duration          `bson:"cooldown,omitempty"`
	Roles    []string               `bson:"roles,omitempty"`
	Misc     map[string]interface{} `bson:"misc,omitempty"`
	Command  string                 `bson:"command"` // internal for DB only
}
