package api

import (
	"encoding/json"
	"fmt"
	"github.com/eldius/rinha-backend-2025/internal/client"
	"net/http"
	"time"
)

func Start(priority, fallback string) error {
	h := handler{
		p: client.New(priority, 1*time.Millisecond),
		f: client.New(fallback, 1*time.Second),
	}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", h.index)
	mux.HandleFunc("POST /payments", h.payments)

	server := &http.Server{Addr: ":8080", Handler: mux}
	return server.ListenAndServe()
}

type handler struct {
	p *client.Client
	f *client.Client
}

func (h *handler) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("<h1>Hello World</h1>"))
}

type Payment struct {
	CorrelationID string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
}

func (h *handler) payments(w http.ResponseWriter, r *http.Request) {
	var payment Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		err = fmt.Errorf("invalid payment format: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	providerPayment := client.ProviderPaymentRequest{
		CorrelationId: payment.CorrelationID,
		Amount:        payment.Amount,
		RequestedAt:   time.Now(),
	}
	resp, err := h.p.Pay(providerPayment)
	if err != nil {
		var fallErr error
		resp, fallErr = h.f.Pay(providerPayment)
		if fallErr != nil {
			err = fmt.Errorf("priority failed: %w > fallback failed: %w", err, fallErr)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
	}
	_ = json.NewEncoder(w).Encode(resp)
}
