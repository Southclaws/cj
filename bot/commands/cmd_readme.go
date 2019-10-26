package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/globalsign/mgo"

	"github.com/Southclaws/cj/types"
)

func (cm *CommandManager) commandReadme(
	args string,
	message discordgo.Message,
	contextual bool,
	settings types.CommandSettings,
) (
	context bool,
	err error,
) {
	// get the channel where readme message is written in
	channel, e := cm.Discord.S.Channel(cm.Config.ReadmeChannel)
	if e != nil {
		err = e
		return
	}

	// the content of upstream gist
	upstream, e := cm.Storage.FetchReadmeMessage(cm.Config.ReadmeGistID, cm.Config.ReadmeGistFileName)
	if e != nil {
		err = e
		return
	}

	// the message id, stored in database
	messageID, e := cm.Storage.GetReadmeMessage()

	// if it's not in the database
	if e == mgo.ErrNotFound {
		// we send a new message in the channel and force-update database
		msg, e := cm.Discord.S.ChannelMessageSend(channel.ID, upstream)
		if e != nil {
			err = e
			return
		}
		err = cm.Storage.UpdateReadmeMessage(cm.Discord.S, msg, upstream)
		if err != nil {
			return
		}
	} else if e != nil {
		err = e
		return
	} else {
		// if it's in the database, then we force it to update if needed
		msg, e := cm.Discord.S.ChannelMessage(channel.ID, messageID)
		if e != nil {
			res, ok := e.(*discordgo.RESTError)
			if ok {
				if res.Message.Code == discordgo.ErrCodeUnknownMessage {
					msg, e := cm.Discord.S.ChannelMessageSend(channel.ID, upstream)
					if e != nil {
						err = e
						return
					}
					err = cm.Storage.UpdateReadmeMessage(cm.Discord.S, msg, upstream)
					if err != nil {
						return
					}
				}
			}
		} else {
			err = cm.Storage.UpdateReadmeMessage(cm.Discord.S, msg, upstream)
			if err != nil {
				return
			}
		}
	}

	return
}
