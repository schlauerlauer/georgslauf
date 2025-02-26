package settings

import (
	"context"
	"encoding/json"
	"georgslauf/internal/db"
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
	Groups Groups `json:"g"`
}

type Groups struct {
	AllowCreate bool  `json:"c"`
	AllowUpdate bool  `json:"u"`
	AllowDelete bool  `json:"d"`
	Min         int64 `json:"min"`
	Max         int64 `json:"max"`
}

func New(queries *db.Queries) *SettingsService {
	service := SettingsService{
		queries:  queries,
		settings: Settings{},
		lock:     sync.RWMutex{},
	}

	if res, err := queries.GetSettings(context.Background()); err != nil {
		slog.Warn("GetSettings", "err", err)

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

func (s *SettingsService) Set(ctx context.Context, settings Settings) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	if err := s.queries.UpdateSettings(ctx, data); err != nil {
		return err
	}

	s.settings = settings

	return nil
}
