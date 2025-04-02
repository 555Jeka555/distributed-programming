package event

type Event interface {
	Type() string
}

type SimilarityCalculated struct {
	TextID     string
	Similarity int
}

func (e SimilarityCalculated) Type() string {
	return "valuator.similarity_calculated"
}
