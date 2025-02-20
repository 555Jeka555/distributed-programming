package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"unicode/utf8"
)

var ErrKeyAlreadyExists = errors.New("key already exists")
var ErrKeyNotFound = errors.New("key not found")

const ALPHABET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZабвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"

type Storage interface {
	Set(ctx context.Context, key string, text string) error
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

func (v *valuatorService) AddText(ctx context.Context, text string) int { // TODO без дубликатов и итерирования
	err := v.storage.Set(ctx, hashText(text), "")
	if errors.Is(err, ErrKeyAlreadyExists) {
		return 1
	} else if err != nil {
		log.Panic(err)
	}

	return 0
}

func (v *valuatorService) CalculateRank(text string) float64 { // TODO Измерять кол-во символов без конвертации в []rune
	tmp := text
	alphabetMap := generateAlphabetMap()
	result := make(map[rune]bool)
	nonAlphabet := 0
	allCount := 0
	for len(tmp) > 0 {
		r, size := utf8.DecodeRuneInString(tmp)
		tmp = tmp[size:]
		result[r] = true
		allCount++
		if !alphabetMap[r] {
			nonAlphabet++
		}
	}

	return float64(nonAlphabet) / float64(allCount)
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

func hashText(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
