package auth

import (
	"errors"
	"fmt"
	"georgslauf/htmx"
	"log/slog"
	"net/http"

	"golang.org/x/oauth2"
)

type OAuthConfig struct {
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
	Endpoint     string `yaml:"endpoint"`
	BaseURL      string `yaml:"baseUrl"`
}

type authHandler struct {
	oauth      *oauth2.Config
	endpoint   string
	onCallback callback
}

type callback interface {
	Callback(w http.ResponseWriter, r *http.Request, accessToken string)
}

var (
	ErrorClientNil = errors.New("client is nil")
)

func NewAuthHandler(cfg OAuthConfig, onCallback callback) (*authHandler, error) {
	config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: fmt.Sprintf("%s/oauth/access_token", cfg.Endpoint),
			AuthURL:  fmt.Sprintf("%s/oauth/authorize", cfg.Endpoint),
		},
		RedirectURL: fmt.Sprintf("%s/oauth/callback", cfg.BaseURL),
	}

	h := authHandler{
		oauth:      config,
		endpoint:   cfg.Endpoint,
		onCallback: onCallback,
	}

	return &h, nil
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	htmxRequest := htmx.IsHTMX(r)

	url := h.oauth.AuthCodeURL("state", oauth2.AccessTypeOnline) // TODO (twice)

	slog.Debug("LOGIN", "htmxRequest", htmxRequest, "url", url)

	if htmxRequest {
		w.Header().Set(htmx.HeaderRedirect, url)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *authHandler) Callback(w http.ResponseWriter, r *http.Request) {
	//Use the authorization code that is pushed to the redirect
	//URL. Exchange will do the handshake to retrieve the
	//initial access token. The HTTP Client returned by
	//conf.Client will refresh the token as necessary.

	// TODO render error

	ctx := r.Context()
	htmxRequest := htmx.IsHTMX(r)

	code := r.URL.Query().Get("code")
	if code == "" {
		slog.Warn("code empty")
		return
	}

	// TODO formvalue state
	// TODO request context or background?
	// PKCE ? CSRF
	tok, err := h.oauth.Exchange(ctx, code, oauth2.AccessTypeOnline) // TODO (twice)
	if err != nil {
		slog.Warn("Exchange", "err", err)

		if htmxRequest {
			w.Header().Set(htmx.HeaderRedirect, "/") // TODO return to error page instead
			return
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect) // return to error page

		return
	}

	// TODO add refresh func

	if h.onCallback != nil {
		h.onCallback.Callback(w, r, tok.AccessToken)
	}
}
