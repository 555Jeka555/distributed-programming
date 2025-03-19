package model

import (
	"context"
)

type TextID string
type RankID string

type Text struct {
	textID TextID
	rankID RankID
	value  string
	rank   float64
}

type TextRepository interface {
	TextReadRepository
	NextTextID(text string) TextID
	NextRankID(text string) RankID
	Store(ctx context.Context, text Text) error
}

type TextReadRepository interface {
	FindByID(ctx context.Context, textID TextID, rankID RankID) (Text, error)
}

func NewText(
	textID TextID,
	rankID RankID,
	value string,
	rank float64,
) Text {
	return Text{
		textID: textID,
		rankID: rankID,
		value:  value,
		rank:   rank,
	}
}

func LoadText(
	textID TextID,
	rankID RankID,
	value string,
	rank float64,
) Text {
	return Text{
		textID: textID,
		rankID: rankID,
		value:  value,
		rank:   rank,
	}
}

func (t *Text) TextID() TextID {
	return t.textID
}

func (t *Text) RankID() RankID {
	return t.rankID
}

func (t *Text) Value() string {
	return t.value
}

func (t *Text) Rank() float64 {
	return t.rank
}
