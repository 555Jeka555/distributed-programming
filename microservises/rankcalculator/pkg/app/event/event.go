package event

type Event interface {
	Type() string
}

type RankCalculated struct {
	TextID string
	Rank   float64
}

func (e RankCalculated) Type() string {
	return "rankcalculator.rank_calculated"
}

type TextSubmitted struct {
	TextValue string
}

func (e TextSubmitted) Type() string {
	return "valuator.text_submitted"
}
