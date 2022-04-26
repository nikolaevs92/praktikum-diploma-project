package objects

import "time"

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accural    float64   `json:"accural"`
	UploudedAt time.Time `json:"uploaded_at"`
}

const (
	ORDER_STATUS_NEW        = "NEW"
	ORDER_STATUS_PROCESSING = "PROCESSING"
	ORDER_STATUS_INVALID    = "INVALID"
	ORDER_STATUS_PROCESSED  = "PROCESSED"
)
