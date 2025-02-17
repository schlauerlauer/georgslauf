package main

import (
	"embed"
	"fmt"
	"georgslauf/config"
	"georgslauf/handler"
	"georgslauf/middleware"
	"georgslauf/persistence"
	"io/fs"

	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

//go:embed all:dist
var embedDist embed.FS

func main() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.TimeOnly,
		}),
	))

	configPath := os.Getenv("CONFIG_PATH")
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		slog.Error("could not read config file", "err", err)
		os.Exit(1)
	}

	repository, err := persistence.NewRepository(&cfg.Database)
	if err != nil {
		slog.Error("error connecting repository", "err", err)
		os.Exit(1)
	}

	handlers := handler.NewHandler(
		repository,
	)

	router := http.NewServeMux()

	subDist, err := fs.Sub(embedDist, "dist")
	if err != nil {
		slog.Error("fs.Sub", "err", err)
		os.Exit(1)
	}
	distServer := http.FileServer(http.FS(subDist))
	router.Handle("GET /dist", http.NotFoundHandler())
	router.Handle("GET /dist/", http.StripPrefix("/dist", distServer))

	router.HandleFunc("GET /ping", handlers.Ping)
	router.HandleFunc("GET /version", handlers.Version)
	router.HandleFunc("GET /robots.txt", handlers.Robots)
	router.HandleFunc("GET /.well-known/security.txt", handlers.Security)

	// router.Handle("GET /metrics", promhttp.Handler())
	router.HandleFunc("/", handlers.GetHome) // TODO optional auth

	// DASH ROUTES
	privateRouter := http.NewServeMux()
	privateRouter.HandleFunc("GET /", handlers.GetHostHome)
	router.Handle("/dash/", http.StripPrefix("/dash", privateRouter)) // TODO authenticated

	// HTMX ROUTES
	apiRouter := http.NewServeMux()
	apiRouter.HandleFunc("GET /schedule", handlers.GetSchedule)
	apiRouter.HandleFunc("GET /tribes", handlers.GetTribes)
	apiRouter.HandleFunc("POST /tribes", handlers.CreateTribe)
	router.Handle("/api/", http.StripPrefix("/api", apiRouter)) // TODO authenticated

	stack := middleware.CreateStack(
		middleware.Logging,
	)

	server := http.Server{
		Addr:    fmt.Sprint(cfg.Server.Host, ":", cfg.Server.Port),
		Handler: stack(router),
	}

	slog.Info("server starting", "host", cfg.Server.Host, "port", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
