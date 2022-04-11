package rest

import (
	"context"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
)

type GofemartInterface interface {
	Run(context.Context)
}

type AuthorizationInterface interface {
	Regist(objects.RegisterMessage) error
}
