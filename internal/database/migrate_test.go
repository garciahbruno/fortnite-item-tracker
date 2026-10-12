package database

import (
	"context"
	"os"
	"testing"
	"time"
)

// Run only against a dedicated test database; the transaction rolls fixture data back.
func TestMigrationsAndConstraints(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not configured")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := Open(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	for range 2 {
		if err := Migrate(ctx, pool); err != nil {
			t.Fatal(err)
		}
	}
	var count int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM schema_migrations").Scan(&count); err != nil || count != 1 {
		t.Fatalf("migration count %d, error %v", count, err)
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(context.Background())
	var userID int64
	if err := tx.QueryRow(ctx, "INSERT INTO users (discord_id, username) VALUES ('123456789', 'test') RETURNING id").Scan(&userID); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(ctx, "INSERT INTO cosmetics (id,name,cosmetic_type) VALUES ('test-cosmetic','Example','outfit')"); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if _, err := tx.Exec(ctx, "INSERT INTO watchlist_items (user_id, cosmetic_id) VALUES ($1,'test-cosmetic') ON CONFLICT DO NOTHING", userID); err != nil {
			t.Fatal(err)
		}
		if _, err := tx.Exec(ctx, "INSERT INTO notifications (user_id, cosmetic_id, shop_date) VALUES ($1,'test-cosmetic','2026-10-11') ON CONFLICT DO NOTHING", userID); err != nil {
			t.Fatal(err)
		}
	}
	for _, table := range []string{"watchlist_items", "notifications"} {
		if err := tx.QueryRow(ctx, "SELECT count(*) FROM "+table+" WHERE user_id=$1", userID).Scan(&count); err != nil || count != 1 {
			t.Fatalf("%s duplicate protection failed: %d %v", table, count, err)
		}
	}
	if _, err := tx.Exec(ctx, "DELETE FROM users WHERE id=$1", userID); err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"watchlist_items", "notifications"} {
		if err := tx.QueryRow(ctx, "SELECT count(*) FROM "+table+" WHERE user_id=$1", userID).Scan(&count); err != nil || count != 0 {
			t.Fatalf("%s cascade failed: %d %v", table, count, err)
		}
	}
}
