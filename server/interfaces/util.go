package interfaces

import (
	"log/slog"
	"net/http"
)

func LogHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request", "method", r.Method, "address", r.RemoteAddr, "path", r.URL)
		// TODO response code
		handler.ServeHTTP(w, r)
	})
}
