package discord

import (
	"encoding/json"
	"net/http"
)

type createDMRequestBody struct {
	RecipientID string `json:"recipient_id"`
}

func CreateDM(recipientId string) Channel {
	response, _ := MakeRequest("/users/@me/channels", http.MethodPost, createDMRequestBody{
		RecipientID: recipientId,
	}, map[string]string{})

	var dmChannel Channel
	json.Unmarshal(response, &dmChannel)

	return dmChannel
}
