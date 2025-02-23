package query

import (
	"context"
	"server/pkg/app/model"
	"unicode/utf8"
)

const ALPHABET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZабвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"

type TextStatisticsQueryService interface {
	GetSummary(ctx context.Context, textID string) (TextStatistics, error)
}

func NewTextStatisticsQueryService(textReadRepository model.TextReadRepository) TextStatisticsQueryService {
	return &textStatisticsQueryService{
		textReadRepository: textReadRepository,
	}
}

type textStatisticsQueryService struct {
	textReadRepository model.TextReadRepository
}

type TextStatistics struct {
	AlphabetCount int
	AllCount      int
}

func (s *textStatisticsQueryService) GetSummary(ctx context.Context, textID string) (TextStatistics, error) {
	text, err := s.textReadRepository.FindByID(ctx, model.TextID(textID))
	if err != nil {
		return TextStatistics{}, err
	}

	alphabetCount, allCount := symbolStatistics(text)

	return TextStatistics{
		AlphabetCount: alphabetCount,
		AllCount:      allCount,
	}, nil
}

func symbolStatistics(text model.Text) (alphabetCount int, allCount int) {
	tmp := text.Value()
	alphabetMap := generateAlphabetMap()
	result := make(map[rune]bool)
	for len(tmp) > 0 {
		r, size := utf8.DecodeRuneInString(tmp)
		tmp = tmp[size:]
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
	for len(tmp) > 0 {
		r, size := utf8.DecodeRuneInString(tmp)
		tmp = tmp[size:]
		result[r] = true
	}
	return result
}
