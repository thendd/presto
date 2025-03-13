package discord

import (
	"strings"
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
