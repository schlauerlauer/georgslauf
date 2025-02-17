package handler

import (
	"fmt"
	"georgslauf/handler/templates"
	"georgslauf/persistence"
	"log/slog"
	"net/http"
	"time"
)

type Public struct {
	repository *persistence.Repository
}

func NewPublic(
	repository *persistence.Repository,
) *Public {
	parsedBuildTime, err := time.Parse("2006-01-02T15:04:05", buildTime)
	if err != nil {
		slog.Error("could not parse build time", "err", err)
	}
	expirationTime = parsedBuildTime.AddDate(2, 0, 0)
	templates.SetVersion(version) // not sure why ldflags don't work for this

	return &Public{
		repository: repository,
	}
}

var (
	version    = "n/a" // set by ldflags
	commitHash = "n/a" // set by ldflags
	buildTime  = "n/a" // set by ldflags

	expirationTime = time.Time{}
)

func (p *Public) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func (p *Public) Version(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(version))
}

func (p *Public) Robots(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User-agent: *\nAllow: /"))
}

func (p *Public) Security(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	text := fmt.Sprintf(`Contact: mailto:security@georgslauf.de
Expires: %s
Preferred-Languages: de, en
Canonical: https://georgslauf.de/.well-known/security.txt`, expirationTime.Format("2006-01-02T15.04.005Z"))
	w.Write([]byte(text))
}

func (p *Public) GetHome() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// TODO userinfo

		schedule, err := p.repository.Queries.GetSchedule(ctx)
		if err != nil {
			slog.Warn("error in GetSchedule query", "err", err)
		}

		w.WriteHeader(http.StatusOK)
		if err := templates.Home(nil, schedule).Render(ctx, w); err != nil {
			slog.Warn("err rendering home", "err", err)
		}
	})
}
