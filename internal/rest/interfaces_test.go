package rest

import (
	"context"
	"errors"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
)

type AuthorizationTest struct {
}

func (a AuthorizationTest) Regist(message objects.RegisterMessage) error {
	if message.Login == "exist_user" {
		return errors.New("user already exist")
	}
	return nil
}

type GofemartTest struct {
}

func (g GofemartTest) Run(end context.Context) {
}
