package discord

import (
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
