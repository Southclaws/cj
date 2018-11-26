package storage

import "github.com/Southclaws/cj/types"

type Memory struct{}

func (m *Memory) RecordChatLog(discordUserID string, discordChannel string, message string) (err error) {
	return
}
func (m *Memory) GetMessagesForUser(discordUserID string) (messages []ChatLog, err error) {
	return
}
func (m *Memory) GetTopMessages(top int) (result TopMessages, err error) {
	return
}
func (m *Memory) StoreVerifiedUser(verification types.Verification) (err error) {
	return
}
func (m *Memory) RemoveUser(id string) (err error) {
	return
}
func (m *Memory) IsUserVerified(discordUserID string) (verified bool, err error) {
	return
}
func (m *Memory) GetDiscordUserForumUser(forumUserID string) (discordUserID string, err error) {
	return
}
func (m *Memory) GetForumUserFromDiscordUser(discordUserID string) (forumUserID string, err error) {
	return
}
func (m *Memory) GetForumNameFromDiscordUser(discordUserID string) (forumUserName string, err error) {
	return
}
func (m *Memory) GetDiscordUserFromForumName(forumName string) (discordUserID string, err error) {
	return
}
