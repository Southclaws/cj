package main

func commandSay(cm CommandManager, text string, channel string) error {
	_, err := cm.App.discordClient.ChannelMessageSend(cm.App.config.PrimaryChannel, text)
	return err
}
