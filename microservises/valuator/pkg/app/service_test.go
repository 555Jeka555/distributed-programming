package app

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Set(ctx context.Context, key string, text string) error {
	args := m.Called(ctx, key, text)
	return args.Error(0)
}

func (m *MockStorage) ListKey(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockStorage) ListValue(ctx context.Context, keys []string) ([]string, error) {
	args := m.Called(ctx, keys)
	return args.Get(0).([]string), args.Error(1)
}

func TestAddText(t *testing.T) {
	mockStorage := new(MockStorage)
	service := NewValuatorService(mockStorage)

	mockStorage.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	result := service.AddText(context.Background(), "newText")
	assert.Equal(t, 0, result)

	mockStorage.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(ErrKeyAlreadyExists).Once()
	result = service.AddText(context.Background(), "newText")
	assert.Equal(t, 1, result)

	mockStorage.AssertExpectations(t)
}

func TestCalculateRank(t *testing.T) {
	service := NewValuatorService(nil)

	text := "abc 123"
	expectedRank := 0.5714285714285714
	rank := service.CalculateRank(text)
	assert.Equal(t, expectedRank, rank)

	text = "abcd"
	expectedRank = 0
	rank = service.CalculateRank(text)
	assert.Equal(t, expectedRank, rank)

	text = "abc 123 абв"
	expectedRank = 0.45454545454545453
	rank = service.CalculateRank(text)
	assert.Equal(t, expectedRank, rank)
}
