package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/garciahbruno/fortnite-item-tracker/internal/web"
	"github.com/gin-gonic/gin"
)

type Pinger interface{ Ping(context.Context) error }

func New(db Pinger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	// No proxy is trusted until a deployment explicitly configures one.
	_ = r.SetTrustedProxies(nil)
	r.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "same-origin")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self'; connect-src 'self'; frame-ancestors 'none'; base-uri 'none'; form-action 'self'")
		c.Next()
	})
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.GET("/readyz", func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if db == nil || db.Ping(ctx) != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})
	r.GET("/api/v1/status", func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.JSON(http.StatusOK, gin.H{"service": "fortnite-item-tracker", "stage": "foundation", "features": gin.H{"shop": false, "authentication": false, "watchlists": false, "discord_alerts": false}})
	})
	for path, asset := range map[string]string{"/": "index.html", "/styles.css": "styles.css", "/app.js": "app.js"} {
		r.GET(path, func(c *gin.Context) {
			contentType := map[string]string{"index.html": "text/html; charset=utf-8", "styles.css": "text/css; charset=utf-8", "app.js": "text/javascript; charset=utf-8"}[asset]
			c.Data(http.StatusOK, contentType, web.Read(asset))
		})
	}
	r.NoRoute(func(c *gin.Context) { c.JSON(http.StatusNotFound, gin.H{"error": "not_found"}) })
	return r
}
