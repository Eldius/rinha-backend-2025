package client

import (
	"bytes"
	"encoding/json"
	"fmt"
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
type ProviderPaymentResponse struct {
	ProviderPaymentRequest
	Message  string `json:"message"`
	Provider string `json:"provider"`
}

func New(backend string, timeout time.Duration) *Client {
	return &Client{
		c: &http.Client{
			Timeout: timeout,
		},
		backend: backend,
	}
}

func (c *Client) Pay(p ProviderPaymentRequest) (*ProviderPaymentResponse, error) {
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

	var pRes ProviderPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("decoding response: %v", err)
	}

	pRes.Provider = c.backend
	pRes.ProviderPaymentRequest = p

	return &pRes, nil
}
