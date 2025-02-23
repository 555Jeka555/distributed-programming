package model

import (
	"context"
)

type TextID string

type Text struct {
	id    TextID
	value string
}

type TextRepository interface {
	TextReadRepository
	NextID(text string) TextID
	Store(ctx context.Context, text Text) error
}

type TextReadRepository interface {
	FindByID(ctx context.Context, textID TextID) (Text, error)
}

func NewText(id TextID, value string) Text {
	return Text{
		id:    id,
		value: value,
	}
}

func LoadText(id TextID, value string) Text {
	return Text{
		id:    id,
		value: value,
	}
}

func (t *Text) ID() TextID {
	return t.id
}

func (t *Text) Value() string {
	return t.value
}
