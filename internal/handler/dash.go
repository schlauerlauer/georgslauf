package handler

import (
	"context"
	"database/sql"
	"fmt"
	"georgslauf/acl"
	"georgslauf/htmx"
	"georgslauf/internal/db"
	"georgslauf/internal/handler/templates"
	"georgslauf/session"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/csrf"
)

func (h *Handler) GetTribeIcon(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.Warn("ParseInt", "err", err)
		return
	}

	image, err := h.queries.GetTribeIcon(ctx, id)
	if err != nil {
		slog.Warn("GetTribeIcon", "err", err)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(image)), 10))
	w.Write(image)
}

func (h *Handler) GetUserIcon(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user *session.UserData
	if userData, ok := ctx.Value(session.ContextKey).(*session.UserData); ok {
		user = userData
	} else {
		slog.Warn("not ok") // TODO
	}
	if user == nil {
		return
	}

	image, err := h.queries.GetUserIcon(ctx, user.ID)
	if err != nil {
		slog.Warn("GetUserIcon", "err", err)
		return // TODO
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(image)), 10))
	w.Write(image)
}

type putGroup struct {
	TribeId  int64  `schema:"tribe" validate:"gte=0"`
	GroupId  int64  `schema:"group" validate:"gte=0"`
	Name     string `schema:"name" validate:"required,min=3,max=30" mod:"trim,sanitize"`
	Size     int64  `schema:"size" validate:"gte=0"`
	Comment  string `schemal:"comment" validate:"max=1024" mod:"trim,sanitize"`
	Grouping int64  `schema:"grouping" validate:"gte=0,lte=3"`
}

func (h *Handler) PutGroup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user *session.UserData
	if userData, ok := ctx.Value(session.ContextKey).(*session.UserData); ok {
		user = userData
	} else {
		slog.Warn("not ok")
		// TODO redirect
		return
	}
	if user == nil {
		return
	}

	var data putGroup
	if err := h.formProcessor.ProcessForm(&data, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	set := h.settings.Get()

	if !set.Groups.AllowUpdate {
		if err := templates.AlertError("Bearbeitung ist ausgeschaltet").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if data.Size < set.Groups.Min || data.Size > set.Groups.Max {
		if err := templates.AlertError("Gruppen Größe nicht erlaubt").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	// TODO allow for stations

	tribeRole, err := h.queries.GetTribeRoleByTribe(ctx, db.GetTribeRoleByTribeParams{
		UserID:  user.ID,
		TribeID: data.TribeId,
	})
	if err != nil {
		slog.Error("GetTribeRoleByTribe", "err", err)
		if err := templates.AlertError("Keine Berechtigung").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if tribeRole < acl.Edit {
		return // TODO
	}

	updatedAt := time.Now()
	if err := h.queries.UpdateGroup(ctx, db.UpdateGroupParams{
		UpdatedAt: updatedAt.Unix(),
		Name:      data.Name,
		Comment: sql.NullString{
			String: data.Comment,
			Valid:  data.Comment != "",
		},
		UpdatedBy: sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		},
		ID: data.GroupId,
		Size: sql.NullInt64{
			Int64: data.Size,
			Valid: true,
		},
		TribeID:  data.TribeId, // just to be sure
		Grouping: data.Grouping,
	}); err != nil {
		slog.Error("UpdateGroup", "err", err)
		if err := templates.AlertError("Speichern fehlgeschlagen").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if err := templates.PutGroup(updatedAt, data.GroupId, data.Name, data.Grouping, user.Firstname).Render(ctx, w); err != nil {
		slog.Error("PutGroup", "err", err)
		return
	}
}

type postJoin struct {
	TribeId int64 `schema:"tribe" validate:"gte=0"`
}

func (h *Handler) PostJoin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user *session.UserData
	if userData, ok := ctx.Value(session.ContextKey).(*session.UserData); ok {
		user = userData
	} else {
		slog.Warn("not ok")
		return // TODO redirect?
	}
	if user == nil {
		return
	}

	var data postJoin
	if err := h.formProcessor.ProcessForm(&data, r); err != nil {
		if err := templates.AlertError("Versuch fehlgeschlagen").Render(ctx, w); err != nil {
			slog.Warn("AlertError", "err", err)
		}
		return
	}

	tribeRole, err := h.queries.GetTribeRoleByTribe(ctx, db.GetTribeRoleByTribeParams{
		UserID:  user.ID,
		TribeID: data.TribeId,
	})
	if err != nil {
		// role does not exist
		if err := h.queries.CreateTribeRole(ctx, db.CreateTribeRoleParams{
			UserID:    user.ID,
			TribeID:   data.TribeId,
			TribeRole: acl.None,
			CreatedBy: sql.NullInt64{
				Int64: user.ID,
				Valid: true,
			},
		}); err != nil {
			slog.Error("CreateTribeRole", "err", err)
			if err := templates.AlertError("Beitritt fehlgeschlagen").Render(ctx, w); err != nil {
				slog.Error("AlertWarning", "err", err)
			}
			return
		}
		slog.Debug("created tribe role", "user", user.ID, "tribe", data.TribeId)
		return
	} else {
		// role exists
		slog.Debug("role already exist", "role", tribeRole)
		if err := templates.AlertWarning("Du bist bereits in dem Stamm").Render(ctx, w); err != nil {
			slog.Error("AlertWarning", "err", err)
		}
		return
	}
}

type postGroup struct {
	TribeId  int64  `schema:"tribe" validate:"gte=0"`
	Name     string `schema:"name" validate:"required,min=3,max=30" mod:"trim,sanitize"`
	Size     int64  `schema:"size" validate:"gte=0"`
	Comment  string `schemal:"comment" validate:"max=1024" mod:"trim,sanitize"`
	Grouping int64  `schema:"grouping" validate:"gte=0,lte=3"`
}

func (h *Handler) PostGroup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user *session.UserData
	if userData, ok := ctx.Value(session.ContextKey).(*session.UserData); ok {
		user = userData
	} else {
		slog.Warn("not ok")
		return // TODO redirect?
	}
	if user == nil {
		return
	}

	set := h.settings.Get()

	if !set.Groups.AllowCreate {
		if err := templates.AlertError("Erstellung ist ausgeschaltet").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	var data postGroup
	if err := h.formProcessor.ProcessForm(&data, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if data.Size < set.Groups.Min || data.Size > set.Groups.Max {
		if err := templates.AlertError("Gruppen Größe nicht erlaubt").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	tribeRole, err := h.queries.GetTribeRoleByTribe(ctx, db.GetTribeRoleByTribeParams{
		UserID:  user.ID,
		TribeID: data.TribeId,
	})
	if err != nil {
		slog.Error("GetTribeRoleByTribe", "err", err)
		if err := templates.AlertError("Keine Berechtigung").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if tribeRole < acl.Edit {
		return // TODO
	}

	timestamp := time.Now().Unix()
	id, err := h.queries.InsertGroup(ctx, db.InsertGroupParams{
		Name: data.Name,
		Comment: sql.NullString{
			String: data.Comment,
			Valid:  data.Comment != "",
		},
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
		UpdatedBy: sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		},
		CreatedBy: sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		},
		Size: sql.NullInt64{
			Int64: data.Size,
			Valid: true,
		},
		TribeID:  data.TribeId, // just to be sure
		Grouping: data.Grouping,
	})
	if err != nil {
		slog.Error("InsertGroup", "err", err)
		if err := templates.AlertError("Erstellen fehlgeschlagen").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if err := templates.DashGroup(db.GetGroupsByTribeRow{
		ID:        id,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
		Name:      data.Name,
		Size: sql.NullInt64{
			Int64: data.Size,
			Valid: true,
		},
		Grouping: data.Grouping,
		Comment: sql.NullString{
			String: data.Comment,
			Valid:  true,
		},
		Firstname: sql.NullString{
			String: user.Firstname,
			Valid:  true,
		},
	}, csrf.Token(r), data.TribeId, set.Groups).Render(ctx, w); err != nil {
		slog.Warn("DashGroup", "err", err)
	}
}

func (h *Handler) GetNewGroup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var tribeId int64
	if query := r.URL.Query().Get("tribe"); query == "" {
		slog.Debug("no url query")
		return // TODO
	} else {
		if id, err := strconv.ParseInt(query, 10, 64); err != nil {
			slog.Debug("ParseInt", "err", err)
			return // TODO
		} else {
			tribeId = id
		}
	}

	settings := h.settings.Get()

	if err := templates.DashNewGroup(csrf.Token(r), tribeId, settings.Groups).Render(ctx, w); err != nil {
		slog.Warn("DashNewGroup", "err", err)
	}
}

func (h *Handler) DashGroups(w http.ResponseWriter, r *http.Request) {
	htmxRequest := htmx.IsHTMX(r)
	if !htmxRequest {
		// TODO allow, add full base template
		slog.Debug("not htmx")
		return
	}

	ctx := r.Context()

	var user *session.UserData
	if userData, ok := ctx.Value(session.ContextKey).(*session.UserData); ok {
		user = userData
	} else {
		slog.Warn("not ok")
		return // TODO redirect?
	}
	if user == nil {
		return
	}

	var tribeId int64
	if query := r.URL.Query().Get("tribe"); query == "" {
		slog.Debug("tribeid is empty")
		return // TODO
	} else {
		if id, err := strconv.ParseInt(query, 10, 64); err != nil {
			slog.Debug("ParseInt", "err", err)
			return // TODO
		} else {
			tribeId = id
		}
	}

	tribeRole, err := h.queries.GetTribeRoleByTribe(ctx, db.GetTribeRoleByTribeParams{
		UserID:  user.ID,
		TribeID: tribeId,
	})
	if err != nil {
		slog.Error("GetTribeRoleByTribe", "err", err)
		return // TODO
	}

	if tribeRole < acl.Edit {
		return // TODO
	}

	groups, err := h.queries.GetGroupsByTribe(ctx, tribeId)
	if err != nil {
		slog.Error("GetGroupsByTribe", "err", err)
		return // TODO
	}

	settings := h.settings.Get()

	if err := templates.DashGroups(
		tribeId,
		groups,
		settings.Groups,
		csrf.Token(r),
	).Render(ctx, w); err != nil {
		slog.Warn("DashStations", "err", err)
	}
}

func (h *Handler) DashStations(w http.ResponseWriter, r *http.Request) {
	htmxRequest := htmx.IsHTMX(r)
	if !htmxRequest {
		slog.Debug("not htmx")
		return
	}

	ctx := r.Context()

	var user *session.UserData
	if userData, ok := ctx.Value(session.ContextKey).(*session.UserData); ok {
		user = userData
	} else {
		slog.Warn("not ok")
		return // TODO redirect?
	}
	if user == nil {
		return
	}

	var tribeId int64
	if query := r.URL.Query().Get("tribe"); query == "" {
		return // TODO
	} else {
		if id, err := strconv.ParseInt(query, 10, 64); err != nil {
			return // TODO
		} else {
			tribeId = id
		}
	}

	tribeRole, err := h.queries.GetTribeRoleByTribe(ctx, db.GetTribeRoleByTribeParams{
		UserID:  user.ID,
		TribeID: tribeId,
	})
	if err != nil {
		slog.Error("GetTribeRoleByTribe", "err", err)
		return // TODO
	}

	if tribeRole < acl.Edit {
		return // TODO
	}

	stations, err := h.queries.GetStationsByTribe(ctx, tribeId)
	if err != nil {
		slog.Error("GetStationsByTribe", "err", err)
		return // TODO
	}

	if err := templates.DashStations(stations).Render(ctx, w); err != nil {
		slog.Warn("DashStations", "err", err)
	}

	// TODO admin assign users to station -> automatically grant access // role state unclaimed / claimed?
}

func (h *Handler) Dash(w http.ResponseWriter, r *http.Request) {

	// FIXME tribe role not required for stations

	htmxRequest := htmx.IsHTMX(r)
	ctx := r.Context()

	var user *session.UserData
	if userData, ok := ctx.Value(session.ContextKey).(*session.UserData); ok {
		user = userData
	} else {
		slog.Warn("not ok")
		return // TODO redirect?
	}
	if user == nil {
		return
	}

	slog.Debug("user role", "role", user.ACL.String())

	var tribeId sql.NullInt64
	if query := r.URL.Query().Get("tribe"); query != "" {
		if id, err := strconv.ParseInt(query, 10, 64); err == nil {
			tribeId = sql.NullInt64{
				Int64: id,
				Valid: true,
			}
		}
	}

	var tribeRole acl.ACL
	var hasIcon bool

	if !tribeId.Valid {
		if tmpRole, err := h.queries.GetTribeRoleWithIcon(ctx, user.ID); err != nil {
			slog.Debug("GetUserCredentials", "err", err)
			if tribeName, ok := h.createTribeRequest(ctx, user.ID, user.Email); ok {
				w.WriteHeader(http.StatusCreated)
				// NTH tribe icon größer in der mitte
				if err := templates.ErrorPage(
					htmxRequest,
					user,
					tribeId.Int64,
					hasIcon,
					http.StatusCreated,
					fmt.Sprintf("Deine Berechtigung für %s wurde angefragt!", tribeName),
				).Render(ctx, w); err != nil {
					slog.Warn("ErrorPage", "err", err)
				}
				return
			} else {
				// show tribe select page
				// TODO show select tribe modal
				tribes, err := h.queries.GetTribes(ctx)
				if err != nil {
					if err := templates.ErrorPage(
						htmxRequest,
						user,
						-1,
						false,
						http.StatusInternalServerError,
						"Serverfehler",
					).Render(ctx, w); err != nil {
						slog.Warn("ErrorPage", "err", err)
					}
					return
				}
				if err := templates.TribeRoleSelect(
					htmxRequest,
					user,
					tribes,
					csrf.Token(r),
				).Render(ctx, w); err != nil {
					slog.Warn("TribeRoleSelect", "err", err)
				}
				return
			}
		} else {
			tribeRole = tmpRole.TribeRole
			tribeId = sql.NullInt64{
				Int64: tmpRole.TribeID,
				Valid: true,
			}
			hasIcon = tmpRole.IconID.Valid
		}
	} else {
		if tmpRole, err := h.queries.GetTribeRoleByTribe(ctx, db.GetTribeRoleByTribeParams{
			UserID:  user.ID,
			TribeID: tribeId.Int64,
		}); err != nil {
			slog.Error("GetTribeRoleByTribe", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			if err := templates.ErrorPage(
				htmxRequest,
				user,
				tribeId.Int64,
				hasIcon,
				http.StatusInternalServerError,
				"Keine Berechtigung gefunden",
			).Render(ctx, w); err != nil {
				slog.Warn("ErrorPage", "err", err)
			}
			return
		} else {
			tribeRole = tmpRole
		}
	}

	if !tribeId.Valid {
		slog.Error("tribeId invalid")
		return
	}

	var isAdmin, isEdit bool

	switch tribeRole {
	case acl.Denied:
		slog.Debug("acl denied")
		w.WriteHeader(http.StatusUnauthorized)
		if err := templates.ErrorPage(
			htmxRequest,
			user,
			tribeId.Int64,
			hasIcon,
			http.StatusUnauthorized,
			"Deine Berechtigung ist abgelaufen",
		).Render(ctx, w); err != nil {
			slog.Warn("ErrorPage", "err", err)
		}
		return
	case acl.None:
		slog.Debug("acl none")
		w.WriteHeader(http.StatusUnauthorized)
		if err := templates.ErrorPage(
			htmxRequest,
			user,
			tribeId.Int64,
			hasIcon,
			http.StatusUnauthorized,
			"Deine Berechtigung muss noch bestätigt werden",
		).Render(ctx, w); err != nil {
			slog.Warn("ErrorPage", "err", err)
		}
		return
	case acl.View:
		// TODO no host view, but could use posten or group
		// showUsers = false
		templates.ErrorPage(
			htmxRequest,
			user,
			tribeId.Int64,
			hasIcon,
			http.StatusUnauthorized,
			"Deine Berechtigung reicht aktuell noch nicht aus",
		).Render(ctx, w)
		return
	case acl.Edit:
		isEdit = true
	case acl.Admin:
		isEdit = true
		isAdmin = true
	}

	if err := templates.Dash(
		htmxRequest,
		user,
		tribeId.Int64,
		hasIcon,
		isEdit,
		isAdmin,
	).Render(ctx, w); err != nil {
		slog.Error("Dash", "err", err)
	}
}

func (h *Handler) createTribeRequest(ctx context.Context, userId int64, userEmail string) (string, bool) {
	emailComponents := strings.Split(userEmail, "@")
	if len(emailComponents) < 2 {
		slog.Warn("email does not have an @ in it", "email", userEmail)
		return "", false
	}

	domain := emailComponents[len(emailComponents)-1]
	if domain == "" {
		slog.Warn("domain not set", "components", emailComponents)
		return "", false
	}

	tribe, err := h.queries.GetTribeByEmail(ctx, sql.NullString{
		String: domain,
		Valid:  true,
	})
	if err != nil {
		slog.Debug("GetTribeByEmail", "err", err)
		return "", false
	}

	slog.Debug("tribe with same domain found, creating role", "tribe", tribe)
	if err := h.queries.CreateTribeRole(ctx, db.CreateTribeRoleParams{
		UserID:    userId,
		TribeID:   tribe.ID,
		TribeRole: acl.View,
		CreatedBy: sql.NullInt64{
			Int64: userId,
			Valid: true,
		},
	}); err != nil {
		slog.Warn("CreateTribeRole", "err", err)
		return "", false
	}

	return tribe.Name, true
}
