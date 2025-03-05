package session

// NTH module

import (
	"context"
	"encoding/gob"
	"errors"
	"georgslauf/acl"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
)

type UserData struct {
	ID         int64
	Username   string
	Firstname  string
	Lastname   string
	Email      string
	ACL        acl.ACL
	HasPicture bool
}

const (
	sessionName            = "georgslauf"
	sessionKey             = "userdata"
	ContextKey  contextKey = "session"
)

type contextKey string

var (
	errSessionNil = errors.New("session is nil")
)

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
			slog.Warn("GetUser", "err", err)
			w.WriteHeader(http.StatusUnauthorized)
			if err := s.errorUnauthorized.Render(r.Context(), w); err != nil {
				slog.Warn("ErrorUnauthorized", "err", err)
			}
		}
		ctx := context.WithValue(r.Context(), ContextKey, userData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Session) RequireRoleFunc(roleFunc acl.RoleFunc, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userData, err := s.GetUser(r)
		if userData == nil || err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			if err := s.errorUnauthorized.Render(r.Context(), w); err != nil {
				slog.Warn("ErrorUnauthorized", "err", err)
			}
			return
		}
		if !roleFunc(userData.ACL) {
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
