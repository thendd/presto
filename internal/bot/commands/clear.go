package commands

import (
	"errors"
	"fmt"
	"presto/internal/bot"
	"presto/internal/discord"
)

var Clear = bot.NewSlashCommand("clear", "Deletes messages from a channel", []bot.ApplicationCommandWithHandlerDataOption{
	{
		Type:         discord.APPLICATION_COMMAND_OPTION_TYPE_INTEGER,
		Name:         "amount",
		Description:  "How many messages should be deleted",
		MinimumValue: 2,
		MaximumValue: 100,
		Required:     true,
	},
	{
		Type:        discord.APPLICATION_COMMAND_OPTION_TYPE_STRING,
		Name:        "reason",
		Description: "Why those messages are being deleted",
	},
	{
		Type:        discord.APPLICATION_COMMAND_OPTION_TYPE_CHANNEL,
		Name:        "channel",
		Description: "The channel where the messages will be deleted",
	},
}, ClearHandler).ToApplicationCommand()

type ClearCommandArguments struct {
	Amount    int
	Reason    string
	ChannelID string
}

func ClearHandler(context bot.Context) error {
	arguments := ClearCommandArguments{}

	for _, option := range context.Interaction.Data.Data.Options {
		switch option.Name {
		case "amount":
			arguments.Amount = int(option.Value.(float64))
		case "reason":
			arguments.Reason = option.Value.(string)
		case "channel":
			arguments.ChannelID = option.Value.(string)
		}
	}

	if arguments.ChannelID == "" {
		arguments.ChannelID = context.Interaction.Data.ChannelID
	}

	channel := discord.Channel{
		ID: arguments.ChannelID,
	}

	messages, err := channel.GetMessages(discord.GetChannelMessagesQueryStringParams{
		Limit: arguments.Amount,
	})
	if err != nil {
		return err
	}

	var messagesIDs []string

	for _, message := range messages {
		messagesIDs = append(messagesIDs, message.ID)
	}

	err = channel.BulkDeleteMessages(messagesIDs, arguments.Reason)
	if err != nil {
		return errors.New("You tried to delete messages older than 2 weeks and this is currently impossible")
	}

	successMessage := fmt.Sprintf("**%d** messages were deleted from <#%s> successfully", arguments.Amount, arguments.ChannelID)
	if arguments.Reason != "" {
		successMessage += fmt.Sprintf("with the following reason: **%s**.", arguments.Reason)
	} else {
		successMessage += "."
	}

	return context.Interaction.RespondWithMessage(discord.Message{
		Embeds: []discord.Embed{
			{
				Description: successMessage,
				Color:       discord.EMBED_COLOR_GREEN,
			},
		},
		Flags: discord.MESSAGE_FLAG_EPHEMERAL,
	})
}
