package service

import (
	"context"
	"errors"
	"log"
	"server/pkg/app/event"

	"server/pkg/app/model"
)

const ALPHABET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZабвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"

var ErrKeyAlreadyExists = errors.New("key already exists")

func NewRankCalculatorService(
	repo model.TextRepository,
	publisher event.Publisher,
) RankCalculatorService {
	return &rankCalculatorService{
		repo:      repo,
		publisher: publisher,
	}
}

type RankCalculatorService interface {
	AddText(ctx context.Context, value string) error
}

type rankCalculatorService struct {
	repo      model.TextRepository
	publisher event.Publisher
}

func (r *rankCalculatorService) AddText(ctx context.Context, value string) error {
	textID := r.repo.GetTextID(value)
	rank := calcRank(value)
	text := model.NewText(textID, 0, value, rank)

	err := r.repo.Store(ctx, text)
	if err != nil {
		if errors.Is(err, ErrKeyAlreadyExists) {
			err := r.repo.Delete(ctx, textID)
			if err != nil {
				return err
			}
			text = model.NewText(textID, 1, value, rank)

			return r.repo.Store(ctx, text)
		}
		log.Panic(err)
	}

	return r.publisher.PublishInExchange(event.RankCalculated{
		TextID: string(textID),
		Rank:   rank,
	})
}

func calcRank(value string) float64 {
	alphabetCount, allCount := symbolStatistics(value)
	if allCount == 0 {
		return 0
	}
	return 1 - float64(alphabetCount)/float64(allCount)
}

func symbolStatistics(text string) (alphabetCount int, allCount int) {
	alphabetMap := generateAlphabetMap()
	result := make(map[rune]bool)
	for _, r := range text {
		result[r] = true
		allCount++
		if alphabetMap[r] {
			alphabetCount++
		}
	}

	return alphabetCount, allCount
}

func generateAlphabetMap() map[rune]bool {
	result := make(map[rune]bool)
	tmp := ALPHABET

	for _, r := range tmp {
		result[r] = true
	}

	return result
}
