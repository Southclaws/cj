package storage

import (
	"github.com/Southclaws/cj/types"
	"github.com/globalsign/mgo/bson"
	"gopkg.in/mgo.v2"
)

// SetCommandSettings upsets command settings
func (m *MongoStorer) SetCommandSettings(command string, settings types.CommandSettings) (err error) {
	settings.Command = command
	_, err = m.settings.Upsert(bson.M{"command": command}, settings)
	return
}

// GetCommandSettings returns command settings and uses a cache
func (m *MongoStorer) GetCommandSettings(command string) (settings types.CommandSettings, found bool, err error) {
	if c, ok := m.cache.Get(command); ok {
		if settings, ok = c.(types.CommandSettings); ok {
			return
		}
	}
	err = m.settings.Find(bson.M{"command": command}).One(&settings)
	if err == mgo.ErrNotFound {
		err = nil
	} else if err == nil {
		found = true
	}
	return
}
