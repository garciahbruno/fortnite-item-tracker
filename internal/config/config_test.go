package config

import "testing"

func TestLoad(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://tracker:secret@localhost:5432/tracker")
	t.Setenv("PORT", "")
	t.Setenv("GIN_MODE", "")
	c, err := Load()
	if err != nil || c.Port != "8080" || c.GinMode != "release" {
		t.Fatalf("unexpected config: %+v %v", c, err)
	}
	for _, value := range []string{"0", "65536", "abc"} {
		t.Setenv("PORT", value)
		if _, err := Load(); err == nil {
			t.Fatalf("accepted invalid port %q", value)
		}
	}
	t.Setenv("PORT", "8080")
	for _, value := range []string{"", "https://example.com/db", "postgres://localhost"} {
		t.Setenv("DATABASE_URL", value)
		if _, err := Load(); err == nil {
			t.Fatalf("accepted invalid database URL")
		}
	}
	t.Setenv("DATABASE_URL", "postgres://localhost/tracker")
	t.Setenv("GIN_MODE", "invalid")
	if _, err := Load(); err == nil {
		t.Fatal("accepted invalid mode")
	}
}
