package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/eldius/rinha-backend-2025/internal/model"
	"net/http"
	"time"
)

type Client struct {
	c       *http.Client
	backend string
}

type ProviderPaymentRequest struct {
	CorrelationId string    `json:"correlationId"`
	Amount        float64   `json:"amount"`
	RequestedAt   time.Time `json:"requestedAt"`
}

func New(backend string, timeout time.Duration) *Client {
	return &Client{
		c: &http.Client{
			Timeout: timeout,
		},
		backend: backend,
	}
}

func (c *Client) Pay(p ProviderPaymentRequest) (*model.PaymentInfo, error) {
	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(p); err != nil {
		return nil, fmt.Errorf("encoding payment: %v", err)
	}
	req, err := http.NewRequest("POST", c.backend+"/payments", bytes.NewReader(buff.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var pRes model.PaymentInfo
	if err := json.NewDecoder(resp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("decoding response: %v", err)
	}

	pRes.Provider = c.backend
	pRes.Amount = p.Amount
	pRes.CorrelationId = p.CorrelationId
	pRes.RequestedAt = p.RequestedAt

	return &pRes, nil
}
