package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"server/pkg/app/event"
	"server/pkg/app/model"
)

type MockTextRepository struct {
	mock.Mock
}

func (m *MockTextRepository) Store(ctx context.Context, text model.Text) error {
	args := m.Called(ctx, text)
	return args.Error(0)
}

func (m *MockTextRepository) GetTextID(value string) model.TextID {
	args := m.Called(value)
	if v, ok := args.Get(0).(model.TextID); ok {
		return v
	}
	return model.TextID(args.String(0))
}

func (m *MockTextRepository) Delete(ctx context.Context, textID model.TextID) error {
	args := m.Called(ctx, textID)
	return args.Error(0)
}

func (m *MockTextRepository) FindByID(ctx context.Context, textID model.TextID) (model.Text, error) {
	args := m.Called(ctx, textID)
	var tеxt model.Text
	if args.Get(0) != nil {
		tеxt = args.Get(0).(model.Text)
	}
	return tеxt, args.Error(1)
}

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) PublishInExchange(ev event.Event) error {
	args := m.Called(ev)
	return args.Error(0)
}

func TestRankCalculatorService_AddText(t *testing.T) {
	type expectedErrors struct {
		storeErr   error
		deleteErr  error
		publishErr error
	}
	type args struct {
		value string
	}
	tests := []struct {
		name           string
		expectedErrors expectedErrors
		args           args
		wantRank       float64
		wantErr        bool
		wantDelete     bool
		wantStoreCnt   int
		wantPublish    bool
	}{
		{
			name:           "Обычный кейс - только буквы",
			expectedErrors: expectedErrors{},
			args:           args{value: "abcXYZ"},
			wantRank:       0,
			wantErr:        false,
			wantDelete:     false,
			wantStoreCnt:   1,
			wantPublish:    true,
		},
		{
			name:           "Кейс с Emoji",
			expectedErrors: expectedErrors{},
			args:           args{value: "a😀b"},
			wantRank:       1 - 2.0/3.0,
			wantErr:        false,
			wantDelete:     false,
			wantStoreCnt:   1,
			wantPublish:    true,
		},
		{
			name:           "Кейс: ключ уже есть, удаляем и заново сохраняем",
			expectedErrors: expectedErrors{storeErr: ErrKeyAlreadyExists},
			args:           args{value: "abc"},
			wantRank:       0,
			wantErr:        false,
			wantDelete:     true,
			wantStoreCnt:   2,
			wantPublish:    false,
		},
		{
			name:           "Кейс: ошибка при удалении",
			expectedErrors: expectedErrors{storeErr: ErrKeyAlreadyExists, deleteErr: errors.New("del error")},
			args:           args{value: "abc"},
			wantRank:       0,
			wantErr:        true,
			wantDelete:     true,
			wantStoreCnt:   1,
			wantPublish:    false,
		},
		{
			name:           "Кейс: ошибка при публикации",
			expectedErrors: expectedErrors{publishErr: errors.New("pub error")},
			args:           args{value: "abc"},
			wantRank:       0,
			wantErr:        true,
			wantDelete:     false,
			wantStoreCnt:   1,
			wantPublish:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTextRepository)
			mockPublisher := new(MockPublisher)
			textID := "id_" + tt.args.value

			mockRepo.On("GetTextID", tt.args.value).Return(model.TextID(textID))

			if tt.wantStoreCnt == 2 {
				mockRepo.On("Store", mock.Anything, mock.AnythingOfType("model.Text")).Return(ErrKeyAlreadyExists).Once()
				mockRepo.On("Delete", mock.Anything, model.TextID(textID)).Return(tt.expectedErrors.deleteErr).Once()
				if tt.expectedErrors.deleteErr == nil {
					mockRepo.On("Store", mock.Anything, mock.AnythingOfType("model.Text")).Return(nil).Once()
				}
			} else {
				mockRepo.On("Store", mock.Anything, mock.AnythingOfType("model.Text")).Return(tt.expectedErrors.storeErr).Once()
				if tt.wantDelete {
					mockRepo.On("Delete", mock.Anything, model.TextID(textID)).Return(tt.expectedErrors.deleteErr).Once()
				}
			}

			if tt.wantPublish {
				mockPublisher.On("PublishInExchange", mock.MatchedBy(func(evt event.Event) bool {
					rc, ok := evt.(event.RankCalculated)
					assert.True(t, ok, "evt должен быть event.RankCalculated")
					if !ok {
						return false
					}

					if tt.wantRank == 0 {
						assert.InDelta(t, tt.wantRank, rc.Rank, 1e-9)
					} else {
						assert.InEpsilon(t, tt.wantRank, rc.Rank, 1e-9)
					}
					return true
				})).Return(tt.expectedErrors.publishErr).Once()
			}

			svc := NewRankCalculatorService(mockRepo, mockPublisher)
			err := svc.AddText(context.Background(), tt.args.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockPublisher.AssertExpectations(t)
		})
	}
}
