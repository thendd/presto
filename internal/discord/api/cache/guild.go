package cache

import "presto/internal/discord"

var Guilds = []discord.Guild{}

func GetGuildById(guildId string) discord.Guild {
	for _, guild := range Guilds {
		if guild.ID == guildId {
			return guild
		}
	}

	return discord.Guild{}
}
