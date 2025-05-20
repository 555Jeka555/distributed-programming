package repo

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"server/pkg/app/model"
	"server/pkg/infrastructure/keyvalue"
)

func NewUserRepository(rdb *redis.Client) model.UserRepository {
	return &userRepository{
		storage: keyvalue.NewStorage[userSerializable](rdb),
	}
}

type userSerializable struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type userRepository struct {
	storage keyvalue.Storage[userSerializable]
}

func (r *userRepository) FindByLogin(ctx context.Context, login string) (model.User, error) {
	user, err := r.storage.Get(ctx, login)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.User{}, nil
		}
		return model.User{}, err
	}

	return model.LoadUser(
		user.Login,
		user.Password,
	), nil
}

func (r *userRepository) Store(ctx context.Context, user model.User) error {
	return r.storage.Set(ctx, user.Login(), userSerializable{
		Login:    user.Login(),
		Password: user.Password(),
	}, 0)
}
