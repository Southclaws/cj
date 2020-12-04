package discord

import (
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"

	"github.com/Southclaws/cj/types"
)

// Session wraps the discordgo session and provides some additional features
type Session struct {
	S         *discordgo.Session
	Config    types.Config
	UserIndex map[string]discordgo.Member
}

// New creates a new wrapped discord session
func New(s *discordgo.Session, c types.Config) (d *Session) {
	d = &Session{
		S:      s,
		Config: c,
	}
	s.AddHandler(d.ready)
	return
}

// GetUserFromName returns a discordgo.Member from a discord username
func (s *Session) GetUserFromName(name string) (user discordgo.Member, exists bool) {
	user, exists = s.UserIndex[name]
	return
}

func (s *Session) GetRandomChannel() (channel string, err error) {
	channels, err := s.S.GuildChannels(s.Config.GuildID)
	if err != nil {
		return "", err
	}

	active := []string{}
	for _, ch := range channels {
		if ch.ParentID == "375285284079665153" || ch.ParentID == "761517355107876864" {
			active = append(active, ch.ID)
		}
	}

	return active[rand.Intn(len(active))], nil
}

// GetCurrentChannelMessageFrequency returns messages-per-second
func (s *Session) GetCurrentChannelMessageFrequency(channelID string) (freq float64, err error) {
	messages, err := s.S.ChannelMessages(channelID, 20, "", "", "")
	if err != nil {
		return
	}

	start, err := messages[0].Timestamp.Parse()
	if err != nil {
		return
	}
	end, err := messages[len(messages)-1].Timestamp.Parse()
	if err != nil {
		return
	}

	windowSize := float64(start.Unix() - end.Unix())

	freq = float64(len(messages)) / windowSize

	return
}

func (s *Session) ready(session *discordgo.Session, event *discordgo.Ready) {
	s.cacheUsernames()
	c := cron.New()
	must(c.AddFunc("@every 2h", s.cacheUsernames))
	c.Start()
}

func (s *Session) cacheUsernames() {
	users, err := s.S.GuildMembers(s.Config.GuildID, "", 1000)
	if err != nil {
		return
	}

	s.UserIndex = make(map[string]discordgo.Member)
	for _, u := range users {
		name := u.Nick
		if name == "" {
			name = u.User.Username
		}
		s.UserIndex[name] = *u
	}

	return
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
