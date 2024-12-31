package api

import (
	"encoding/json"

	"presto/configs"
)

type sessionStartLimitObject struct {
	Total          int `json:"total"`
	Remaining      int `json:"remaining"`
	ResetAfter     int `json:"reset_after"`
	MaxConcurrency int `json:"max_concurrency"`
}

type getGatewayResponseBody struct {
	URL               string                  `json:"url"`
	Shards            int                     `json:"shards"`
	SessionStartLimit sessionStartLimitObject `json:"session_start_limit"`
}

// Gets the WSS URL so that the bot can listen to interactions, created messages, etc
func GetGateway() (response getGatewayResponseBody) {
	rawResponse := MakeHTTPRequestToDiscord("/gateway/bot", configs.MethodGet, nil)
	json.Unmarshal(rawResponse, &response)

	return
}
