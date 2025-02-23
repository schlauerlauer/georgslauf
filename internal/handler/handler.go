package handler

import (
	"georgslauf/internal/db"
	"georgslauf/internal/handler/templates"
	"georgslauf/session"
	"log/slog"
	"time"

	"github.com/schlauerlauer/go-forms"
)

var (
	version   = "n/a"                 // set by ldflags
	buildTime = "2006-01-02T15:04:05" // set by ldflags

	expirationTime = time.Time{}
)

type Handler struct {
	repository    *db.Repository
	formProcessor *forms.FormProcessor
	session       *session.Session
}

func NewHandler(
	repository *db.Repository,
	session *session.Session,
) (*Handler, error) {
	parsedBuildTime, err := time.Parse("2006-01-02T15:04:05", buildTime)
	if err != nil {
		slog.Error("could not parse build time", "err", err)
	}
	expirationTime = parsedBuildTime.AddDate(2, 0, 0)

	templates.SetVars(version, parsedBuildTime.Format("2006"))

	// FormProcessor
	formProcessor, err := forms.NewFormProcessor()
	if err != nil {
		slog.Error("NewFormProcessor", "err", err)
		return nil, err
	}

	return &Handler{
		repository:    repository,
		formProcessor: formProcessor,
		session:       session,
	}, nil
}
