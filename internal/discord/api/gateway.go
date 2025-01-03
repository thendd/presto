package api

import (
	"encoding/json"
	"log"
	"net/http"
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
	log.Println("Started fetching Discord's websocket URL")

	rawResponse, statusCode := MakeRequest("/gateway/bot", http.MethodGet, nil)
	json.Unmarshal(rawResponse, &response)

	if statusCode != http.StatusOK {
		log.Fatalf("Could not get gateway URL. Expected status code 200 and got %d", statusCode)
	}

	log.Println("Finished fetching Discord's WSS URL successfully")

	return
}
