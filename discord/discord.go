package discord

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
	"go.uber.org/zap"

	"github.com/Southclaws/cj/types"
)

const memberPageSize = 1000

type Session struct {
	S      *discordgo.Session
	Config types.Config

	mu               sync.RWMutex
	userIndex        map[string]discordgo.Member
	membersByID      map[string]*discordgo.Member
	members          []*discordgo.Member
	membersUpdatedAt time.Time
}

func New(s *discordgo.Session, c types.Config) (d *Session) {
	d = &Session{
		S:      s,
		Config: c,
	}
	s.AddHandler(d.ready)
	return
}

func (s *Session) GetUserFromName(name string) (user discordgo.Member, exists bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, exists = s.userIndex[name]
	return
}

func (s *Session) Members() ([]*discordgo.Member, time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*discordgo.Member, len(s.members))
	copy(out, s.members)
	return out, s.membersUpdatedAt
}

func (s *Session) MemberByID(id string) (discordgo.Member, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.membersByID[id]
	if !ok {
		return discordgo.Member{}, false
	}
	return *m, true
}

func (s *Session) RefreshMembers() error {
	members, err := paginateAllMembers(func(after string) ([]*discordgo.Member, error) {
		return s.S.GuildMembers(s.Config.GuildID, after, memberPageSize)
	})
	if err != nil {
		return err
	}

	index := make(map[string]discordgo.Member, len(members))
	byID := make(map[string]*discordgo.Member, len(members))
	for _, m := range members {
		m.GuildID = s.Config.GuildID

		name := m.Nick
		if name == "" && m.User != nil {
			name = m.User.Username
		}
		index[name] = *m
		if m.User != nil {
			byID[m.User.ID] = m
		}
	}

	s.mu.Lock()
	s.userIndex = index
	s.membersByID = byID
	s.members = members
	s.membersUpdatedAt = time.Now()
	s.mu.Unlock()

	return nil
}

func paginateAllMembers(fetch func(after string) ([]*discordgo.Member, error)) ([]*discordgo.Member, error) {
	var all []*discordgo.Member
	after := ""
	for {
		batch, err := fetch(after)
		if err != nil {
			return nil, err
		}
		all = append(all, batch...)
		if len(batch) < memberPageSize {
			return all, nil
		}
		after = batch[len(batch)-1].User.ID
	}
}

func (s *Session) GetCurrentChannelMessageFrequency(channelID string) (freq float64, err error) {
	messages, err := s.S.ChannelMessages(channelID, 20, "", "", "")
	if err != nil {
		return
	}

	if len(messages) == 0 {
		return 0.0, nil
	}

	start, err := time.Parse(time.RFC3339, messages[0].Timestamp.Format(time.RFC3339))
	if err != nil {
		return
	}
	end, err := time.Parse(time.RFC3339, messages[len(messages)-1].Timestamp.Format(time.RFC3339))
	if err != nil {
		return
	}

	windowSize := float64(start.Unix() - end.Unix())

	freq = float64(len(messages)) / windowSize

	return
}

func (s *Session) ready(session *discordgo.Session, event *discordgo.Ready) {
	go func() {
		if err := s.RefreshMembers(); err != nil {
			zap.L().Error("failed to build initial member cache", zap.Error(err))
		}
	}()

	c := cron.New()
	must(c.AddFunc("@every 2h", func() {
		if err := s.RefreshMembers(); err != nil {
			zap.L().Error("failed to refresh member cache", zap.Error(err))
		}
	}))
	c.Start()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
