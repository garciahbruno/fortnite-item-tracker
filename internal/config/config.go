package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	Port        string
	DatabaseURL string
	GinMode     string
}

func Load() (Config, error) {
	c := Config{Port: os.Getenv("PORT"), DatabaseURL: os.Getenv("DATABASE_URL"), GinMode: os.Getenv("GIN_MODE")}
	if c.Port == "" {
		c.Port = "8080"
	}
	n, err := strconv.Atoi(c.Port)
	if err != nil || n < 1 || n > 65535 {
		return Config{}, fmt.Errorf("PORT must be between 1 and 65535")
	}
	u, err := url.Parse(c.DatabaseURL)
	if err != nil || u == nil || (u.Scheme != "postgres" && u.Scheme != "postgresql") || u.Host == "" || u.Path == "" || u.Path == "/" {
		return Config{}, fmt.Errorf("DATABASE_URL must be a PostgreSQL URL with a host and database name")
	}
	if c.GinMode == "" {
		c.GinMode = "release"
	}
	if c.GinMode != "release" && c.GinMode != "debug" && c.GinMode != "test" {
		return Config{}, fmt.Errorf("GIN_MODE must be release, debug, or test")
	}
	return c, nil
}
