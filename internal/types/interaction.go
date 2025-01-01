package types

type InteractionCallbackType int

type InteractionCreateResponse struct {
	ID            string                        `json:"id"`
	ApplicationID any                           `json:"application_id"`
	Type          int                           `json:"type"`
	Data          InteractionCreateResponseData `json:"data,omitempty"`
	GuildID       any                           `json:"guild_id,omitempty"`
	ChannelID     any                           `json:"channel_id,omitempty"`
	Member        GuildMember                   `json:"member,omitempty"`
	User          User                          `json:"user,omitempty"`
	Token         string                        `json:"token"`
	Version       int                           `json:"version"`
	Message       *Message                      `json:"message,omitempty"`
	Locale        string                        `json:"locale,omitempty"`
	GuildLocale   string                        `json:"guild_locale,omitempty"`
}

type InteractionCreateResponseData struct {
	ID            any                                   `json:"id"`
	Name          string                                `json:"name"`
	Type          int                                   `json:"type"`
	Resolved      ResolvedData                          `json:"resolved,omitempty"`
	Options       []InteractionCreateResponseDataOption `json:"options,omitempty"`
	CustomID      string                                `json:"custom_id,omitempty"`
	ComponentType int                                   `json:"component_type,omitempty"`
	TargetID      any                                   `json:"target_id,omitempty"`
}

type InteractionCreateResponseDataOption struct {
	Name    string                                `json:"name"`
	Type    int                                   `json:"type"`
	Value   interface{}                           `json:"value,omitempty"`
	Options []InteractionCreateResponseDataOption `json:"options,omitempty"`
}

type ResolvedData struct {
	Users    map[any]User        `json:"users,omitempty"`
	Members  map[any]GuildMember `json:"members,omitempty"`
	Roles    map[any]Role        `json:"roles,omitempty"`
	Channels map[any]Channel     `json:"channels,omitempty"`
	Messages map[any]Message     `json:"messages,omitempty"`
}
