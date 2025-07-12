package main

import (
	"log/slog"
	"time"
)

var (
	AppVersion string
)

func main() {
	for {
		slog.With("version", AppVersion).Info("Hello World")
		time.Sleep(1 * time.Second)
	}
}
