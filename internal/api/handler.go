package api

import (
	"encoding/json"
	"fmt"
	"github.com/eldius/rinha-backend-2025/internal/client"
	"github.com/eldius/rinha-backend-2025/internal/persistence"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
	"time"
)

type customResponseWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (h *customResponseWriter) WriteHeader(code int) {
	h.status = code
	h.ResponseWriter.WriteHeader(code)
}

func (h *customResponseWriter) Write(b []byte) (int, error) {
	h.body = b
	return h.ResponseWriter.Write(b)
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logData := map[string]interface{}{
			"mrthod":     r.Method,
			"url":        r.URL.String(),
			"start_time": start,
		}
		writer := customResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(&writer, r)
		logData["status"] = writer.status
		logData["duration"] = time.Since(start)

		logData["response"] = string(writer.body)

		slog.With("request", logData).Info("request completed", logData)
	})
}

func Start(priority, fallback string) error {
	mux := http.NewServeMux()

	db, err := sqlx.Connect("postgres", "postgres://app:MyStrongP%40ss@db:5432/rinha?sslmode=disable")
	if err != nil {
		return err
	}
	h := newHandler(db, priority, fallback)
	defer func() {
		_ = db.Close()
	}()
	mux.HandleFunc("GET /", h.index)
	mux.HandleFunc("POST /payments", h.payments)
	mux.HandleFunc("GET /payments-summary", h.summary)

	server := &http.Server{
		Addr:    ":8080",
		Handler: logMiddleware(mux),
	}
	return server.ListenAndServe()
}

type handler struct {
	p *client.Client
	f *client.Client
	r *persistence.PaymentRepository
}

func newHandler(db *sqlx.DB, primary, fallback string) *handler {
	return &handler{
		r: persistence.New(db),
		p: client.New("default", primary, 500*time.Millisecond),
		f: client.New("fallback", fallback, 5*time.Second),
	}
}

func (h *handler) index(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("<h1>Hello World</h1>"))
}

func (h *handler) summary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.r.Summary()
	if err != nil {
		slog.With("error", err).Error("failed to get summary")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	_ = json.NewEncoder(w).Encode(summary)
}

type Payment struct {
	CorrelationID string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
}

func (h *handler) payments(w http.ResponseWriter, r *http.Request) {
	var payment Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		err = fmt.Errorf("invalid payment format: %v", err)
		slog.With("error", err).Error("failed to decode payment")
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
		slog.With("error", err).Warn("failed to pay with primary provider, trying fallback...")
		var fallErr error
		resp, fallErr = h.f.Pay(providerPayment)
		if fallErr != nil {
			err = fmt.Errorf("priority failed: %w > fallback failed: %w", err, fallErr)
			slog.With("error", err).Error("failed to pay using fallback provider")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
	}

	go func() {
		_ = h.r.Save(*resp)
	}()
	_ = json.NewEncoder(w).Encode(resp)
}
