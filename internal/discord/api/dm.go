package api

import (
	"encoding/json"
	"net/http"

	"presto/internal/discord"
	"presto/internal/discord/api/cache"
)

type createDMRequestBody struct {
	RecipientID string `json:"recipient_id"`
}

func CreateDM(recipientId string) discord.Channel {
	response, _ := MakeRequest("/users/@me/channels", http.MethodPost, createDMRequestBody{
		RecipientID: recipientId,
	})

	var dmChannel discord.Channel
	json.Unmarshal(response, &dmChannel)

	cache.DMChannels = append(cache.DMChannels, dmChannel)

	return dmChannel
}
