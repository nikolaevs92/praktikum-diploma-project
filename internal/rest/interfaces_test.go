package rest

import (
	"context"
	"errors"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
)

type AuthorizationTest struct {
}

func (a AuthorizationTest) Registration(message objects.RegisterMessage) error {
	if message.Login == "exist_user" {
		return errors.New("user already exist")
	}
	return nil
}

func (a AuthorizationTest) Login(message objects.LoginMessage) (objects.TokenMessage, error) {
	if message.Login == "exist_user" && message.Password == "correctpassword" {
		return objects.TokenMessage{Token: "token"}, nil
	}
	return objects.TokenMessage{}, errors.New("user doesnt exist")
}

type GofemartTest struct {
}

func (g GofemartTest) Run(end context.Context) {
}
