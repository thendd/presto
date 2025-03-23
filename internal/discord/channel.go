package discord

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

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

type GetChannelMessagesQueryStringParams struct {
	Around string
	Before string
	After  string
	Limit  int
}

type BulkDeleteMessagesFromChannelRequestBody struct {
	Messages []string `json:"messages"`
}

func (channel *Channel) GetMessages(params GetChannelMessagesQueryStringParams) ([]Message, error) {
	endpoint := "/channels/" + channel.ID + "/messages?"

	if params.Around != "" {
		endpoint += "around=" + params.Around + "&"
	}

	if params.Before != "" {
		endpoint += "before=" + params.Before + "&"
	}

	if params.After != "" {
		endpoint += "after=" + params.After + "&"
	}

	if params.Limit != 0 {
		endpoint += "limit=" + strconv.Itoa(params.Limit)
	}

	response, statusCode := MakeRequest(endpoint, http.MethodGet, nil, map[string]string{})
	if statusCode != 200 {
		return []Message{}, errors.New(string(response))
	}

	var messages []Message
	json.Unmarshal(response, &messages)

	return messages, nil
}

func (channel *Channel) BulkDeleteMessages(IDs []string, reason string) error {
	_, statusCode := MakeRequest("/channels/"+channel.ID+"/messages/bulk-delete", http.MethodPost, BulkDeleteMessagesFromChannelRequestBody{
		Messages: IDs,
	}, map[string]string{
		"X-Audit-Log-Reason": reason,
	})
	if statusCode != 204 {
		return errors.New(fmt.Sprintf("Expected status code 204 but got %d instead", statusCode))
	}

	return nil
}
