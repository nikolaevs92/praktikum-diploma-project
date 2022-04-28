package accuralclient

import "time"

type Config struct {
	AccuralHost string

	Retries int
	Timeout time.Duration
}

func GetDefaultConfig() Config {
	return Config{
		AccuralHost: "localhost:4444",
		Retries:     2,
		Timeout:     time.Second,
	}
}
