package types

type ApplicationCommandType int

type ApplicationCommandOptionType int

type ApplicationCommandOption struct {
	Type         ApplicationCommandOptionType `json:"type"`
	Name         string                       `json:"name"`
	Description  string                       `json:"description"`
	Required     bool                         `json:"required"`
	Autocomplete bool                         `json:"autocomplete"`
}
