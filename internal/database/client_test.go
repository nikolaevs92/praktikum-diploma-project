package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorization(t *testing.T) {
	// cancelChan := make(chan os.Signal, 1)
	// signal.Notify(cancelChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	// ctx, cancel := context.WithCancel(context.Background())

	// go func() {
	// 	<-cancelChan
	// 	cancel()
	// }()

	t.Run("test_push_order", func(t *testing.T) {
		assert.NoError(t, nil)
	})
}
