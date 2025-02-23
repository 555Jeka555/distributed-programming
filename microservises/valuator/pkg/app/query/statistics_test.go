package query_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"server/pkg/app/model"
	"server/pkg/app/query"
)

type MockTextReadRepository struct {
	texts map[model.TextID]model.Text
}

func NewMockTextReadRepository() *MockTextReadRepository {
	return &MockTextReadRepository{
		texts: make(map[model.TextID]model.Text),
	}
}

func (m *MockTextReadRepository) FindByID(ctx context.Context, textID model.TextID) (model.Text, error) {
	text, exists := m.texts[textID]
	if !exists {
		return model.Text{}, errors.New("text not found")
	}
	return text, nil
}

func TestGetSummary_Success(t *testing.T) {
	repo := NewMockTextReadRepository()
	textStatisticsService := query.NewTextStatisticsQueryService(repo)

	textValue := "Hello, World!"
	textID := model.TextID("1")
	repo.texts[textID] = model.NewText(textID, textValue)

	stats, err := textStatisticsService.GetSummary(context.Background(), string(textID))
	assert.NoError(t, err)
	assert.Equal(t, 10, stats.AlphabetCount)
	assert.Equal(t, 13, stats.AllCount)
}

func TestGetSummary_TextNotFound(t *testing.T) {
	repo := NewMockTextReadRepository()
	textStatisticsService := query.NewTextStatisticsQueryService(repo)

	stats, err := textStatisticsService.GetSummary(context.Background(), "nonexistent_id")
	assert.Error(t, err)
	assert.Equal(t, query.TextStatistics{}, stats)
}
