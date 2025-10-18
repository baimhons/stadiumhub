package response

import (
	"time"

	"github.com/google/uuid"
)

type PaymentResponse struct {
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
	SessionURL StripeModel `json:"seesion_url"`
}

type StripeModel struct {
	UserID     uuid.UUID `json:"user_id"`
	SessionURL string    `json:"session_url"`
	SessionID  string    `json:"session_id"`
	Amount     float32   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}
