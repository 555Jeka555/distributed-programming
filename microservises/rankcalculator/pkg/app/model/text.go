package model

import (
	"context"
)

type TextID string

type Text struct {
	textID     TextID
	similarity int
	value      string
	rank       float64
}

type TextRepository interface {
	NextTextID(text string) TextID
	Store(ctx context.Context, text Text) error
	Delete(ctx context.Context, textID TextID) error
}

func NewText(
	textID TextID,
	similarity int,
	value string,
	rank float64,
) Text {
	return Text{
		textID:     textID,
		similarity: similarity,
		value:      value,
		rank:       rank,
	}
}

func (t *Text) TextID() TextID {
	return t.textID
}

func (t *Text) Value() string {
	return t.value
}

func (t *Text) Similarity() int {
	return t.similarity
}

func (t *Text) Rank() float64 {
	return t.rank
}
