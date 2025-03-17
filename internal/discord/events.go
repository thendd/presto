package discord

const (
	READY              = "READY"
	INTERACTION_CREATE = "INTERACTION_CREATE"
	GUILD_CREATE       = "GUILD_CREATE"
)

const (
	DISPATCH_EVENT_OPCODE = iota
	HEARTBEAT_EVENT_OPCODE
	IDENTIFY_EVENT_OPCODE
	PRESENCE_UPDATE_EVENT_OPCODE
	VOICE_STATE_UPDATE_EVENT_OPCODE
	RESUME_EVENT_OPCODE = iota + 1
	RECONNECT_EVENT_OPCODE
	REQUEST_GUILD_MEMBERS_EVENT_OPCODE
	INVALID_SESSION_EVENT_OPCODE
	HELLO_EVENT_OPCODE
	HEARTBEAT_ACK_EVENT_OPCODE
	REQUEST_SOUNDBOARD_SOUNDS_EVENT_OPCODE
)

type HeartbeatEventData = int

type HelloEventData struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
}

type IdentifyEventData struct {
	Token          string                      `json:"token"`
	Properties     IdentifyEventDataProperties `json:"properties"`
	Compress       bool                        `json:"compress,omitempty"`
	LargeThreshold int                         `json:"large_threshold,omitempty"`
	Shard          *[2]int                     `json:"shard,omitempty"`
	Presence       *UpdatePresence             `json:"presence,omitempty"`
	Intents        int                         `json:"intents"`
}

type IdentifyEventDataProperties struct {
	OS      string `json:"os"`
	Browser string `json:"browser"`
	Device  string `json:"device"`
}
