package gofermart

import "context"

type GofemartInterface interface {
	Run(end context.Context)
}
