package discord

type ChannelType int

const (
	CHANNEL_TYPE_GUILD_TEXT          ChannelType = 0
	CHANNEL_TYPE_DM                  ChannelType = 1
	CHANNEL_TYPE_GUILD_VOICE         ChannelType = 2
	CHANNEL_TYPE_GROUP_DM            ChannelType = 3
	CHANNEL_TYPE_GUILD_CATEGORY      ChannelType = 4
	CHANNEL_TYPE_GUILD_ANNOUNCEMENT  ChannelType = 5
	CHANNEL_TYPE_ANNOUNCEMENT_THREAD ChannelType = 10
	CHANNEL_TYPE_PUBLIC_THREAD       ChannelType = 11
	CHANNEL_TYPE_PRIVATE_THREAD      ChannelType = 12
	CHANNEL_TYPE_GUILD_STAGE_VOICE   ChannelType = 13
	CHANNEL_TYPE_GUILD_DIRECTORY     ChannelType = 14
	CHANNEL_TYPE_GUILD_FORUM         ChannelType = 15
	CHANNEL_TYPE_GUILD_MEDIA         ChannelType = 16
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
