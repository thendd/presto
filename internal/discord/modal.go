package discord

type Modal struct {
	CustomID   string             `json:"custom_id,omitempty"`
	Title      string             `json:"title,omitempty"`
	Components []MessageComponent `json:"components,omitempty"`
}
