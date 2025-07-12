package model

import "time"

type PaymentInfo struct {
	CorrelationId string    `json:"correlationId" db:"correlation_id"`
	Amount        float64   `json:"amount" db:"amount"`
	RequestedAt   time.Time `json:"requestedAt" db:"requested_at"`
	Message       string    `json:"message" db:"message"`
	Provider      string    `json:"provider" db:"provider"`
	Status        string    `json:"status" db:"status"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
}
