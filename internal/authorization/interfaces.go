package authorization

import "context"

type AuthorizationDBInterface interface {
	IsLoginExist(string) bool
	CheckLoginPasswordHash(string, string) (bool, string)
	CreateUser(string, string) error
	Run(context.Context)
}
