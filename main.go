package main

import (
	"context"
	"embed"
	"fmt"
	"georgslauf/auth"
	"georgslauf/authsession"
	"georgslauf/internal/config"
	"georgslauf/internal/db"
	"georgslauf/internal/handler"
	"georgslauf/internal/settings"
	"georgslauf/session"
	"io/fs"

	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/csrf"
	"github.com/lmittmann/tint"
	"github.com/schlauerlauer/go-middleware"
)

const (
	csrfCookieName = "georgslauf.csrf"
)

//go:embed all:resources
var embedRes embed.FS

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

	repository, err := db.NewRepository(&cfg.Database)
	if err != nil {
		slog.Error("error connecting repository", "err", err)
		os.Exit(1)
	}

	sessionService := session.NewSessionService(cfg.SessionKey)

	settings := settings.New(repository.Queries)

	handlers, err := handler.NewHandler(repository.Queries, sessionService, settings)
	if err != nil {
		slog.Error("NewHandler", "err", err)
		os.Exit(1)
	}

	a2s := authsession.New(repository.Queries, sessionService, cfg.OAuth.Endpoint, "/dash/")

	authHandler, err := auth.NewAuthHandler(cfg.OAuth, a2s)
	if err != nil {
		slog.Error("NewAuthHandler", "err", err)
		os.Exit(1)
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /debug", func(w http.ResponseWriter, r *http.Request) {
		set := settings.Get()
		set.Groups.Min = 4
		set.Groups.Max = 13
		settings.Set(context.Background(), set)
	})
	router.HandleFunc("GET /debug2", func(w http.ResponseWriter, r *http.Request) {
		set := settings.Get()
		set.Groups.Min = 2
		set.Groups.Max = 8
		settings.Set(context.Background(), set)
	})

	// auth
	router.HandleFunc("GET /login", authHandler.Login)
	router.HandleFunc("GET /oauth/callback", authHandler.Callback)

	// ./uploads => /res/
	resRouter := http.FileServer(neuteredFileSystem{http.Dir(cfg.UploadDir)})
	router.Handle("GET /res/", http.StripPrefix("/res", resRouter))

	// ./resources => /dist/
	subDist, err := fs.Sub(embedRes, "resources")
	if err != nil {
		slog.Error("fs.Sub", "err", err)
		os.Exit(1)
	}
	distServer := http.FileServer(http.FS(subDist))
	router.Handle("GET /dist/", http.StripPrefix("/dist", neuter(distServer)))

	// public pages
	router.Handle("/", sessionService.OptionalAuth(http.HandlerFunc(handlers.GetHome)))
	router.HandleFunc("GET /ping", handlers.Ping)
	router.HandleFunc("GET /version", handlers.Version)
	router.HandleFunc("GET /robots.txt", handlers.Robots)
	router.HandleFunc("GET /.well-known/security.txt", handlers.Security)

	// dash routes
	dashRouter := http.NewServeMux()
	dashRouter.HandleFunc("GET /", handlers.Dash)
	dashRouter.HandleFunc("GET /stations", handlers.DashStations)
	dashRouter.HandleFunc("GET /groups", handlers.DashGroups)
	dashRouter.HandleFunc("PUT /groups", handlers.PutGroup)
	router.Handle("/dash/", http.StripPrefix("/dash", sessionService.RequiredAuth(dashRouter)))

	// host routes
	hostRouter := http.NewServeMux()
	// hostRouter.HandleFunc("GET /", handlers.GetHostHome) // TODO
	// hostRouter.HandleFunc("GET /schedule", handlers.GetSchedule)
	router.Handle("/host/", http.StripPrefix("/host", sessionService.RequireRoleFunc(session.RoleAtLeastElevated, hostRouter)))

	router.Handle("GET /icon/user", sessionService.RequiredAuth(http.HandlerFunc(handlers.GetUserIcon)))
	router.Handle("GET /icon/tribe", sessionService.RequiredAuth(http.HandlerFunc(handlers.GetTribeIcon)))

	stack := middleware.CreateStack(
		middleware.Logging,
		csrf.Protect(
			cfg.CsrfKey,
			csrf.Secure(true),
			csrf.SameSite(csrf.SameSiteStrictMode),
			csrf.Path("/"),
			csrf.CookieName(csrfCookieName),
		),
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
