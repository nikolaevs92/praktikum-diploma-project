package accuralclient

import "time"

type Config struct {
	AccuralHost string

	Retries int
	Timeout time.Duration
}
