package discord

type UpdatePresence struct {
	Since  *int      `json:"since"`
	Game   *Activity `json:"activities"`
	Status string    `json:"status"`
	AFK    bool      `json:"afk"`
}

type Activity struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}
