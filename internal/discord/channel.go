package discord

type ChannelType int

const (
	CHANNEL_TYPE_GUILD_TEXT ChannelType = iota
	CHANNEL_TYPE_DM
	CHANNEL_TYPE_GUILD_VOICE
	CHANNEL_TYPE_GROUP_DM
	CHANNEL_TYPE_GUILD_CATEGORY
	CHANNEL_TYPE_GUILD_ANNOUNCEMENT
	CHANNEL_TYPE_ANNOUNCEMENT_THREAD ChannelType = iota + 4
	CHANNEL_TYPE_PUBLIC_THREAD
	CHANNEL_TYPE_PRIVATE_THREAD
	CHANNEL_TYPE_GUILD_STAGE_VOICE
	CHANNEL_TYPE_GUILD_DIRECTORY
	CHANNEL_TYPE_GUILD_FORUM
	CHANNEL_TYPE_GUILD_MEDIA
)

type Channel struct {
	ID         string      `json:"id"`
	Type       ChannelType `json:"type"`
	GuildID    any         `json:"guild_id"`
	Position   int         `json:"position"`
	Name       string      `json:"name"`
	Topic      string      `json:"topic"`
	NSFW       bool        `json:"nsfw"`
	Bitrate    int         `json:"bitrate"`
	UserLimit  int         `json:"user_limit"`
	Recipients []*User     `json:"recipient,omitempty"`
}
