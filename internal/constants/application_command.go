package constants

import "presto/internal/types"

const (
	APPLICATION_COMMAND_TYPE_CHAT_INPUT          types.ApplicationCommandType = 1
	APPLICATION_COMMAND_TYPE_USER                types.ApplicationCommandType = 2
	APPLICATION_COMMAND_TYPE_MESSAGE             types.ApplicationCommandType = 3
	APPLICATION_COMMAND_TYPE_PRIMARY_ENTRY_POINT types.ApplicationCommandType = 4
)

const (
	APPLICATION_COMMAND_OPTION_TYPE_SUB_COMMAND       types.ApplicationCommandOptionType = 1
	APPLICATION_COMMAND_OPTION_TYPE_SUB_COMMAND_GROUP types.ApplicationCommandOptionType = 2
	APPLICATION_COMMAND_OPTION_TYPE_STRING            types.ApplicationCommandOptionType = 3
	APPLICATION_COMMAND_OPTION_TYPE_INTEGER           types.ApplicationCommandOptionType = 4
	APPLICATION_COMMAND_OPTION_TYPE_BOOLEAN           types.ApplicationCommandOptionType = 5
	APPLICATION_COMMAND_OPTION_TYPE_USER              types.ApplicationCommandOptionType = 6
	APPLICATION_COMMAND_OPTION_TYPE_CHANNEL           types.ApplicationCommandOptionType = 7
	APPLICATION_COMMAND_OPTION_TYPE_ROLE              types.ApplicationCommandOptionType = 8
	APPLICATION_COMMAND_OPTION_TYPE_MENTIONABLE       types.ApplicationCommandOptionType = 9
	APPLICATION_COMMAND_OPTION_TYPE_NUMBER            types.ApplicationCommandOptionType = 10
	APPLICATION_COMMAND_OPTION_TYPE_ATTACHMENT        types.ApplicationCommandOptionType = 11
)
