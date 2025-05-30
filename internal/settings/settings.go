package settings

import (
	"context"
	"database/sql"
	"encoding/json"
	"georgslauf/internal/db"
	"georgslauf/md"
	"log/slog"
	"os"
	"sync"
)

type SettingsService struct {
	queries  *db.Queries
	settings Settings
	lock     sync.RWMutex
}

type Settings struct {
	Login    Login    `json:"l"`
	Groups   Groups   `json:"g"`
	Stations Stations `json:"s"`
	Help     Help     `json:"h"`
	Home     md.Input `json:"hp"`
}

type Login struct {
	Title   string `json:"t" schema:"title" validate:"max=64" mod:"trim,sanitize"`
	Welcome string `json:"w" schema:"welcome" validate:"max=1024" mod:"trim,sanitize"`
}

type Groups struct {
	AllowCreate bool  `json:"c" schema:"group-create"`
	AllowUpdate bool  `json:"u" schema:"group-update"`
	AllowDelete bool  `json:"d" schema:"group-delete"`
	ShowAbbr    bool  `json:"s" schema:"group-abbr"`
	Min         int64 `json:"min" schema:"group-min" validate:"gte=0"`
	Max         int64 `json:"max" schema:"group-max" validate:"gte=0,gtfield=Min"`
}

type Stations struct {
	AllowCreate         bool `json:"c" schema:"station-create"`
	AllowUpdate         bool `json:"u" schema:"station-update"`
	AllowDelete         bool `json:"d" schema:"station-delete"`
	EnableCategories    bool `json:"ca" schema:"station-categories"`
	ShowAbbr            bool `json:"s" schema:"station-abbr"`
	AllowScoring        bool `json:"x" schema:"station-scoring"`
	EditAccounts        bool `json:"e" schema:"station-accounts"`
	EditAccountsStation bool `json:"es" schema:"station-accounts-self"`
}

type Help struct {
	Footer string `json:"f" schema:"footer" validate:"http_url,max=128"`
}

func New(queries *db.Queries) *SettingsService {
	service := SettingsService{
		queries: queries,
		settings: Settings{
			Groups: Groups{
				Max: 1,
			},
		},
		lock: sync.RWMutex{},
	}

	if res, err := queries.GetSettings(context.Background()); err != nil {
		slog.Info("no settings found, setting defaults")

		if data, err := json.Marshal(service.settings); err != nil {
			slog.Error("Marshal", "err", err)
			os.Exit(1)
		} else {
			if err := queries.InsertSettings(context.Background(), data); err != nil {
				slog.Error("InsertSettings", "err", err)
				os.Exit(1)
			}
		}
	} else {
		if err := json.Unmarshal(res.Data, &service.settings); err != nil {
			slog.Error("Unmarshal", "err", err)
			os.Exit(1)
		}
	}

	return &service
}

func (s *SettingsService) Get() Settings {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.settings
}

func (s *SettingsService) Set(ctx context.Context, settings Settings, userId int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	if err := s.queries.UpdateSettings(ctx, db.UpdateSettingsParams{
		Data: data,
		UpdatedBy: sql.NullInt64{
			Int64: userId,
			Valid: true,
		},
	}); err != nil {
		return err
	}

	s.settings = settings

	return nil
}
