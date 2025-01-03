package discord

type Channel struct {
	ID        any    `json:"id"`
	Type      int    `json:"type"`
	GuildID   any    `json:"guild_id"`
	Position  int    `json:"position"`
	Name      string `json:"name"`
	Topic     string `json:"topic"`
	NSFW      bool   `json:"nsfw"`
	Bitrate   int    `json:"bitrate"`
	UserLimit int    `json:"user_limit"`
}
