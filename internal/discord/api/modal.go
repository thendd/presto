package api

import "presto/internal/discord"

type Modal struct {
	CustomID   string                     `json:"custom_id,omitempty"`
	Title      string                     `json:"title,omitempty"`
	Components []discord.MessageComponent `json:"components,omitempty"`
}
