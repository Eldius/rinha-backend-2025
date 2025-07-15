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

type ProviderSummary struct {
	TotalRequests int     `json:"totalRequests"`
	TotalAmount   float64 `json:"totalAmount"`
}
type PaymentsSummary struct {
	Default  ProviderSummary `json:"default"`
	Fallback ProviderSummary `json:"fallback"`
}
