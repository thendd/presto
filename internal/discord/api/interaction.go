package api

import (
	"errors"
	"net/http"

	"presto/internal/discord"
	"presto/internal/discord/config"
)

const (
	INTERACTION_TYPE_PING discord.InteractionType = iota + 1
	INTERACTION_TYPE_APPLICATION_COMMAND
	INTERACTION_TYPE_MESSAGE_COMPONENT
	INTERACTION_TYPE_APPLICATION_COMMAND_AUTOCOMPLETE
	INTERACTION_TYPE_MODAL_SUBMIT
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
	Data discord.InteractionCreatePayload
}

type RespondInteractionRequestBody struct {
	Type discord.InteractionCallbackType `json:"type"`
	Data any                             `json:"data"`
}

func (ctx Interaction) EditOriginalInteraction(message discord.Message, originalInteractionToken string) error {
	response, statusCode := MakeRequest("/webhooks/"+config.APPLICATION_ID+"/"+originalInteractionToken+"/messages/@original", http.MethodPatch, message)
	if statusCode != http.StatusOK {
		return errors.New(string(response))
	}

	body := RespondInteractionRequestBody{
		Type: INTERACTION_CALLBACK_TYPE_UPDATE_MESSAGE,
	}

	response, statusCode = MakeRequest("/interactions/"+ctx.Data.ID+"/"+ctx.Data.Token+"/callback", http.MethodPost, body)

	if statusCode != http.StatusNoContent {
		return errors.New(string(response))
	}

	return nil
}

func (ctx Interaction) RespondWithMessage(message discord.Message) error {
	body := RespondInteractionRequestBody{
		Type: INTERACTION_CALLBACK_TYPE_CHANNEL_MESSAGE_WITH_SOURCE,
		Data: message,
	}

	response, statusCode := MakeRequest("/interactions/"+ctx.Data.ID+"/"+ctx.Data.Token+"/callback", http.MethodPost, body)
	if statusCode != http.StatusNoContent {
		return errors.New(string(response))
	}

	return nil
}

func (ctx Interaction) RespondWithModal(modal Modal) error {
	body := RespondInteractionRequestBody{
		Type: INTERACTION_CALLBACK_TYPE_MODAL,
		Data: modal,
	}

	response, statusCode := MakeRequest("/interactions/"+ctx.Data.ID+"/"+ctx.Data.Token+"/callback", http.MethodPost, body)
	if statusCode != http.StatusNoContent {
		return errors.New(string(response))
	}

	return nil
}
