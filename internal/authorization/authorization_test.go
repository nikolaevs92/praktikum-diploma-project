package authorization

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
	"github.com/stretchr/testify/assert"
)

func TestAuthorization(t *testing.T) {
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-cancelChan
		cancel()
	}()

	auth := Authorization{DB: &AuthorizationDBTest{}, Tokens: make(map[string]string)}
	auth.Run(ctx)

	t.Run("test_registration", func(t *testing.T) {
		_, err := auth.Registration(objects.RegisterMessage{Login: "", Password: ""})
		assert.Error(t, err)
		_, err = auth.Registration(objects.RegisterMessage{Login: "", Password: "asdsd"})
		assert.Error(t, err)
		_, err = auth.Registration(objects.RegisterMessage{Login: "12d", Password: ""})
		assert.Error(t, err)

		_, err = auth.Registration(objects.RegisterMessage{Login: "login", Password: "password"})
		assert.NoError(t, err)
		_, err = auth.Registration(objects.RegisterMessage{Login: "login", Password: "passwdsford"})
		assert.Error(t, err)
		_, err = auth.Registration(objects.RegisterMessage{Login: "login", Password: "password"})
		assert.Error(t, err)
		_, err = auth.Registration(objects.RegisterMessage{Login: "login2", Password: "password"})
		assert.NoError(t, err)
	})

	t.Run("test_login", func(t *testing.T) {
		_, err := auth.Login(objects.LoginMessage{Login: "", Password: ""})
		assert.Error(t, err)
		_, err = auth.Login(objects.LoginMessage{Login: "asd", Password: ""})
		assert.Error(t, err)
		_, err = auth.Login(objects.LoginMessage{Login: "", Password: "asdd"})
		assert.Error(t, err)
		_, err = auth.Login(objects.LoginMessage{Login: "e3w", Password: "ewqe"})
		assert.Error(t, err)

		_, err = auth.Login(objects.LoginMessage{Login: "login", Password: "asdasd"})
		assert.Error(t, err)
		_, err = auth.Login(objects.LoginMessage{Login: "login2", Password: ""})
		assert.Error(t, err)

		_, err = auth.Login(objects.LoginMessage{Login: "login", Password: "password"})
		assert.NoError(t, err)
		_, err = auth.Login(objects.LoginMessage{Login: "login2", Password: "password"})
		assert.NoError(t, err)
	})

	t.Run("test_GetUser", func(t *testing.T) {
		token1, err := auth.Login(objects.LoginMessage{Login: "login", Password: "password"})
		assert.NoError(t, err)
		userID1, err := auth.GetUser(token1.Token)
		assert.NoError(t, err)

		token2, err := auth.Login(objects.LoginMessage{Login: "login2", Password: "password"})
		assert.NoError(t, err)
		userID2, err := auth.GetUser(token2.Token)
		assert.NoError(t, err)

		assert.NotEqual(t, token1, token2)
		assert.NotEqual(t, userID1, userID2)
	})
}
