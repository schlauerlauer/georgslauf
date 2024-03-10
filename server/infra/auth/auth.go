package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	ory "github.com/ory/client-go"
)

type AuthInterface interface {
	SessionMiddleware(next http.Handler) http.Handler
}

type ClientData struct {
	client *ory.APIClient
}

var _ AuthInterface = &ClientData{}

func NewAuth(client *ory.APIClient) *ClientData {
	return &ClientData{client: client}
}

func (cd *ClientData) validateSession(request *http.Request) (*ory.Session, error) {
	cookie, err := request.Cookie("georgslauf_session")
	if err != nil {
		return nil, err
	}
	if cookie == nil {
		return nil, errors.New("no session found in cookie")
	}

	resp, _, err := cd.client.FrontendAPI.ToSession(context.Background()).Cookie(cookie.String()).Execute()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (cd *ClientData) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := cd.validateSession(r)
		if err != nil {
			slog.Warn("error validating session", "err", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if session == nil || !*session.Active {
			slog.Warn("session nil or not active")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		slog.Info("context", "rq", r.Context())

		next.ServeHTTP(w, r)
	})
}
