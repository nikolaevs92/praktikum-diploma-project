package objects

import "time"

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accural    float64   `json:"accural"`
	UploudedAt time.Time `json:"uploaded_at"`
}

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)
