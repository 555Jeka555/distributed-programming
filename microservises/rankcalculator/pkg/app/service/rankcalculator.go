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

func (v *rankCalculatorService) AddText(ctx context.Context, value string) error {
	textID := v.repo.GetTextID(value)
	alphabetCount, allCount := symbolStatistics(value)
	rank := 1 - float64(alphabetCount)/float64(allCount)
	text := model.NewText(textID, 0, value, rank)

	err := v.repo.Store(ctx, text)
	if err != nil {
		if errors.Is(err, ErrKeyAlreadyExists) {
			err := v.repo.Delete(ctx, textID)
			if err != nil {
				return err
			}
			text = model.NewText(textID, 1, value, rank)

			return v.repo.Store(ctx, text)
		}
		log.Panic(err)
	}

	return v.publisher.PublishInExchange(event.RankCalculated{
		TextID: string(textID),
		Rank:   rank,
	})
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
