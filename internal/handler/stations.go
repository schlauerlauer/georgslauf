package handler

import (
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
	"time"

	"github.com/gorilla/csrf"
)

func (h *Handler) StationPostStationRole(w http.ResponseWriter, r *http.Request) {
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

	var data postRole
	if err := h.formProcessor.Process(&data, r); err != nil {
		slog.Error("Process", "err", err)
		return // TODO
	}

	set := h.settings.Get()
	if !set.Stations.EditAccountsStation {
		if err := templates.AlertError("Rollen bearbeitung ist ausgestellt").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if err := h.checkTribeRole(user.ID, data.TribeID, acl.None); err != nil {
		if err := templates.AlertError("Keine Berechtigung").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	stationId, _, err := h.checkStationRole(user.ID, acl.Admin)
	if err != nil {
		slog.Warn("checkStationRole", "err", err)
		if err := templates.AlertError("Keine Berechtigung").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}
	if stationId != data.StationID {
		if err := templates.AlertError("Keine Berechtigung").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if count, err := h.queries.GetUserWithTribeRole(ctx, db.GetUserWithTribeRoleParams{
		ID:      data.UserID,
		TribeID: data.TribeID,
	}); err != nil || count == 0 {
		slog.Warn("GetUserWithTribeRole", "err", err)
		if err := templates.AlertError("Account gehört nicht zum Stamm").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if count, err := h.queries.CountStationRoleByUser(ctx, data.UserID); err != nil {
		slog.Warn("CountStationRoleByUser", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		if err := templates.AlertError("Ein Fehler ist aufgetreten").Render(ctx, w); err != nil {
			slog.Error("templ", "err", err)
		}
		return
	} else {
		if count > 0 {
			w.WriteHeader(http.StatusBadRequest)
			if err := templates.AlertWarning("Der Account ist bereits einem Posten zugeordnet").Render(ctx, w); err != nil {
				slog.Error("templ", "err", err)
			}
			return
		}
	}

	if err := h.queries.CreateStationRole(ctx, db.CreateStationRoleParams{
		UserID:      data.UserID,
		StationID:   data.StationID,
		StationRole: acl.ACL(data.Role),
		CreatedBy: sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		},
	}); err != nil {
		slog.Error("UpdateStationRole", "err", err)
		return // TODO
	}

	if err := templates.AlertSuccess("Berechtigung gesetzt").Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

// userRole view, tribeRole none, stationRole admin
func (h *Handler) StationPutStationRole(w http.ResponseWriter, r *http.Request) {
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
	if !set.Stations.EditAccountsStation {
		if err := templates.AlertError("Rollen bearbeitung ist ausgestellt").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	var data putRole
	if err := h.formProcessor.Process(&data, r); err != nil {
		slog.Error("Process", "err", err)
		return // TODO
	}

	if err := h.checkTribeRole(user.ID, data.TribeID, acl.None); err != nil {
		if err := templates.AlertError("Keine Berechtigung").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	stationId, _, err := h.checkStationRole(user.ID, acl.Admin)
	if err != nil {
		slog.Warn("checkStationRole", "err", err)
		if err := templates.AlertError("Keine Berechtigung").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	roleStation, err := h.queries.GetStationRoleStationWithUser(ctx, data.RoleID)
	if err != nil {
		slog.Warn("role not found", "err", err)
		return
	}

	if stationId != roleStation.StationID {
		if err := templates.AlertError("Keine Berechtigung").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if count, err := h.queries.GetUserWithTribeRole(ctx, db.GetUserWithTribeRoleParams{
		ID:      roleStation.UserID,
		TribeID: data.TribeID,
	}); err != nil || count == 0 {
		slog.Warn("GetUserWithTribeRole", "err", err)
		if err := templates.AlertError("Account gehört nicht zum Stamm").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if err := h.queries.UpdateStationRole(ctx, db.UpdateStationRoleParams{
		ID:          data.RoleID,
		StationRole: acl.ACL(data.Role),
		UpdatedBy: sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		},
		UpdatedAt: time.Now().Unix(),
	}); err != nil {
		slog.Error("UpdateStationRole", "err", err)
		return // TODO
	}

	if err := templates.AlertSuccess("Berechtigung gesetzt").Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

func (h *Handler) StationDeleteStationRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.Warn("ParseInt", "err", err)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

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
	if !set.Stations.EditAccountsStation {
		if err := templates.AlertError("Rollen bearbeitung ist ausgestellt").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
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
	if tribeId <= 0 {
		return // TODO
	}

	if err := h.checkTribeRole(user.ID, tribeId, acl.None); err != nil {
		if err := templates.AlertError("Keine Berechtigung").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	stationId, _, err := h.checkStationRole(user.ID, acl.Admin)
	if err != nil {
		slog.Warn("checkStationRole", "err", err)
		if err := templates.AlertError("Keine Berechtigung").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	roleStation, err := h.queries.GetStationRoleStation(ctx, id)
	if err != nil {
		slog.Warn("role not found", "err", err)
		return
	}

	if stationId != roleStation {
		if err := templates.AlertError("Keine Berechtigung").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if err := h.queries.DeleteStationRole(ctx, id); err != nil {
		slog.Warn("sqlc", "err", err)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if err := templates.AlertSuccess("Gelöscht").Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

func (h *Handler) GetStationAccounts(w http.ResponseWriter, r *http.Request) {
	htmxRequest := htmx.IsHTMX(r)
	if !htmxRequest {
		slog.Debug("not htmx")
		return
	}

	ctx := r.Context()

	var tribeId int64
	if query := r.URL.Query().Get("tribe"); query == "" {
		slog.Warn("tribe query empty")
		return // TODO
	} else {
		if id, err := strconv.ParseInt(query, 10, 64); err != nil {
			slog.Warn("ParseInt", "err", err)
			return // TODO
		} else {
			tribeId = id
		}
	}
	if tribeId <= 0 {
		slog.Warn("tribeId <= 0")
		return // TODO
	}

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

	if err := h.checkTribeRole(user.ID, tribeId, acl.None); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		if err := templates.AlertError("Keine Berechtigung für den Stamm").Render(ctx, w); err != nil {
			slog.Error("templ", "err", err)
		}
		return
	}

	stationId, stationRole, err := h.checkStationRole(user.ID, acl.View)
	if err != nil {
		// // no tribe access and no station roles
		if err := templates.StationPointsTab(
			templates.ErrorMessage("Du wurdest noch keinem Posten zugeordnet"),
		).Render(ctx, w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	}

	station, err := h.queries.GetStationName(ctx, stationId)
	if err != nil {
		slog.Warn("GetStationName", "err", err)
		w.WriteHeader(http.StatusNotFound)
		if err := templates.AlertError("Posten nicht gefunden").Render(ctx, w); err != nil {
			slog.Error("templ", "err", err)
		}
		return
	}

	// query single role
	var roleId sql.NullInt64
	if query := r.URL.Query().Get("id"); query != "" {
		if id, err := strconv.ParseInt(query, 10, 64); err != nil {
			slog.Warn("ParseInt", "err", err)
			return // TODO
		} else {
			roleId.Int64 = id
			roleId.Valid = true
		}
	}
	if roleId.Valid {
		userRole, err := h.queries.GetStationRoleById(ctx, roleId.Int64)
		if err != nil {
			slog.Warn("role not found", "err", err)
			return
		}

		if err := templates.StationRoleModalStation(
			userRole,
			csrf.Token(r),
		).Render(ctx, w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	}

	// query all roles
	users, err := h.queries.GetUsersByTribeRole(ctx, tribeId)
	if err != nil {
		slog.Error("sqlc", "err", err)
		return // TODO
	}

	userRoles, err := h.queries.GetStationRolesInStation(ctx, stationId)
	if err != nil {
		slog.Error("sqlc", "err", err)
		return // TODO
	}

	set := h.settings.Get()

	if err := templates.StationRolesModal(
		users,
		station,
		csrf.Token(r),
		stationRole,
		userRoles,
		set.Stations.EditAccountsStation,
	).Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

type putStationPointForm struct {
	GroupId int64 `schema:"group" validate:"gte=0"`
	Points  int64 `schema:"points" validate:"gte=0,lte=100"`
}

func (h *Handler) PutStationGroupPoint(w http.ResponseWriter, r *http.Request) {
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

	var formData putStationPointForm
	if err := h.formProcessor.Process(&formData, r); err != nil {
		slog.Error("Process", "err", err)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	stationId, _, err := h.checkStationRole(user.ID, acl.Edit)
	if err != nil {
		// // no tribe access and no station roles
		if err := templates.StationPointsTab(
			templates.ErrorMessage("Du wurdest noch keinem Posten zugeordnet"),
		).Render(ctx, w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	}

	set := h.settings.Get()

	if !set.Stations.AllowScoring {
		if err := templates.StationPointsTab(
			templates.ErrorMessage("Die Bewertungen sind ausgestellt"),
		).Render(ctx, w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	}

	if err := h.queries.UpsertPointToGroup(ctx, db.UpsertPointToGroupParams{
		CreatedBy: sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		},
		UpdatedBy: sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		},
		StationID: stationId,
		GroupID:   formData.GroupId,
		Points:    formData.Points,
	}); err != nil {
		// TODO
		slog.Error("sqlc", "err", err)
		return
	}

	if err := templates.AlertSuccess("Gespeichert").Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

// NTH duplicate
func (h *Handler) GetStationGroupPointsReload(w http.ResponseWriter, r *http.Request) {
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

	stationId, stationRole, err := h.checkStationRole(user.ID, acl.Edit)
	if err != nil {
		// // no tribe access and no station roles
		if err := templates.StationPointsTab(
			templates.ErrorMessage("Du wurdest noch keinem Posten zugeordnet"),
		).Render(ctx, w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	}

	set := h.settings.Get()

	station, err := h.queries.GetStationInfo(ctx, stationId)
	if err != nil {
		slog.Error("sqlc", "err", err)
		return // TODO
	}

	if !set.Stations.AllowScoring {
		if err := templates.PointsList(
			[]db.GetPointsToGroupsRow{},
			csrf.Token(r),
			station,
			set.Groups.ShowAbbr,
			false,
			stationRole,
			false,
		).Render(ctx, w); err != nil {
			slog.Error("templ", "err", err)
		}
		return
	}

	onlyMissing := false
	if query := r.URL.Query().Get("missing"); query == "true" {
		onlyMissing = true
	}

	// points from station to group
	points, err := h.queries.GetPointsToGroups(ctx, stationId)
	if err != nil {
		slog.Error("sqlc", "err", err)
		return // TODO
	}

	if err := templates.PointsList(
		points,
		csrf.Token(r),
		station,
		set.Groups.ShowAbbr,
		true,
		stationRole,
		onlyMissing,
	).Render(ctx, w); err != nil {
		slog.Error("templ", "err", err)
	}
}

func (h *Handler) GetStationGroupPoints(w http.ResponseWriter, r *http.Request) {
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

	stationId, stationRole, err := h.checkStationRole(user.ID, acl.Edit)
	if err != nil {
		// // no tribe access and no station roles
		if err := templates.StationPointsTab(
			templates.ErrorMessage("Du wurdest noch keinem Posten zugeordnet"),
		).Render(ctx, w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	}

	set := h.settings.Get()

	station, err := h.queries.GetStationInfo(ctx, stationId)
	if err != nil {
		slog.Error("sqlc", "err", err)
		return // TODO
	}

	if !set.Stations.AllowScoring {
		if err := templates.StationPointsTab(
			templates.PointsList(
				[]db.GetPointsToGroupsRow{},
				csrf.Token(r),
				station,
				set.Groups.ShowAbbr,
				false,
				stationRole,
				false,
			),
		).Render(ctx, w); err != nil {
			slog.Error("templ", "err", err)
		}
		return
	}

	// points from station to group
	points, err := h.queries.GetPointsToGroups(ctx, stationId)
	if err != nil {
		slog.Error("sqlc", "err", err)
		return // TODO
	}

	if err := templates.StationPointsTab(
		templates.PointsList(
			points,
			csrf.Token(r),
			station,
			set.Groups.ShowAbbr,
			true,
			stationRole,
			false,
		),
	).Render(ctx, w); err != nil {
		slog.Error("templ", "err", err)
	}
}

func (h *Handler) GetStationSettings(w http.ResponseWriter, r *http.Request) {
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

	stationRole, err := h.queries.GetStationRoleByUser(ctx, user.ID)
	if err != nil {
		// // no tribe access and no station roles
		if err := templates.ErrorPage(
			htmxRequest,
			user,
			-1,
			false,
			http.StatusUnauthorized,
			fmt.Sprintf("Du wurdest noch keinem Posten zugeordnet"),
		).Render(ctx, w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	}

	// if len(stationRoles) > 1 {
	// 	slog.Warn("TODO not yet implemented")
	// 	return // TODO
	// }
	// stationRole := stationRoles[0] // TODO query param

	if stationRole.StationRole < acl.Edit {
		slog.Debug("TODO station_role acl not implemented")
		return // TODO
	}

	station, err := h.queries.GetStation(ctx, stationRole.StationID)
	if err != nil {
		slog.Warn("GetStation", "err", err)
		return // TODO
	}

	categories, err := h.queries.GetStationCategories(ctx)
	if err != nil {
		slog.Error("GetStationCategories", "err", err)
		return // TODO
	}

	positions, err := h.queries.GetStationPositionsOpen(ctx)
	if err != nil {
		slog.Error("sqlc", "err", err)
		return // TODO
	}

	set := h.settings.Get()

	if err := templates.StationSettings(
		db.GetStationsByTribeRow{
			ID:           station.ID,
			CreatedAt:    station.CreatedAt,
			UpdatedAt:    station.UpdatedAt,
			Name:         station.Name,
			Size:         station.Size,
			CategoryID:   station.CategoryID,
			Description:  station.Description,
			Requirements: station.Requirements,
			Vegan:        station.Vegan,
			PositionID:   station.PositionID,
			PositionName: station.PositionName,
			CategoryName: station.CategoryName,
			Firstname:    station.Firstname,
			UserImage:    station.UserImage,
		},
		set.Stations,
		csrf.Token(r),
		station.TribeID, // NTH multiple stations
		categories,
		positions,
	).Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}
