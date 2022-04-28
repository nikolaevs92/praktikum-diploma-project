package objects

import "time"

type OrderRow struct {
	UserID     string
	Number     string
	Status     string
	Accural    float64
	UploudedAt time.Time
}
