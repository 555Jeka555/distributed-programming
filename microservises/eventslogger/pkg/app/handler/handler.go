package handler

import (
	"log"
	"server/pkg/app/event"
)

type Handler interface {
	Handle(evt event.Event) error
}

func NewHandler() Handler {
	return &handler{}
}

type handler struct {
}

func (h *handler) Handle(evt event.Event) error {
	switch e := evt.(type) {
	case event.RankCalculated:
		log.Printf("RankCalculated event - TextID: %s, Rank: %.2f", e.TextID, e.Rank)
	case event.SimilarityCalculated:
		log.Printf("SimilarityCalculated event - TextID: %s, Similarity: %d", e.TextID, e.Similarity)
	default:
		log.Printf("Unknown event type: %s", e.Type())
	}

	return nil
}
