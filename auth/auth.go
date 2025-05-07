package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

type OAuthConfig struct {
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
	Endpoint     string `yaml:"endpoint"`
	BaseURL      string `yaml:"baseUrl"`
	Hash         []byte `yaml:"hash"`
}

type authHandler struct {
	oauth      *oauth2.Config
	endpoint   string
	onCallback callback
	store      *sessions.CookieStore
}

type callback interface {
	Callback(w http.ResponseWriter, r *http.Request, accessToken string)
}

var (
	errSessionNil = errors.New("session is nil")
)

const (
	sessionName = "georgslauf.oauth"
	sessionKey  = "state"
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

	store := sessions.NewCookieStore(cfg.Hash)

	store.Options.HttpOnly = true
	store.Options.SameSite = http.SameSiteLaxMode
	store.Options.Secure = true
	store.MaxAge(60 * 5)

	h := authHandler{
		oauth:      config,
		endpoint:   cfg.Endpoint,
		onCallback: onCallback,
		store:      store,
	}

	return &h, nil
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	state, err := generateState()
	if err != nil {
		slog.Error("generateState", "err", err)
		return // TODO
	}

	if err := h.saveSession(w, r, state); err != nil {
		slog.Error("saveSession", "err", err)
		return // TODO
	}

	url := h.oauth.AuthCodeURL(state, oauth2.AccessTypeOnline)

	http.Redirect(w, r, url, http.StatusFound)
}

func generateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (h *authHandler) saveSession(w http.ResponseWriter, r *http.Request, state string) error {
	session, err := h.store.Get(r, sessionName)
	if err != nil {
		return err
	}
	if session == nil {
		return errSessionNil
	}

	session.Values[sessionKey] = state
	if err := h.store.Save(r, w, session); err != nil {
		return err
	}
	return nil
}

func (h *authHandler) getSession(r *http.Request) (string, error) {
	session, err := h.store.Get(r, sessionName)
	if err != nil {
		return "", err
	}
	if session == nil {
		return "", errSessionNil
	}

	if data, ok := session.Values[sessionKey]; ok {
		if cast, ok := data.(string); ok {
			return cast, nil
		} else {
			slog.Warn("session data cast not ok")
		}
	} else {
		slog.Warn("session does not exist")
	}

	return "", nil
}

func (h *authHandler) Callback(w http.ResponseWriter, r *http.Request) {
	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.

	state, err := h.getSession(r)
	if err != nil {
		slog.Error("getSession", "err", err)
		return // TODO
	}

	queryState := r.URL.Query().Get("state")

	if state == "" || queryState == "" {
		slog.Warn("state is empty", "state", state, "query", queryState)
		return
	}

	if state != queryState {
		slog.Warn("state does not match query", "state", state, "query", queryState)
		return // TODO
	}

	ctx := r.Context()

	code := r.URL.Query().Get("code")
	if code == "" {
		slog.Warn("code empty")
		return // TODO
	}

	// NTH PKCE
	tok, err := h.oauth.Exchange(ctx, code, oauth2.AccessTypeOnline)
	if err != nil {
		slog.Warn("Exchange", "err", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect) // TODO return to error page
		return
	}

	// NTH add refresh func

	if h.onCallback != nil {
		h.onCallback.Callback(w, r, tok.AccessToken)
	}
}
