package app

import (
	"context"
	"github.com/gofrs/uuid"
	"log"
	"strings"
)

type Storage interface {
	Set(ctx context.Context, key string, text string)
	ListKey(ctx context.Context) ([]string, error)
	ListValue(ctx context.Context, keys []string) ([]string, error)
}

func NewValuatorService(storage Storage) *valuatorService {
	return &valuatorService{
		storage: storage,
	}
}

type ValuatorService interface {
	AddText(ctx context.Context, value string) int
	CalculateRank(text string) float64
}

type valuatorService struct {
	storage Storage
}

func (v *valuatorService) AddText(ctx context.Context, text string) int {
	keys, err := v.storage.ListKey(ctx)
	if err != nil {
		log.Println(err)
	}

	values, err := v.storage.ListValue(ctx, keys)
	if err != nil {
		log.Println(err)
	}

	for _, value := range values {
		if value == text {
			return 1
		}
	}

	uid := uuid.Must(uuid.NewV4())
	v.storage.Set(ctx, uid.String(), text)

	return 0
}

func (v *valuatorService) CalculateRank(text string) float64 {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZабвгдеёжзийклмнопрстуфхцчшщъыьэюя"
	var nonAlphaCount int
	for _, char := range text {
		if !strings.ContainsRune(alphabet, char) {
			nonAlphaCount++
		}
	}
	return float64(nonAlphaCount) / float64(len(text))
}
