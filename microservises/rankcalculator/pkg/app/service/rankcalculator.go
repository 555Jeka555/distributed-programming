package service

import (
	"context"
	"errors"
	"log"

	"server/pkg/app/model"
)

const ALPHABET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZабвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"

var ErrKeyAlreadyExists = errors.New("key already exists")

func NewRankCalculatorService(repo model.TextRepository) RankCalculatorService {
	return &rankCalculatorService{
		repo: repo,
	}
}

type RankCalculatorService interface {
	AddText(ctx context.Context, value string) (model.TextID, model.RankID, error)
}

type rankCalculatorService struct {
	repo model.TextRepository
}

func (v *rankCalculatorService) AddText(ctx context.Context, value string) (model.TextID, model.RankID, error) { // TODO без дубликатов и итерирования
	textID := v.repo.NextTextID(value)
	rankID := v.repo.NextRankID(value)
	alphabetCount, allCount := symbolStatistics(value)
	rank := float64(alphabetCount) / float64(allCount) // TODO Хранить сразу rank
	text := model.NewText(textID, rankID, value, rank)

	err := v.repo.Store(ctx, text)
	if err != nil {
		if errors.Is(err, ErrKeyAlreadyExists) {
			return textID, rankID, err
		}
		log.Panic(err)
	}

	return textID, rankID, nil
}

func symbolStatistics(text string) (alphabetCount int, allCount int) {
	alphabetMap := generateAlphabetMap()
	result := make(map[rune]bool)
	// TODO for range по строке
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
