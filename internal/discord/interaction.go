package discord

import (
	"errors"
	"net/http"
	"presto/internal/config"
	"strings"
)

const (
	INTERACTION_TYPE_PING InteractionType = iota + 1
	INTERACTION_TYPE_APPLICATION_COMMAND
	INTERACTION_TYPE_MESSAGE_COMPONENT
	INTERACTION_TYPE_APPLICATION_COMMAND_AUTOCOMPLETE
	INTERACTION_TYPE_MODAL_SUBMIT
)

const (
	INTERACTION_CALLBACK_TYPE_PONG                        InteractionCallbackType = iota + 1
	INTERACTION_CALLBACK_TYPE_CHANNEL_MESSAGE_WITH_SOURCE InteractionCallbackType = iota + 3
	INTERACTION_CALLBACK_TYPE_DEFERRED_CHANNEL_MESSAGE_WITH_SOURCE
	INTERACTION_CALLBACK_TYPE_DEFERRED_UPDATE_MESSAGE
	INTERACTION_CALLBACK_TYPE_UPDATE_MESSAGE
	INTERACTION_CALLBACK_TYPE_APPLICATION_COMMAND_AUTOCOMPLETE_RESULT
	INTERACTION_CALLBACK_TYPE_MODAL
	INTERACTION_CALLBACK_TYPE_PREMIUM_REQUIRED
	INTERACTION_CALLBACK_TYPE_LAUNCH_ACTIVITY InteractionCallbackType = iota + 4
)

type (
	InteractionType         int
	InteractionCallbackType int
)

type InteractionCreatePayload struct {
	ID            string                       `json:"id"`
	ApplicationID any                          `json:"application_id"`
	Type          InteractionType              `json:"type"`
	Data          InteractionCreatePayloadData `json:"data,omitempty"`
	Guild         Guild                        `json:"guild,omitempty"`
	GuildID       string                       `json:"guild_id,omitempty"`
	ChannelID     string                       `json:"channel_id,omitempty"`
	Member        GuildMember                  `json:"member,omitempty"`
	User          User                         `json:"user,omitempty"`
	Token         string                       `json:"token"`
	Version       int                          `json:"version"`
	Message       *Message                     `json:"message,omitempty"`
	Locale        string                       `json:"locale,omitempty"`
	GuildLocale   string                       `json:"guild_locale,omitempty"`
}

type InteractionCreatePayloadData struct {
	ID            any                                  `json:"id"`
	Name          string                               `json:"name"`
	Type          ApplicationCommandType               `json:"type"`
	Resolved      ResolvedData                         `json:"resolved,omitempty"`
	Options       []InteractionCreatePayloadDataOption `json:"options,omitempty"`
	CustomID      string                               `json:"custom_id,omitempty"`
	ComponentType int                                  `json:"component_type,omitempty"`
	TargetID      string                               `json:"target_id,omitempty"`
	Components    []MessageComponent                   `json:"components,omitempty"`
	Values        []string                             `json:"values,omitempty"`
}

type InteractionCreatePayloadDataOption struct {
	Name    string                               `json:"name"`
	Type    int                                  `json:"type"`
	Value   any                                  `json:"value,omitempty"`
	Options []InteractionCreatePayloadDataOption `json:"options,omitempty"`
}

type ResolvedData struct {
	Users    map[string]User        `json:"users,omitempty"`
	Members  map[string]GuildMember `json:"members,omitempty"`
	Roles    map[string]Role        `json:"roles,omitempty"`
	Channels map[string]Channel     `json:"channels,omitempty"`
	Messages map[string]Message     `json:"messages,omitempty"`
}

type Interaction struct {
	Data InteractionCreatePayload
}

type RespondInteractionRequestBody struct {
	Type InteractionCallbackType `json:"type"`
	Data any                     `json:"data"`
}

func (ctx Interaction) EditOriginalInteraction(message Message, originalInteractionToken string) error {
	response, statusCode := MakeRequest("/webhooks/"+config.DISCORD_APPLICATION_ID+"/"+originalInteractionToken+"/messages/@original", http.MethodPatch, message)
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

func (ctx Interaction) RespondWithMessage(message Message) error {
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

// Joins the name of the interaction options if they are sub commands or sub command groups.
// This is used for comparision purposes with names of application commands
func JoinInteractionOptionsNames(options []InteractionCreatePayloadDataOption) string {
	var names []string

	var appendNames func([]InteractionCreatePayloadDataOption)
	appendNames = func(opts []InteractionCreatePayloadDataOption) {
		for _, option := range opts {
			if option.Type == int(APPLICATION_COMMAND_OPTION_TYPE_SUB_COMMAND) || option.Type == int(APPLICATION_COMMAND_OPTION_TYPE_SUB_COMMAND_GROUP) {
				names = append(names, option.Name)
				appendNames(option.Options)
			}
		}
	}

	appendNames(options)

	return strings.Join(names, " ")
}

// Gets the "whole name" of an interaction, joining its base name, sub command names and sub command group names
func GetInteractionName(interaction InteractionCreatePayloadData) string {
	optionsNames := JoinInteractionOptionsNames(interaction.Options)

	if len(optionsNames) == 0 {
		return interaction.Name
	}

	fullName := []string{interaction.Name, optionsNames}
	return strings.Join(fullName, " ")
}
