package handler

import (
	"georgslauf/db"
	"georgslauf/handler/templates"
	"georgslauf/htmx"
	"log/slog"
	"net/http"
	"time"
)

func (h *Handler) GetHostHome(w http.ResponseWriter, r *http.Request) {
	// authCtx := h.authInterceptor.Context(r.Context()) // TODO

	w.WriteHeader(http.StatusOK)
	htmxRequest := htmx.IsHTMX(r)

	if err := templates.HostHome(htmxRequest, nil).Render(r.Context(), w); err != nil {
		slog.Warn("err rendering HostHome", "err", err)
	}
}

func (h *Handler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	schedule, err := h.repository.Queries.GetSchedule(ctx)
	if err != nil {
		slog.Warn("GetSchedule", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := templates.HostSchedule(schedule).Render(ctx, w); err != nil {
		slog.Warn("HostSchedule", "err", err)
		return
	}
}

func (h *Handler) GetTribes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tribes, err := h.repository.Queries.GetTribes(ctx)
	if err != nil {
		slog.Warn("GetTribes", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := templates.HostTribes(tribes).Render(ctx, w); err != nil {
		slog.Warn("HostTribes", "err", err)
		return
	}
}

func (h *Handler) CreateTribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tribe, err := h.repository.Queries.CreateTribe(ctx, db.CreateTribeParams{
		Name:      "TESTING",
		UpdatedAt: time.Now().Unix(),
	})
	if err != nil {
		slog.Warn("error in CreateTribe query", "err", err)
	}

	w.WriteHeader(http.StatusCreated)
	if err := templates.HostTribe(tribe).Render(ctx, w); err != nil {
		slog.Warn("err rendering HostTribe", "err", err)
	}
}
