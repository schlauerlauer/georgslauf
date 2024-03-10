package interfaces

import (
	"georgslauf/infra/persistence"
	"georgslauf/view/public"
	"log/slog"
	"net/http"
)

type Public struct {
	repository *persistence.Repository
	version    string
}

func NewPublic(
	repository *persistence.Repository,
	version string,
) *Public {
	return &Public{
		repository: repository,
		version:    version,
	}
}

func (p *Public) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (p *Public) Version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(p.version))
}

func (p *Public) GetHome(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	schedule, err := p.repository.Queries.GetSchedule(ctx)
	if err != nil {
		slog.Warn("error in GetSchedule query", "err", err)
	}

	w.WriteHeader(http.StatusOK)
	err = public.Home(schedule, p.repository.Location).Render(ctx, w)
	if err != nil {
		slog.Warn("err rendering home", "err", err)
	}
}
