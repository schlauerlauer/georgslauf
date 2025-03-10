package handler

import (
	"georgslauf/internal/db"
	"georgslauf/internal/handler/templates"
	"georgslauf/internal/settings"
	"georgslauf/md"
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
	queries       *db.Queries
	formProcessor *forms.FormProcessor
	session       *session.Session
	settings      *settings.SettingsService
	md            *md.MdService
}

func NewHandler(
	queries *db.Queries,
	session *session.Session,
	settings *settings.SettingsService,
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

	set := settings.Get()
	mdData := md.New()
	if _, err := mdData.Update(set.Home); err != nil {
		slog.Error("Invalid Markdown data", "err", err)
		return nil, err
	}

	return &Handler{
		queries:       queries,
		formProcessor: formProcessor,
		session:       session,
		settings:      settings,
		md:            mdData,
	}, nil
}
