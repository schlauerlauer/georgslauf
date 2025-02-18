package handler

import (
	"georgslauf/handler/templates"
	"georgslauf/persistence"
	"log/slog"
	"time"

	"github.com/schlauerlauer/go-forms"
)

var (
	version   = "n/a" // set by ldflags
	buildTime = "n/a" // set by ldflags

	expirationTime = time.Time{}
)

type Handler struct {
	repository    *persistence.Repository
	formProcessor *forms.FormProcessor
}

func NewHandler(
	repository *persistence.Repository,
) (*Handler, error) {
	parsedBuildTime, err := time.Parse("2006-01-02T15:04:05", buildTime)
	if err != nil {
		slog.Error("could not parse build time", "err", err)
	}
	expirationTime = parsedBuildTime.AddDate(2, 0, 0)

	templates.SetVersion(version) // not sure why ldflags don't work for this
	templates.SetYear(parsedBuildTime.Format("2006"))

	// FormProcessor
	formProcessor, err := forms.NewFormProcessor()
	if err != nil {
		slog.Error("NewFormProcessor", "err", err)
		return nil, err
	}

	return &Handler{
		repository:    repository,
		formProcessor: formProcessor,
	}, nil
}
