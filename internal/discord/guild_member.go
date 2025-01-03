package discord

type GuildMember struct {
	User         *User  `json:"user,omitempty"`
	Nick         string `json:"nick,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	Roles        []any  `json:"roles"`
	JoinedAt     string `json:"joined_at"`
	PremiumSince string `json:"premium_since,omitempty"`
	Deaf         bool   `json:"deaf"`
	Mute         bool   `json:"mute"`
	Pending      bool   `json:"pending,omitempty"`
	Permissions  string `json:"permissions,omitempty"`
}
