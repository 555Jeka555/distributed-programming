package service

import (
	"context"
	"errors"
	"log"
	"server/pkg/app/model"
)

var ErrKeyAlreadyExists = errors.New("key already exists")

func NewUserService(hashService HashService, userRepository model.UserRepository) UserService {
	return UserService{
		hashService:    hashService,
		userRepository: userRepository,
	}
}

type UserService struct {
	hashService    HashService
	userRepository model.UserRepository
}

func (a *UserService) CreateUser(ctx context.Context, login, password string) error {
	user, err := a.userRepository.FindByLogin(ctx, login)
	if err != nil {
		return err
	}

	if user.Login() == login {
		log.Printf("User already exist %s", user.Login())
		if user.Password() != a.hashService.Hash(password) {
			return errors.New("password not matched")
		}
		return nil
	}

	return a.userRepository.Store(ctx, model.NewUser(login, a.hashService.Hash(password)))
}

func (a *UserService) Authenticate(ctx context.Context, login, password string) (bool, error) {
	user, err := a.userRepository.FindByLogin(ctx, login)
	if err != nil {
		return false, err
	}

	if user.Login() != login {
		log.Printf("User not exist %s", user.Login())
		return false, nil
	}

	if user.Password() != a.hashService.Hash(password) {
		log.Printf("Password not match %s", user.Password())
		return false, nil
	}

	return true, nil
}
