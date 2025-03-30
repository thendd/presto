package bot

import (
	"errors"
	"presto/internal/discord"
	"slices"

	"github.com/google/uuid"
)

type SelectMenusWithHandler []SelectMenuWithHandler

func (selectMenus *SelectMenusWithHandler) Append(selectMenu SelectMenuWithHandler) {
	*selectMenus = append(*selectMenus, selectMenu)
}

func (selectMenus *SelectMenusWithHandler) Remove(index int) {
	*selectMenus = slices.Delete(*selectMenus, index, index+1)
}

type ModalsWithHandler []ModalWithHandler

func (modals *ModalsWithHandler) Append(modal ModalWithHandler) {
	*modals = append(*modals, modal)
}

func (modals *ModalsWithHandler) Remove(index int) {
	*modals = slices.Delete(*modals, index, index+1)
}

type DMChannels []discord.Channel

func (dmChannels *DMChannels) Append(dmChannel discord.Channel) {
	*dmChannels = append(*dmChannels, dmChannel)
}

func (dmChannels *DMChannels) GetByRecipientID(recipientId string) (discord.Channel, error) {
	for _, channel := range *dmChannels {
		for _, recipient := range channel.Recipients {
			if recipient.ID == recipientId {
				return channel, nil
			}
		}
	}

	return discord.Channel{}, errors.New("Not found")
}

type Guilds []discord.Guild

func (guilds *Guilds) Append(guild discord.Guild) {
	*guilds = append(*guilds, guild)
}

func (guilds *Guilds) GetByID(id string) (discord.Guild, error) {
	for _, guild := range *guilds {
		if guild.ID == id {
			return guild, nil
		}
	}

	return discord.Guild{}, errors.New("Not found")
}

type Cache struct {
	DMChannels  DMChannels
	Guilds      Guilds
	SelectMenus SelectMenusWithHandler
	Modals      ModalsWithHandler
}

func NewSelectMenuWithHandler(data discord.MessageComponent, handler func(Context, ...any) error, args []any) SelectMenuWithHandler {
	return SelectMenuWithHandler{
		Data: discord.MessageComponent{
			CustomID:    uuid.NewString(),
			Type:        discord.MESSAGE_COMPONENT_TYPE_SELECT_MENU,
			Title:       data.Title,
			Style:       data.Style,
			Label:       data.Label,
			Placeholder: data.Placeholder,
			MinLength:   data.MinLength,
			MaxLength:   data.MaxLength,
			Required:    data.Required,
			Value:       data.Value,
			Options:     data.Options,
			Components:  data.Components,
		},
		Handler: handler,
		Args:    args,
	}
}

func NewRoleSelectMenuWithHandler(data discord.MessageComponent, handler func(Context, ...any) error, args []any) SelectMenuWithHandler {
	selectMenu := NewSelectMenuWithHandler(data, handler, args)
	selectMenu.Data.Type = discord.MESSAGE_COMPONENT_TYPE_ROLE_SELECT

	return selectMenu
}

func NewModalWithHandler(data discord.Modal, handler func(Context, ...any) error, args []any) ModalWithHandler {
	return ModalWithHandler{
		Data: discord.Modal{
			CustomID:   uuid.NewString(),
			Title:      data.Title,
			Components: data.Components,
		},
		Handler: handler,
		Args:    args,
	}
}
