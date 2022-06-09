package authorization

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
)

type Authorization struct {
	DB  AuthorizationDBInterface
	Cfg Config

	Tokens map[string]string
}

func New(db AuthorizationDBInterface, config Config) Authorization {
	return Authorization{DB: db, Cfg: config, Tokens: map[string]string{}}
}

func (a Authorization) Run(ctx context.Context) {
	// a.DB.Run(ctx)
}

func (a Authorization) Registration(message objects.RegisterMessage) (objects.TokenMessage, error) {
	if message.Login == "" || message.Password == "" {
		return objects.TokenMessage{}, errors.New("Login and password should be empty")
	}

	if a.DB.IsLoginExist(message.Login) {
		return objects.TokenMessage{}, errors.New("user already exist")
	}
	err := a.DB.CreateUser(message.Login, getPasswordHash(message.Password))
	if err == nil {
		log.Printf("User with login: %s was succesfully created", message.Login)
	}
	token := fmt.Sprintf("token%d", len(a.Tokens))
	a.Tokens[token] = message.Login
	return objects.TokenMessage{Token: token}, nil
}

func (a Authorization) Login(message objects.LoginMessage) (objects.TokenMessage, error) {
	ok, userID := a.DB.CheckLoginPasswordHash(message.Login, getPasswordHash(message.Password))
	if !ok {
		return objects.TokenMessage{}, errors.New("wrong login or password")
	}

	token := fmt.Sprintf("token%d", len(a.Tokens))
	a.Tokens[token] = userID
	return objects.TokenMessage{Token: token}, nil
}

func (a Authorization) GetUser(token string) (string, error) {
	userID, ok := a.Tokens[token]
	if !ok {
		return "", errors.New("wrong token")
	}
	return userID, nil
}

func getPasswordHash(password string) string {
	// TODO
	return password
}
