package rest

import (
	"context"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
)

type GofemartInterface interface {
	Run(context.Context)
}

type AuthorizationInterface interface {
	Registration(objects.RegisterMessage) error
	Login(objects.LoginMessage) (objects.TokenMessage, error)
}
