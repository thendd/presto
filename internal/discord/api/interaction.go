package api

import (
	"net/http"

	"presto/internal/discord"

	"github.com/coder/websocket"
)

const (
	INTERACTION_TYPE_PING                             discord.InteractionType = 1
	INTERACTION_TYPE_APPLICATION_COMMAND              discord.InteractionType = 2
	INTERACTION_TYPE_MESSAGE_COMPONENT                discord.InteractionType = 3
	INTERACTION_TYPE_APPLICATION_COMMAND_AUTOCOMPLETE discord.InteractionType = 4
	INTERACTION_TYPE_MODAL_SUBMIT                     discord.InteractionType = 5
)

const (
	INTERACTION_CALLBACK_TYPE_PONG                                    discord.InteractionCallbackType = 1
	INTERACTION_CALLBACK_TYPE_CHANNEL_MESSAGE_WITH_SOURCE             discord.InteractionCallbackType = 4
	INTERACTION_CALLBACK_TYPE_DEFERRED_CHANNEL_MESSAGE_WITH_SOURCE    discord.InteractionCallbackType = 5
	INTERACTION_CALLBACK_TYPE_DEFERRED_UPDATE_MESSAGE                 discord.InteractionCallbackType = 6
	INTERACTION_CALLBACK_TYPE_UPDATE_MESSAGE                          discord.InteractionCallbackType = 7
	INTERACTION_CALLBACK_TYPE_APPLICATION_COMMAND_AUTOCOMPLETE_RESULT discord.InteractionCallbackType = 8
	INTERACTION_CALLBACK_TYPE_MODAL                                   discord.InteractionCallbackType = 9
	INTERACTION_CALLBACK_TYPE_PREMIUM_REQUIRED                        discord.InteractionCallbackType = 10
	INTERACTION_CALLBACK_TYPE_LAUNCH_ACTIVITY                         discord.InteractionCallbackType = 12
)

type Interaction struct {
	Data      discord.InteractionCreatePayload
	Websocket *websocket.Conn
}

type RespondInteractionRequestBody struct {
	Type discord.InteractionCallbackType `json:"type"`
	Data any                             `json:"data"`
}

func (ctx Interaction) RespondWithMessage(message discord.Message) {
	body := RespondInteractionRequestBody{
		Type: INTERACTION_CALLBACK_TYPE_CHANNEL_MESSAGE_WITH_SOURCE,
		Data: message,
	}

	MakeRequest("/interactions/"+ctx.Data.ID+"/"+ctx.Data.Token+"/callback", http.MethodPost, body)
}

func (ctx Interaction) RespondWithModal(modal Modal) {
	body := RespondInteractionRequestBody{
		Type: INTERACTION_CALLBACK_TYPE_MODAL,
		Data: modal,
	}

	MakeRequest("/interactions/"+ctx.Data.ID+"/"+ctx.Data.Token+"/callback", http.MethodPost, body)
}
