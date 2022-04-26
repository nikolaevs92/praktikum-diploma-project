package authorization

import (
	"context"
	"fmt"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
)

type AuthorizationDBTest struct {
	Users []objects.User
}

func (adb *AuthorizationDBTest) Run(ctx context.Context) {
	adb.Users = make([]objects.User, 0)
}

func (adb AuthorizationDBTest) IsLoginExist(login string) bool {
	for _, user := range adb.Users {
		if user.Login == login {
			return true
		}
	}
	return false
}

func (adb AuthorizationDBTest) CheckLoginPasswordHash(login string, paswordHash string) (bool, string) {
	for _, user := range adb.Users {
		if user.Login == login && user.PasswordHash == paswordHash {
			return true, user.UserId
		}
	}
	return false, ""
}

func (adb *AuthorizationDBTest) CreateUser(login string, paswordHash string) error {
	adb.Users = append(adb.Users, objects.User{Login: login, PasswordHash: paswordHash, UserId: fmt.Sprintf("user%d", len(adb.Users))})
	return nil
}
