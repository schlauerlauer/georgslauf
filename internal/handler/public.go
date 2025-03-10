package handler

import (
	"fmt"
	"georgslauf/htmx"
	"georgslauf/internal/handler/templates"
	"georgslauf/session"
	"log/slog"
	"net/http"
)

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func (h *Handler) Version(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(version))
}

func (h *Handler) Robots(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User-agent: *\nAllow: /"))
}

func (h *Handler) Security(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	text := fmt.Sprintf(`Contact: mailto:security@georgslauf.de
Expires: %s
Preferred-Languages: de, en
Canonical: https://georgslauf.de/.well-known/security.txt`, expirationTime.Format("2006-01-02T15.04.005Z"))
	w.Write([]byte(text))
}

func (h *Handler) GetHome(w http.ResponseWriter, r *http.Request) {
	htmxRequest := htmx.IsHTMX(r)
	ctx := r.Context()

	var user *session.UserData
	if userData, ok := ctx.Value(session.ContextKey).(*session.UserData); ok {
		user = userData
	}

	schedule, err := h.queries.GetSchedule(ctx)
	if err != nil {
		slog.Warn("GetSchedule", "err", err)
	}

	setMd := h.md.Get()

	w.WriteHeader(http.StatusOK)
	if err := templates.Home(htmxRequest, user, schedule, setMd).Render(ctx, w); err != nil {
		slog.Warn("Home", "err", err)
	}
}
