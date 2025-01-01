package gateway

import (
	"encoding/json"
	"log"

	"presto/internal/constants"
	"presto/internal/handlers/api"
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
	rawResponse, statusCode := api.MakeHTTPRequestToDiscord("/gateway/bot", constants.METHOD_GET, nil)
	json.Unmarshal(rawResponse, &response)

	if statusCode != 200 {
		log.Fatalf("Could not get gateway URL. Expected status code 200 and got %d", statusCode)
	}

	return
}
