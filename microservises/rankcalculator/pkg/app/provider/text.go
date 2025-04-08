package provider

import "context"

type TextProvider interface {
	GetTextID(text string) string
	GetByTextID(ctx context.Context, textID string) (TextData, error)
}

type TextData struct {
	Similarity int
	Value      string
	Rank       float64
}
