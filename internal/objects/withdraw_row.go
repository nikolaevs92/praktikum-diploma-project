package objects

import "time"

type WithdrawRow struct {
	UserID      string
	Order       string
	Sum         float64
	ProcessedAt time.Time
}
