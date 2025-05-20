package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"server/pkg/app/service"
)

func NewHashService() service.HashService {
	return &hashService{}
}

type hashService struct {
}

func (s *hashService) Hash(value string) string {
	hash := sha256.New()
	hash.Write([]byte(value))
	return hex.EncodeToString(hash.Sum(nil))
}
