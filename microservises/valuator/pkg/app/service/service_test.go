package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"server/pkg/app/model"
	"server/pkg/app/service"
)

type MockTextRepository struct {
	texts map[model.TextID]model.Text
}

func NewMockTextRepository() *MockTextRepository {
	return &MockTextRepository{
		texts: make(map[model.TextID]model.Text),
	}
}

func (m *MockTextRepository) NextID(text string) model.TextID {
	id := model.TextID(text)
	return id
}

func (m *MockTextRepository) Store(_ context.Context, text model.Text) error {
	if _, exists := m.texts[text.ID()]; exists {
		return service.ErrKeyAlreadyExists
	}
	m.texts[text.ID()] = text
	return nil
}

func (m *MockTextRepository) FindByID(ctx context.Context, textID model.TextID) (model.Text, error) {
	text, exists := m.texts[textID]
	if !exists {
		return model.Text{}, errors.New("text not found")
	}
	return text, nil
}

func TestAddText_Success(t *testing.T) {
	repo := NewMockTextRepository()
	valuatorService := service.NewValuatorService(repo)

	textValue := "Hello, World!"
	textID, err := valuatorService.AddText(context.Background(), textValue)

	assert.NoError(t, err)
	assert.Equal(t, model.TextID("Hello, World!"), textID)
	assert.Len(t, repo.texts, 1)
}

func TestAddText_Duplicate(t *testing.T) {
	repo := NewMockTextRepository()
	valuatorService := service.NewValuatorService(repo)

	textValue := "Hello, World!"
	_, err := valuatorService.AddText(context.Background(), textValue)
	assert.NoError(t, err)

	textID, err := valuatorService.AddText(context.Background(), textValue)
	assert.True(t, errors.Is(err, service.ErrKeyAlreadyExists))

	assert.Equal(t, model.TextID("Hello, World!"), textID)
}
