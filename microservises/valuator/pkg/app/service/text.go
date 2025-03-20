package service

import "server/pkg/app/model"

type TextService struct {
	repo model.TextRepository
}

func NewTextService(repo model.TextRepository) *TextService {
	return &TextService{
		repo: repo,
	}
}

func (s *TextService) GetNextTextID(text string) model.TextID {
	return s.repo.NextTextID(text)
}
