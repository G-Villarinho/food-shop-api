package config

import (
	"log/slog"
	"os"
)

func ConfigureLogger() {
	handler := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
	}))
	slog.SetDefault(handler)
}
