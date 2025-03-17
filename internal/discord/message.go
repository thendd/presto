package discord

import (
	"errors"
	"net/http"
	"presto/internal/log"
)

const (
	MESSAGE_COMPONENT_TYPE_ACTION_ROW MessageComponentType = iota + 1
	MESSAGE_COMPONENT_TYPE_BUTTON
	MESSAGE_COMPONENT_TYPE_SELECT_MENU
	MESSAGE_COMPONENT_TYPE_TEXT_INPUT
	MESSAGE_COMPONENT_TYPE_USER_SELECT
	MESSAGE_COMPONENT_TYPE_ROLE_SELECT
	MESSAGE_COMPONENT_TYPE_MENTIONABLE_SELECT
	MESSAGE_COMPONENT_TYPE_CHANNEL_SELECT
)

const (
	TEXT_INPUT_STYLE_SHORT MessageComponentStyle = iota + 1
	TEXT_INPUT_STYLE_PARAGRAPH
	BUTTON_STYLE_PRIMARY
	BUTTON_STYLE_SECONDARY
	BUTTON_STYLE_SUCCESS
	BUTTON_STYLE_DANGER
	BUTTON_STYLE_LINK
	BUTTON_STYLE_PREMIUM
)

const (
	EMBED_COLOR_GREEN  EmbedColor = 0x95E06C
	EMBED_COLOR_RED    EmbedColor = 0xF0544F
	EMBED_COLOR_YELLOW EmbedColor = 0xF6AE2D
)

const (
	MESSAGE_FLAG_CROSSPOSTED                            MessageFlag = 1              // 1 << 0
	MESSAGE_FLAG_IS_CROSSPOST                           MessageFlag = 2 ^ (iota - 1) // 1 << 1
	MESSAGE_FLAG_SUPPRESS_EMBEDS                        MessageFlag = 4              // 1 << 2
	MESSAGE_FLAG_SOURCE_MESSAGE_DELETED                 MessageFlag = 8              // 1 << 3
	MESSAGE_FLAG_URGENT                                 MessageFlag = 16             // 1 << 4
	MESSAGE_FLAG_HAS_THREAD                             MessageFlag = 32             // 1 << 5
	MESSAGE_FLAG_EPHEMERAL                              MessageFlag = 64             // 1 << 6
	MESSAGE_FLAG_LOADING                                MessageFlag = 128            // 1 << 7
	MESSAGE_FLAG_FAILED_TO_MENTION_SOME_ROLES_IN_THREAD MessageFlag = 256            // 1 << 8
	MESSAGE_FLAG_SUPPRESS_NOTIFICATIONS                 MessageFlag = 4096           // 1 << 12
	MESSAGE_FLAG_IS_VOICE_MESSAGE                       MessageFlag = 8192           // 1 << 13
)

type (
	MessageComponentType  int
	MessageComponentStyle int
	EmbedColor            int
	MessageFlag           int
)

type Message struct {
	ID              string             `json:"id,omitempty"`
	ChannelID       string             `json:"channel_id,omitempty"`
	Author          *User              `json:"author,omitempty"`
	Content         string             `json:"content,omitempty"`
	Timestamp       string             `json:"timestamp,omitempty"`
	EditedTimestamp string             `json:"edited_timestamp,omitempty"`
	TTS             bool               `json:"tts,omitempty"`
	MentionEveryone bool               `json:"mention_everyone,omitempty"`
	Mentions        []User             `json:"mentions,omitempty"`
	Embeds          []Embed            `json:"embeds,omitempty"`
	Flags           MessageFlag        `json:"flags,omitempty"`
	Components      []MessageComponent `json:"components"`
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
	Required    bool                  `json:"required"`
	Value       string                `json:"value,omitempty"`
	Options     []SelectOption        `json:"options,omitempty"`
	Components  []MessageComponent    `json:"components,omitempty"`
}

type SelectOption struct {
	Label       string `json:"label,omitempty"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
	Emoji       *Emoji `json:"emoji,omitempty"`
	Default     bool   `json:"default,omitempty"`
}

type Emoji struct {
	Name     string `json:"name,omitempty"`
	ID       string `json:"id,omitempty"`
	Animated bool   `json:"animated,omitempty"`
}

func (message *Message) Send() error {
	response, statusCode := MakeRequest("/channels/"+message.ChannelID+"/messages", http.MethodPost, message)
	if statusCode != http.StatusCreated {
		log.Errorf("Could not send message: expected status code 201 but received %d. The API response was:\n%s", statusCode, string(response))
		return errors.New(string(response))
	}

	return nil
}
