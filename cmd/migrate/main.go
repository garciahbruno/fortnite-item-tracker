package main

import (
	"context"
	"github.com/garciahbruno/fortnite-item-tracker/internal/config"
	"github.com/garciahbruno/fortnite-item-tracker/internal/database"
	"log/slog"
	"os"
	"time"
)

func main() {
	c, err := config.Load()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	pool, err := database.Open(ctx, c.DatabaseURL)
	if err != nil {
		slog.Error("database unavailable; check connectivity and credentials")
		os.Exit(1)
	}
	defer pool.Close()
	if err := database.Migrate(ctx, pool); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}
	slog.Info("database migrations applied")
}
