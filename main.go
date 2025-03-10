package main

import (
	"embed"
	"fmt"
	"georgslauf/acl"
	"georgslauf/auth"
	"georgslauf/authsession"
	"georgslauf/internal/config"
	"georgslauf/internal/db"
	"georgslauf/internal/handler"
	"georgslauf/internal/handler/templates"
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

//go:generate sqlc generate -f ./sqlc.yaml
//go:generate templ generate -path internal/handler/templates

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

	repository, err := db.NewLibsql(&cfg.Database)
	if err != nil {
		slog.Error("error connecting repository", "err", err)
		os.Exit(1)
	}

	sessionService := session.NewSessionService(
		cfg.SessionKey,
		templates.ErrorUnauthorized(false, nil),
	)

	settings := settings.New(repository.Queries)
	templates.SetHelp(settings.Get().Help.Footer) // NTH move somewhere else

	handlers, err := handler.NewHandler(repository.Queries, sessionService, settings)
	if err != nil {
		slog.Error("NewHandler", "err", err)
		os.Exit(1)
	}

	a2s := authsession.New(
		repository.Queries,
		sessionService,
		cfg.OAuth.Endpoint,
		"/dash/",
		settings,
	)

	authHandler, err := auth.NewAuthHandler(cfg.OAuth, a2s)
	if err != nil {
		slog.Error("NewAuthHandler", "err", err)
		os.Exit(1)
	}

	router := http.NewServeMux()

	// auth
	router.HandleFunc("GET /login", authHandler.Login)
	router.HandleFunc("GET /oauth/callback", authHandler.Callback)

	// ./uploads => /res/
	upl, err := os.OpenRoot(cfg.UploadDir)
	if err != nil {
		slog.Error("os.OpenRoot", "err", err)
		os.Exit(1)
	}
	router.Handle("GET /res/", http.StripPrefix("/res", http.FileServer(neuteredFileSystem{http.FS(upl.FS())})))

	// ./resources => /dist/
	subDist, err := fs.Sub(embedRes, "resources")
	if err != nil {
		slog.Error("fs.Sub", "err", err)
		os.Exit(1)
	}
	distServer := http.FileServer(http.FS(subDist))
	router.Handle("GET /dist/", http.StripPrefix("/dist", neuter(distServer)))

	// public pages
	router.Handle("GET /{$}", sessionService.OptionalAuth(http.HandlerFunc(handlers.GetHome)))
	router.HandleFunc("GET /ping", handlers.Ping)
	router.HandleFunc("GET /version", handlers.Version)
	router.HandleFunc("GET /robots.txt", handlers.Robots)
	router.HandleFunc("GET /.well-known/security.txt", handlers.Security)

	// dash routes
	dashRouter := http.NewServeMux()
	router.Handle("GET /dash/{$}", sessionService.RequiredAuth(http.HandlerFunc(handlers.Dash))) // check for permissions
	router.Handle("POST /dash/join", sessionService.RequiredAuth(http.HandlerFunc(handlers.PostJoin)))

	dashRouter.HandleFunc("GET /stations", handlers.DashStations)
	dashRouter.HandleFunc("GET /stations/new", handlers.GetNewStation)
	dashRouter.HandleFunc("GET /groups", handlers.DashGroups)
	dashRouter.HandleFunc("GET /groups/new", handlers.GetNewGroup)
	dashRouter.HandleFunc("POST /groups", handlers.PostGroup)
	dashRouter.HandleFunc("PUT /groups", handlers.PutGroup)
	dashRouter.HandleFunc("PUT /stations", handlers.PutStation)
	dashRouter.HandleFunc("POST /stations", handlers.PostStation)
	dashRouter.HandleFunc("DELETE /groups/{id}", handlers.DeleteGroup)
	dashRouter.HandleFunc("DELETE /stations/{id}", handlers.DeleteStation)
	router.Handle("/dash/", http.StripPrefix("/dash", sessionService.RequireRoleFunc(acl.ACLViewUp, dashRouter)))

	// host routes
	// NTH acl.View auf get requests
	hostRouter := http.NewServeMux()
	hostRouter.HandleFunc("GET /{$}", handlers.GetTribes)
	hostRouter.HandleFunc("GET /users", handlers.GetUsers)
	hostRouter.HandleFunc("PUT /users/role", handlers.PutUserRole)
	hostRouter.HandleFunc("GET /tribes", handlers.GetTribes)
	hostRouter.HandleFunc("POST /tribes/icon/{id}", handlers.PostTribeIcon)
	hostRouter.HandleFunc("PUT /tribes/icon/{id}", handlers.PutTribeIcon)
	hostRouter.HandleFunc("GET /settings", handlers.GetSettings)
	hostRouter.HandleFunc("PUT /settings/groups", handlers.PutSettingsGroups)
	hostRouter.HandleFunc("PUT /settings/stations", handlers.PutSettingsStations)
	hostRouter.HandleFunc("PUT /settings/login", handlers.PutSettingsLogin)
	hostRouter.HandleFunc("PUT /settings/help", handlers.PutSettingsHelp)
	hostRouter.HandleFunc("PUT /settings/home", handlers.PutSettingsHome)
	hostRouter.HandleFunc("PUT /tribes/role", handlers.PutTribeRole)
	hostRouter.HandleFunc("GET /tribes/role", handlers.GetTribeRoleModal)
	router.Handle("/host/", http.StripPrefix("/host", sessionService.RequireRoleFunc(acl.ACLEditUp, hostRouter)))

	router.Handle("GET /icon/user", sessionService.RequiredAuth(http.HandlerFunc(handlers.GetUserIcon)))
	router.Handle("GET /icon/tribe/{id}", http.HandlerFunc(handlers.GetTribeIcon))

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
