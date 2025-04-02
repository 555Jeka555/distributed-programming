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

type SimilarityCalculated struct {
	TextID     string
	Similarity int
}

func (e SimilarityCalculated) Type() string {
	return "valuator.similarity_calculated"
}
