package storage

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type GuildSettings struct {
	SearchMessageChannelID string   `bson:"search_message_channel_id,omitempty" json:"searchMessageChannelId"`
	ErrorReportChannelID   string   `bson:"error_report_channel_id,omitempty" json:"errorReportChannelId"`
	LeaderboardChannelID   string   `bson:"leaderboard_channel_id,omitempty" json:"leaderboardChannelId"`
	AdsChannelID           string   `bson:"ads_channel_id,omitempty" json:"adsChannelId"`
	LTFChannelID           string   `bson:"ltf_channel_id,omitempty" json:"ltfChannelId"`
	LTFUserIDs             []string `bson:"ltf_user_ids,omitempty" json:"ltfUserIds"`
}

const guildSettingsCacheKey = "guild_settings"

var guildSettingsFilter = bson.M{"kind": guildSettingsCacheKey}

func (m *MongoStorer) GetGuildSettings() (settings GuildSettings, err error) {
	if c, ok := m.cache.Get(guildSettingsCacheKey); ok {
		if settings, ok = c.(GuildSettings); ok {
			return settings, nil
		}
	}

	ctx, cancel := m.newContext()
	defer cancel()

	err = m.settings.FindOne(ctx, guildSettingsFilter).Decode(&settings)
	if err == mongo.ErrNoDocuments {
		return GuildSettings{}, nil
	}
	if err != nil {
		return
	}

	m.cache.SetDefault(guildSettingsCacheKey, settings)
	return
}

func (m *MongoStorer) SetGuildSettings(settings GuildSettings) (err error) {
	ctx, cancel := m.newContext()
	defer cancel()

	_, err = m.settings.UpdateOne(
		ctx,
		guildSettingsFilter,
		bson.M{"$set": settings, "$setOnInsert": guildSettingsFilter},
		options.UpdateOne().SetUpsert(true),
	)
	if err != nil {
		return
	}

	m.cache.SetDefault(guildSettingsCacheKey, settings)
	return
}
