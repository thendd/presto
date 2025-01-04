package discord

type (
	InteractionType         int
	InteractionCallbackType int
)

type InteractionCreatePayload struct {
	ID            string                       `json:"id"`
	ApplicationID any                          `json:"application_id"`
	Type          InteractionType              `json:"type"`
	Data          InteractionCreatePayloadData `json:"data,omitempty"`
	Guild         Guild                        `json:"guild,omitempty"`
	GuildID       string                       `json:"guild_id,omitempty"`
	ChannelID     string                       `json:"channel_id,omitempty"`
	Member        GuildMember                  `json:"member,omitempty"`
	User          User                         `json:"user,omitempty"`
	Token         string                       `json:"token"`
	Version       int                          `json:"version"`
	Message       *Message                     `json:"message,omitempty"`
	Locale        string                       `json:"locale,omitempty"`
	GuildLocale   string                       `json:"guild_locale,omitempty"`
}

type InteractionCreatePayloadData struct {
	ID            any                                  `json:"id"`
	Name          string                               `json:"name"`
	Type          int                                  `json:"type"`
	Resolved      ResolvedData                         `json:"resolved,omitempty"`
	Options       []InteractionCreatePayloadDataOption `json:"options,omitempty"`
	CustomID      string                               `json:"custom_id,omitempty"`
	ComponentType int                                  `json:"component_type,omitempty"`
	TargetID      string                               `json:"target_id,omitempty"`
	Components    []MessageComponent                   `json:"components,omitempty"`
}

type InteractionCreatePayloadDataOption struct {
	Name    string                               `json:"name"`
	Type    int                                  `json:"type"`
	Value   interface{}                          `json:"value,omitempty"`
	Options []InteractionCreatePayloadDataOption `json:"options,omitempty"`
}

type ResolvedData struct {
	Users    map[string]User        `json:"users,omitempty"`
	Members  map[string]GuildMember `json:"members,omitempty"`
	Roles    map[string]Role        `json:"roles,omitempty"`
	Channels map[string]Channel     `json:"channels,omitempty"`
	Messages map[string]Message     `json:"messages,omitempty"`
}
