package constants

import "presto/internal/types"

const (
	INTERACTION_CALLBACK_TYPE_PONG                                    types.InteractionCallbackType = 1
	INTERACTION_CALLBACK_TYPE_CHANNEL_MESSAGE_WITH_SOURCE             types.InteractionCallbackType = 4
	INTERACTION_CALLBACK_TYPE_DEFERRED_CHANNEL_MESSAGE_WITH_SOURCE    types.InteractionCallbackType = 5
	INTERACTION_CALLBACK_TYPE_DEFERRED_UPDATE_MESSAGE                 types.InteractionCallbackType = 6
	INTERACTION_CALLBACK_TYPE_UPDATE_MESSAGE                          types.InteractionCallbackType = 7
	INTERACTION_CALLBACK_TYPE_APPLICATION_COMMAND_AUTOCOMPLETE_RESULT types.InteractionCallbackType = 8
	INTERACTION_CALLBACK_TYPE_MODAL                                   types.InteractionCallbackType = 9
	INTERACTION_CALLBACK_TYPE_PREMIUM_REQUIRED                        types.InteractionCallbackType = 10
	INTERACTION_CALLBACK_TYPE_LAUNCH_ACTIVITY                         types.InteractionCallbackType = 12
)
