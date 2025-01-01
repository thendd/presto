package types

type Message struct {
	ID              any    `json:"id,omitempty"`
	ChannelID       any    `json:"channel_id,omitempty"`
	Author          *User  `json:"author,omitempty"`
	Content         string `json:"content,omitempty"`
	Timestamp       string `json:"timestamp,omitempty"`
	EditedTimestamp string `json:"edited_timestamp,omitempty"`
	TTS             bool   `json:"tts,omitempty"`
	MentionEveryone bool   `json:"mention_everyone,omitempty"`
	Mentions        []User `json:"mentions,omitempty"`
}
