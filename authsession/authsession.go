package authsession

import (
	"database/sql"
	"georgslauf/acl"
	"georgslauf/internal/db"
	"georgslauf/internal/handler/templates"
	"georgslauf/internal/settings"
	"georgslauf/mattermost"
	"georgslauf/session"
	"log/slog"
	"net/http"
	"time"
)

type authToSession struct {
	queries     *db.Queries
	client      *mattermost.Client
	session     *session.Session
	redirectUrl string
	settings    *settings.SettingsService
}

func New(
	queries *db.Queries,
	session *session.Session,
	authEndpoint string,
	redirectUrl string,
	settings *settings.SettingsService,
) *authToSession {
	client := mattermost.NewClient(authEndpoint)

	return &authToSession{
		client:      client,
		queries:     queries,
		session:     session,
		redirectUrl: redirectUrl,
		settings:    settings,
	}
}

// TODO add callback functions for different auth errors, connecting them to templates from here

func (a2s *authToSession) Callback(w http.ResponseWriter, r *http.Request, token string) {
	// sessionData, err := a2s.session.GetUser(r)
	// if err != nil {
	// 	slog.Warn("sesion.GetUser", "err", err)
	// 	return // TODO render
	// }

	userInfo, err := a2s.client.GetUser(token)
	if err != nil {
		// NTH add var error returns in client to render different errors here
		slog.Error("GetUser", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		if err := templates.ErrorPageSlim(
			http.StatusBadRequest,
			"User info Fehler",
		).Render(r.Context(), w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	}
	if userInfo == nil {
		slog.Error("userInfo is nil")
		w.WriteHeader(http.StatusBadRequest)
		if err := templates.ErrorPageSlim(
			http.StatusBadRequest,
			"User info nicht erhalten",
		).Render(r.Context(), w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	}
	if !userInfo.EmailVerified {
		slog.Error("userInfo.EmailVerified false", "userInfo", userInfo)
		w.WriteHeader(http.StatusUnauthorized)
		if err := templates.ErrorPageSlim(
			http.StatusUnauthorized,
			"Email ist nicht verifiziert. Bitte aus Mattermost abmelden und neu anmelden, um die Bestätigungsmail zu erhalten",
		).Render(r.Context(), w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	}
	if userInfo.ID == "" {
		slog.Error("userInfo.ID empty", "userInfo", userInfo)
		w.WriteHeader(http.StatusBadRequest)
		if err := templates.ErrorPageSlim(
			http.StatusBadRequest,
			"User ID ist leer",
		).Render(r.Context(), w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	}

	ctx := r.Context()
	existing, err := a2s.queries.GetUserIdByExt(ctx, sql.NullString{
		String: userInfo.ID,
		Valid:  true,
	})
	if err != nil {
		// user does not exist, yet
		existing.ID, err = a2s.queries.CreateUser(ctx, db.CreateUserParams{
			ExtID: sql.NullString{
				String: userInfo.ID,
				Valid:  true,
			},
			Username:  userInfo.Username,
			Firstname: userInfo.Firstname,
			Lastname:  userInfo.Lastname,
			Email:     userInfo.Email,
		})
		if err != nil {
			slog.Error("CreateUser", "err", err)
			return // TODO render
		}
	} else {
		err = a2s.queries.UpdateUser(ctx, db.UpdateUserParams{
			ID:        existing.ID,
			LastLogin: time.Now().Unix(),
			Username:  userInfo.Username,
			Firstname: userInfo.Firstname,
			Lastname:  userInfo.Lastname,
			Email:     userInfo.Email, // NTH allow email change?
		})
		if err != nil {
			slog.Error("UpdateLogin", "err", err)
			return // TODO render
		}
	}

	hasPicture := false
	slog.Debug("picture", "update", userInfo.PictureUpdate, "lastLogin", existing.LastLogin)
	if userInfo.PictureUpdate > existing.LastLogin {
		slog.Debug("updating user icon")
		image, err := a2s.client.GetUserPicture(token, userInfo.ID)
		if err != nil {
			slog.Warn("GetUserPicture", "err", err)
		} else {
			slog.Debug("got image", "len", len(image))
			if err := a2s.queries.CreateUserIcon(ctx, db.CreateUserIconParams{
				ID:    existing.ID,
				Image: image,
			}); err != nil {
				slog.Warn("CreateUserIcon", "err", err)
			} else {
				slog.Debug("CreateUserIcon success", "id", existing.ID)
				hasPicture = true
			}
		}
	}

	sessionData := &session.UserData{
		ID:         existing.ID,
		Username:   userInfo.Username,
		Firstname:  userInfo.Firstname,
		Lastname:   userInfo.Lastname,
		Email:      userInfo.Email,
		ACL:        acl.ACL(existing.Role),
		HasPicture: hasPicture,
	}

	if err := a2s.session.SaveSession(w, r, sessionData); err != nil {
		slog.Error("SaveSession", "err", err)
		return // TODO render
	}

	set := a2s.settings.Get()

	if err := templates.Login(sessionData, set.Login).Render(ctx, w); err != nil {
		slog.Warn("Login", "err", err)
	}
}
