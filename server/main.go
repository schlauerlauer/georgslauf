package main

import (
	"fmt"
	"georgslauf/infra/auth"
	"georgslauf/infra/persistence"
	"georgslauf/interfaces"
	"georgslauf/interfaces/config"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/libsql/go-libsql"
	"github.com/lmittmann/tint"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	version = "24.2.0-alpha" // bumpver
)

func main() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))

	configPath := os.Getenv("FLOW_CONFIG_FILE")
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		slog.Error("could not read config file", "err", err)
		os.Exit(1)
	}

	repository, err := persistence.NewRepository(&cfg.Config.Database)
	if err != nil {
		slog.Error("error connecting repository", "err", err)
		os.Exit(1)
	}

	authService := auth.NewKratosClient(
		cfg.Config.Auth.KratosLocalURL,
		repository,
	)

	publicService := interfaces.NewPublic(
		repository,
		version,
	)

	hostService := interfaces.NewHost(
		repository,
	)

	mux := http.NewServeMux()

	mux.Handle("GET /dist/", http.FileServer(http.Dir("./"))) // TODO ../ safety ?

	mux.HandleFunc("GET /ping", publicService.Ping)
	mux.HandleFunc("GET /version", publicService.Version)

	// mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("GET /", publicService.GetHome)

	mux.Handle("GET /host", authService.Auth.SessionMiddleware(hostService.GetHostHome()))

	slog.Info("server starting", "host", cfg.Config.Server.Host, "port", cfg.Config.Server.Port)
	if err := http.ListenAndServe(fmt.Sprint(cfg.Config.Server.Host, ":", cfg.Config.Server.Port), interfaces.LogHandler(mux)); err != nil {
		slog.Error("err", err)
		os.Exit(1)
	}
}
