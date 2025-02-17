package handler

import (
	"georgslauf/db"
	"georgslauf/handler/templates"
	"georgslauf/persistence"
	"log/slog"
	"net/http"
	"time"
)

type Host struct {
	repository *persistence.Repository
}

func NewHost(
	respository *persistence.Repository,
) *Host {
	return &Host{
		repository: respository,
	}
}

func (h *Host) GetHostHome() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// authCtx := h.authInterceptor.Context(r.Context()) // TODO

		w.WriteHeader(http.StatusOK)
		if err := templates.HostHome(nil).Render(r.Context(), w); err != nil {
			slog.Warn("err rendering HostHome", "err", err)
		}
	})
}

func (h *Host) GetSchedule() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		schedule, err := h.repository.Queries.GetSchedule(ctx)
		if err != nil {
			slog.Warn("error in GetSchedule query", "err", err)
		}

		w.WriteHeader(http.StatusOK)
		if err := templates.HostSchedule(schedule).Render(ctx, w); err != nil {
			slog.Warn("err rendering HostSchedule", "err", err)
		}
	})
}

func (h *Host) GetTribes() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tribes, err := h.repository.Queries.GetTribes(ctx)
		if err != nil {
			slog.Warn("error in GetTribes query", "err", err)
		}

		w.WriteHeader(http.StatusOK)
		if err := templates.HostTribes(tribes).Render(ctx, w); err != nil {
			slog.Warn("err rendering HostTribes", "err", err)
		}
	})
}

func (h *Host) CreateTribe() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
}
