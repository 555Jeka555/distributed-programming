package model

import "context"

type User struct {
	login    string
	password string
}

type UserRepository interface {
	UserReadRepository
	Store(ctx context.Context, user User) error
}

type UserReadRepository interface {
	FindByLogin(ctx context.Context, login string) (User, error)
}

func NewUser(login string, password string) User {
	return User{
		login:    login,
		password: password,
	}
}

func (u *User) Login() string {
	return u.login
}

func (u *User) Password() string {
	return u.password
}

func LoadUser(login string, password string) User {
	return User{
		login:    login,
		password: password,
	}
}
