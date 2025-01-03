package discord

const (
	MESSAGE_COMPONENT_TYPE_ACTION_ROW         MessageComponentType = 1
	MESSAGE_COMPONENT_TYPE_BUTTON             MessageComponentType = 2
	MESSAGE_COMPONENT_TYPE_SELECT_MENU        MessageComponentType = 3
	MESSAGE_COMPONENT_TYPE_TEXT_INPUT         MessageComponentType = 4
	MESSAGE_COMPONENT_TYPE_ROLE_SELECT        MessageComponentType = 5
	MESSAGE_COMPONENT_TYPE_USER_SELECT        MessageComponentType = 6
	MESSAGE_COMPONENT_TYPE_MENTIONABLE_SELECT MessageComponentType = 7
	MESSAGE_COMPONENT_TYPE_CHANNEL_SELECT     MessageComponentType = 8
)

const (
	TEXT_INPUT_STYLE_SHORT     MessageComponentStyle = 1
	TEXT_INPUT_STYLE_PARAGRAPH MessageComponentStyle = 2
)

const (
	EMBED_COLOR_GREEN  EmbedColor = 0x95E06C
	EMBED_COLOR_RED    EmbedColor = 0xF0544F
	EMBED_COLOR_YELLOW EmbedColor = 0xF6AE2D
)

type (
	MessageComponentType  int
	MessageComponentStyle int
	EmbedColor            int
)

type Message struct {
	ID              any     `json:"id,omitempty"`
	ChannelID       any     `json:"channel_id,omitempty"`
	Author          *User   `json:"author,omitempty"`
	Content         string  `json:"content,omitempty"`
	Timestamp       string  `json:"timestamp,omitempty"`
	EditedTimestamp string  `json:"edited_timestamp,omitempty"`
	TTS             bool    `json:"tts,omitempty"`
	MentionEveryone bool    `json:"mention_everyone,omitempty"`
	Mentions        []User  `json:"mentions,omitempty"`
	Embeds          []Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       string         `json:"title,omitempty"`
	Type        string         `json:"type,omitempty"`
	Description string         `json:"description,omitempty"`
	URL         string         `json:"url,omitempty"`
	Timestamp   string         `json:"timestamp,omitempty"`
	Color       EmbedColor     `json:"color,omitempty"`
	Footer      *EmbedFooter   `json:"footer,omitempty"`
	Image       *EmbedImage    `json:"image,omitempty"`
	Thumbnail   *EmbedImage    `json:"thumbnail,omitempty"`
	Video       *EmbedVideo    `json:"video,omitempty"`
	Provider    *EmbedProvider `json:"provider,omitempty"`
	Author      *EmbedAuthor   `json:"author,omitempty"`
	Fields      []EmbedField   `json:"fields,omitempty"`
}

type EmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type EmbedImage struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedVideo struct {
	URL    string `json:"url,omitempty"`
	Height int    `json:"height,omitempty"`
	Width  int    `json:"width,omitempty"`
}

type EmbedProvider struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type EmbedAuthor struct {
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type MessageComponent struct {
	Type        MessageComponentType  `json:"type,omitempty"`
	CustomID    string                `json:"custom_id,omitempty"`
	Title       string                `json:"title,omitempty"`
	Style       MessageComponentStyle `json:"style,omitempty"`
	Label       string                `json:"label,omitempty"`
	Placeholder string                `json:"placeholder,omitempty"`
	MinLength   int                   `json:"min_length,omitempty"`
	MaxLength   int                   `json:"max_length,omitempty"`
	Required    bool                  `json:"required,omitempty"`
	Value       string                `json:"value,omitempty"`
	Options     []SelectOption        `json:"options,omitempty"`
	Components  []MessageComponent    `json:"components,omitempty"`
}

type SelectOption struct {
	Label       string `json:"label,omitempty"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
	Emoji       Emoji  `json:"emoji,omitempty"`
	Default     bool   `json:"default,omitempty"`
}

type Emoji struct {
	Name     string `json:"name,omitempty"`
	ID       string `json:"id,omitempty"`
	Animated bool   `json:"animated,omitempty"`
}
