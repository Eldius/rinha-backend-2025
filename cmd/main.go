package main

import (
	"github.com/eldius/rinha-backend-2025/internal/api"
	"log/slog"
	"os"
	"slices"
	"strings"
)

const (
	appName = "payment-gateway"
)

var (
	AppVersion string

	logKeys = []string{
		"host",
		"service.name",
		"level",
		"message",
		"time",
		"error",
		"source",
		"function",
		"file",
		"line",
	}
)

func logAttrsReplacerFunc() func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if slices.Contains(logKeys, a.Key) {
			return a
		}
		if strings.HasPrefix(a.Key, "request") ||
			strings.HasPrefix(a.Key, "response") ||
			strings.HasPrefix(a.Key, "service") {
			return a
		}

		if slices.Contains(logKeys, a.Key) {
			return a
		}
		if a.Key == "msg" {
			a.Key = "message"
			return a
		}
		return a
	}
}

func init() {

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: logAttrsReplacerFunc(),
	})
	logger := slog.New(handler)
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}

	slog.SetDefault(logger.With(
		slog.String("name", appName),
		slog.String("service_name", appName),
		slog.String("service.version", AppVersion),
		slog.String("host", host),
	))
}

func main() {
	slog.Default().Info("starting")
	backend, ok := os.LookupEnv("API_PRIMARY_BACKEND")
	if !ok {
		panic("API_PRIMARY_BACKEND env var is not set")
	}
	fallback, ok := os.LookupEnv("API_FALLBACK_BACKEND")
	if !ok {
		panic("API_FALLBACK_BACKEND env var is not set")
	}
	if err := api.Start(backend, fallback); err != nil {
		panic(err)
	}
}
