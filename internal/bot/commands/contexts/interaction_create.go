package contexts

import (
	"presto/internal/constants"
	"presto/internal/handlers/api"
	"presto/internal/types"
)

type InteractionCreateContext types.InteractionCreateResponse

type RespondInteractionRequestBody struct {
	Type types.InteractionCallbackType `json:"type"`
	Data any                           `json:"data"`
}

func (ctx InteractionCreateContext) RespondWithMessage(message types.Message) {
	body := RespondInteractionRequestBody{
		Type: constants.INTERACTION_CALLBACK_TYPE_CHANNEL_MESSAGE_WITH_SOURCE,
		Data: message,
	}

	api.MakeHTTPRequestToDiscord("/interactions/"+ctx.ID+"/"+ctx.Token+"/callback", constants.METHOD_POST, body)
}
