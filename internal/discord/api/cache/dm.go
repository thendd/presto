package cache

import "presto/internal/discord"

var DMChannels = []discord.Channel{}

func GetDMChannelByRecipientID(recipientId string) discord.Channel {
	for _, channel := range DMChannels {
		for _, recipient := range channel.Recipients {
			if recipient.ID == recipientId {
				return channel
			}
		}
	}

	return discord.Channel{}
}
