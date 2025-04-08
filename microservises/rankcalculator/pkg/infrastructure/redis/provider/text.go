package provider

import (
	"crypto/sha256"
	"encoding/hex"
	"server/pkg/app/provider"
)

func NewTextProvider() provider.TextProvider {
	return &textProvider{}
}

type textProvider struct {
}

func (p *textProvider) GetTextID(text string) string {
	return hashText(text)
}

func hashText(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
