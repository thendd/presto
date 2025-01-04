package api

import (
	"net/http"

	"presto/internal/discord"
)

func SendMessage(message discord.Message) (response []byte, statusCode int) {
	response, statusCode = MakeRequest("/channels/"+message.ChannelID+"/messages", http.MethodPost, message)
	return
}
