package readmodel

import "github.com/Southclaws/cj/storage"

const (
	defaultUsersPageSize = 50
	maxUsersPageSize     = 200
)

type UserWithPerson struct {
	storage.User
	Person PersonRef `json:"person"`
}

type UsersPage struct {
	Users  []UserWithPerson `json:"users"`
	Total  int              `json:"total"`
	Limit  int              `json:"limit"`
	Offset int              `json:"offset"`
}

type UsersProvider struct {
	storer  storage.Storer
	discord *DiscordProvider
}

func NewUsersProvider(storer storage.Storer, discord *DiscordProvider) *UsersProvider {
	return &UsersProvider{storer: storer, discord: discord}
}

func (p *UsersProvider) List(query string, limit, offset int) (UsersPage, error) {
	if limit <= 0 || limit > maxUsersPageSize {
		limit = defaultUsersPageSize
	}
	if offset < 0 {
		offset = 0
	}

	users, total, err := p.storer.ListUsers(query, limit, offset)
	if err != nil {
		return UsersPage{}, err
	}

	out := make([]UserWithPerson, len(users))
	seen := make(map[string]bool, len(users))
	for i := range users {
		normalizeUser(&users[i])
		out[i] = UserWithPerson{User: users[i], Person: p.discord.PersonRef(users[i].DiscordUserID)}
		seen[users[i].DiscordUserID] = true
	}

	if query != "" && offset == 0 {
		for _, person := range p.discord.MatchMembers(query, maxUsersPageSize) {
			if seen[person.ID] || len(out) >= limit {
				break
			}
			seen[person.ID] = true
			out = append(out, UserWithPerson{
				User:   storage.User{DiscordUserID: person.ID, ReceivedReactions: []storage.ReactionCounter{}},
				Person: person,
			})
			total++
		}
	}

	return UsersPage{Users: out, Total: total, Limit: limit, Offset: offset}, nil
}

func (p *UsersProvider) Get(discordUserID string) (storage.User, bool, error) {
	user, found, err := p.storer.GetUser(discordUserID)
	if err != nil {
		return user, found, err
	}
	normalizeUser(&user)
	return user, found, nil
}

func normalizeUser(u *storage.User) {
	if u.ReceivedReactions == nil {
		u.ReceivedReactions = []storage.ReactionCounter{}
	}
}
