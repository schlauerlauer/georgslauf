package handler

import (
	"georgslauf/handler/templates"
	"georgslauf/persistence"
	"log/slog"
	"time"
)

var (
	version   = "n/a" // set by ldflags
	buildTime = "n/a" // set by ldflags

	expirationTime = time.Time{}
)

type Handler struct {
	repository *persistence.Repository
}

func NewHandler(
	repository *persistence.Repository,
) *Handler {
	parsedBuildTime, err := time.Parse("2006-01-02T15:04:05", buildTime)
	if err != nil {
		slog.Error("could not parse build time", "err", err)
	}
	expirationTime = parsedBuildTime.AddDate(2, 0, 0)

	templates.SetVersion(version) // not sure why ldflags don't work for this
	templates.SetYear(parsedBuildTime.Format("2006"))

	return &Handler{
		repository: repository,
	}
}
