package types

import (
	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/forum"
)

// VerificationState represents the state a user's verification process is in.
type VerificationState int32

const (
	// VerificationStateNone represents the state a user is in before or after
	// the verification process. In other words, if a Verification is in this
	// state it means an error has occurred and the Verification should be
	// purged from the cache.
	VerificationStateNone VerificationState = iota

	// VerificationStateAwaitProfileURL is when the bot is waiting for the user
	// to provide their user profile page URL.
	VerificationStateAwaitProfileURL VerificationState = iota

	// VerificationStateAwaitConfirmation is when the bot is waiting for the
	// user to reply with either "done" or "cancel"
	VerificationStateAwaitConfirmation VerificationState = iota
)

// Verification holds all the state for a verification process
type Verification struct {
	ChannelID   string
	DiscordUser discordgo.User
	ForumUser   string
	UserProfile forum.UserProfile
	Code        string
	VerifyState VerificationState
}
