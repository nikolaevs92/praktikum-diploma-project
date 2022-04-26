package authorization

import "context"

type AuthorizationDB interface {
	IsLoginExist(string) bool
	CheckLoginPasswordHash(string, string) (bool, string)
	CreateUser(string, string) error
	Run(context.Context)
}
