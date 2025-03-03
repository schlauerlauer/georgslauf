package session

// NTH module

import (
	"context"
	"encoding/gob"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
)

type UserData struct {
	ID         int64
	Username   string
	Firstname  string
	Lastname   string
	Email      string
	LastUpdate time.Time
	Role       Role
	HasPicture bool
}

const (
	sessionName            = "georgslauf"
	sessionKey             = "userdata"
	ContextKey  contextKey = "session"
)

type contextKey string

type Role int64

const (
	RoleDefault Role = iota
	RoleElevated
	RoleAdmin
)

var roles = []string{"Default", "Elevated", "Admin"}

func (r Role) String() string {
	if int64(r) >= int64(len(roles)) {
		return "Invalid"
	}
	return roles[r]
}

var (
	errSessionNil = errors.New("session is nil")
)

func RoleAtLeastElevated(userRole Role) bool {
	return userRole >= RoleElevated
}

func RoleAtLeastAdmin(userRole Role) bool {
	return userRole >= RoleAdmin
}

type Session struct {
	store             *sessions.CookieStore
	errorUnauthorized templ.Component
}

func NewSessionService(hash []byte, unauthorizedComponent templ.Component) *Session {
	gob.Register(UserData{})

	store := sessions.NewCookieStore(hash)

	store.Options.HttpOnly = true
	store.Options.SameSite = http.SameSiteStrictMode

	return &Session{
		store:             store,
		errorUnauthorized: unauthorizedComponent,
	}
}

// does not check session nil
func (s *Session) SaveSession(w http.ResponseWriter, r *http.Request, userData *UserData) error {
	// TODO duplicate call in Callback
	session, err := s.store.Get(r, sessionName)
	if err != nil {
		return err
	}

	session.Values[sessionKey] = userData
	if err := s.store.Save(r, w, session); err != nil {
		return err
	}
	return nil
}

func (s *Session) GetUser(r *http.Request) (*UserData, error) {
	session, err := s.store.Get(r, sessionName)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errSessionNil
	}
	if data, ok := session.Values[sessionKey]; ok {
		if cast, ok := data.(UserData); ok {
			return &cast, nil
		}
	}
	return nil, nil
}

func (s *Session) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userData, err := s.GetUser(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			if err := s.errorUnauthorized.Render(r.Context(), w); err != nil {
				slog.Warn("ErrorUnauthorized", "err", err)
			}
		}
		ctx := context.WithValue(r.Context(), ContextKey, userData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Session) RequiredAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userData, err := s.GetUser(r)
		if userData == nil || err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			if err := s.errorUnauthorized.Render(r.Context(), w); err != nil {
				slog.Warn("ErrorUnauthorized", "err", err)
			}
		}
		ctx := context.WithValue(r.Context(), ContextKey, userData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type roleFunc func(userRole Role) bool

func (s *Session) RequireRoleFunc(roleFunc roleFunc, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userData, err := s.GetUser(r)
		if userData == nil || err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			if err := s.errorUnauthorized.Render(r.Context(), w); err != nil {
				slog.Warn("ErrorUnauthorized", "err", err)
			}
			return
		}
		if !roleFunc(userData.Role) {
			w.WriteHeader(http.StatusUnauthorized)
			if err := s.errorUnauthorized.Render(r.Context(), w); err != nil {
				slog.Warn("ErrorUnauthorized", "err", err)
			}
			return
		}
		ctx := context.WithValue(r.Context(), ContextKey, userData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
