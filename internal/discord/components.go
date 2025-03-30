package discord

import "github.com/google/uuid"

func NewTextInputComponent(data MessageComponent) MessageComponent {
	return MessageComponent{
		Type:        MESSAGE_COMPONENT_TYPE_TEXT_INPUT,
		CustomID:    uuid.NewString(),
		Style:       data.Style,
		Label:       data.Label,
		MinLength:   data.MinLength,
		MaxLength:   data.MaxLength,
		Required:    data.Required,
		Value:       data.Value,
		Placeholder: data.Placeholder,
	}
}
