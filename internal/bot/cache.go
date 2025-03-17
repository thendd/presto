package bot

import (
	"errors"
	"presto/internal/discord"
	"slices"
)

type SelectMenuWithHandler struct {
	Data    discord.MessageComponent
	Handler func(Context, ...any) error
	Args    []any
}

type SelectMenusWithHandler []SelectMenuWithHandler

func (selectMenus *SelectMenusWithHandler) Append(selectMenu SelectMenuWithHandler) {
	*selectMenus = append(*selectMenus, selectMenu)
}

func (selectMenus *SelectMenusWithHandler) Remove(index int) {
	*selectMenus = slices.Delete(*selectMenus, index, index+1)
}

type ModalWithHandler struct {
	Data    discord.Modal
	Handler func(Context, ...any) error
	Args    []any
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
