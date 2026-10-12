package httpapi

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeDB struct{ err error }

func (f fakeDB) Ping(context.Context) error { return f.err }
func TestRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		name, path string
		db         Pinger
		status     int
		contains   string
	}{
		{"live despite database failure", "/healthz", fakeDB{errors.New("secret connection detail")}, 200, "ok"},
		{"ready", "/readyz", fakeDB{}, 200, "ready"},
		{"database down", "/readyz", fakeDB{errors.New("secret connection detail")}, 503, "unavailable"},
		{"missing database", "/readyz", nil, 503, "unavailable"},
		{"status", "/api/v1/status", fakeDB{}, 200, `"shop":false`},
		{"home", "/", fakeDB{}, 200, "Fortnite Item Tracker"},
		{"css", "/styles.css", fakeDB{}, 200, "background"},
		{"script", "/app.js", fakeDB{}, 200, "fetch"},
		{"unknown API", "/api/v1/missing", fakeDB{}, 404, "not_found"},
		{"private files", "/.env", fakeDB{}, 404, "not_found"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			New(tc.db).ServeHTTP(w, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if w.Code != tc.status || !strings.Contains(w.Body.String(), tc.contains) {
				t.Fatalf("got %d %s", w.Code, w.Body.String())
			}
			if strings.Contains(w.Body.String(), "secret") {
				t.Fatal("leaked internal error")
			}
			if w.Header().Get("Content-Security-Policy") == "" {
				t.Fatal("missing security headers")
			}
		})
	}
}
