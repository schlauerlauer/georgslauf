package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
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

type IdentityData struct {
	ID     uuid.UUID
	Traits identityTraits
}

type identityTraits struct {
	Email string
	Name  identityNames
}

type identityNames struct {
	First string
	Last  string
}

type contextKey string
var IdentityKey = contextKey("identity")

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

		id, err := uuid.Parse(session.Identity.Id)
		if err != nil {
			slog.Warn("invalid uuid identity id")
		}
		traits := (session.Identity.Traits).(map[string]interface{})
		email, ok := traits["email"].(string)
		if !ok {
			slog.Warn("no email in context")
		}
		names, ok := traits["name"].(map[string]interface{})
		if !ok {
			slog.Warn("no name map in context")
		}
		first, ok := names["first"].(string)
		if !ok {
			slog.Warn("no first name in context")
		}
		last, ok := names["last"].(string)
		if !ok {
			slog.Warn("no last name in context")
		}

		ctx := context.WithValue(r.Context(), IdentityKey, IdentityData{
			ID: id,
			Traits: identityTraits{
				Email: email,
				Name: identityNames{
					First: first,
					Last:  last,
				},
			},
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
