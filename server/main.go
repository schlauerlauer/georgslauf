package main

import (
	"context"
	"embed"
	"fmt"
	"georgslauf/infra/persistence"
	"georgslauf/interfaces"
	"georgslauf/interfaces/config"
	"georgslauf/middleware"
	"path/filepath"

	"log/slog"
	"net/http"
	"os"
	"time"

	// "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/lmittmann/tint"
	_ "github.com/tursodatabase/go-libsql"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication"
	openid "github.com/zitadel/zitadel-go/v3/pkg/authentication/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
	// "github.com/justinas/nosurf"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

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

	ctx := context.Background()
	authRouter, err := authentication.New(ctx, zitadel.New(cfg.Config.Auth.Domain), cfg.Config.Auth.Key, openid.DefaultAuthentication(cfg.Config.Auth.ClientID, cfg.Config.Auth.CallbackURL, cfg.Config.Auth.Key))
	if err != nil {
		slog.Error("zitadel sdk could not initialize", "err", err)
		os.Exit(1)
	}
	authInterceptor := authentication.Middleware(authRouter)

	repository, err := persistence.NewRepository(&cfg.Config.Database, &embedMigrations)
	if err != nil {
		slog.Error("error connecting repository", "err", err)
		os.Exit(1)
	}

	publicService := interfaces.NewPublic(
		repository,
		authInterceptor,
	)

	privateService := interfaces.NewHost(
		repository,
		authInterceptor,
	)

	router := http.NewServeMux()

	distServer := http.FileServer(neuteredFileSystem{http.Dir("./dist")})
	router.Handle("GET /dist", http.NotFoundHandler())
	router.Handle("GET /dist/", http.StripPrefix("/dist", distServer))

	router.HandleFunc("GET /ping", publicService.Ping)
	router.HandleFunc("GET /version", publicService.Version)
	router.HandleFunc("GET /robots.txt", publicService.Robots)
	router.HandleFunc("GET /.well-known/security.txt", publicService.Security)

	router.Handle("/auth/", authRouter)

	optionalAuthStack := middleware.CreateStack(
		authInterceptor.CheckAuthentication(),
	)

	// router.Handle("GET /metrics", promhttp.Handler())
	router.Handle("/", optionalAuthStack(publicService.GetHome()))

	privateRouter := http.NewServeMux()
	privateRouter.Handle("GET /", privateService.GetHostHome())

	requiredAuthStack := middleware.CreateStack(
		authInterceptor.RequireAuthentication(),
	)

	router.Handle("/dash/", http.StripPrefix("/dash", requiredAuthStack(privateRouter)))

	stack := middleware.CreateStack(
		middleware.Logging,
	)

	server := http.Server{
		Addr:    fmt.Sprint(cfg.Config.Server.Host, ":", cfg.Config.Server.Port),
		Handler: stack(router),
	}

	slog.Info("server starting", "host", cfg.Config.Server.Host, "port", cfg.Config.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, _ := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
