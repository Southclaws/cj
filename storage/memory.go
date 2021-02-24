package storage

import (
	"github.com/Southclaws/cj/types"
	"github.com/bwmarrin/discordgo"
	"github.com/google/go-github/v28/github"
)

type Memory struct{}

func (m *Memory) GetUserOrCreate(discordUserID string) (result User) {
	return
}
func (m *Memory) RecordChatLog(discordUserID string, discordChannel string, message string, messageID string) (err error) {
	return
}
func (m *Memory) GetMessagesForUser(discordUserID string) (messages []ChatLog, err error) {
	return
}
func (m *Memory) GetTopMessages(top int) (result TopMessages, err error) {
	return
}
func (m *Memory) GetUserMessageCount(discordUserID string) (messageCount int, err error) {
	return
}
func (m *Memory) GetUserRank(discordUserID string) (rank int, err error) {
	return
}
func (m *Memory) GetRandomMessage() (result ChatLog, err error) {
	return
}
func (m *Memory) GetRandomUser() (result string, err error) {
	return
}
func (m *Memory) GetMessageByID(messageID string) (message ChatLog, err error) {
	return
}
func (m *Memory) UpdateUserUsername(discordUserID string, username string) (err error) {
	return
}
func (m *Memory) UpdateUser(user User) (err error) {
	return
}
func (m *Memory) RemoveUser(id string) (err error) {
	return
}
func (m *Memory) IsUserVerified(discordUserID string) (verified bool, err error) {
	return
}
func (m *Memory) IsUserLegacyVerified(discordUserID string) (verified bool, err error) {
	return
}
func (m *Memory) GetDiscordUserForumUser(forumUserID string) (discordUserID string, err error) {
	return
}
func (m *Memory) GetForumUserFromDiscordUser(discordUserID string) (legacyUserID string, burgerUserID string, err error) {
	return
}
func (m *Memory) GetForumNameFromDiscordUser(discordUserID string) (legacyUserName string, burgerUserName string, err error) {
	return
}
func (m *Memory) GetDiscordUserFromForumName(forumName string) (legacyUserID string, burgerUserID string, err error) {
	return
}
func (m *Memory) AddEmojiReactionToUser(discordUserID string, emoji string) (err error) {
	return
}
func (m *Memory) RemoveEmojiReactionFromUser(discordUserID string, emoji string) (err error) {
	return
}
func (m *Memory) GetTopReactions(top int, reaction string) (result []TopReactionEntry, err error) {
	return
}
func (m *Memory) SetCommandSettings(command string, settings types.CommandSettings) (err error) {
	return
}
func (m *Memory) GetCommandSettings(command string) (settings types.CommandSettings, found bool, err error) {
	return
}
func (m *Memory) GetReadmeMessage() (message string, err error) {
	return
}
func (m *Memory) FetchReadmeMessage(gistID string, gistFile github.GistFilename) (message string, err error) {
	return
}
func (m *Memory) UpdateReadmeMessage(session *discordgo.Session, original *discordgo.Message, upstream string) (err error) {
	return
}
