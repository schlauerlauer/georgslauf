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

	"github.com/gorilla/csrf"
)

// FIXME

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
	if err := h.formProcessor.ProcessForm(&formData, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	stationId, err := h.checkStationRole(user.ID, acl.Edit)
	if err != nil {
		// // no tribe access and no station roles
		if err := templates.StationPointsTab(
			templates.ErrorMessage("Du wurdest noch keinem Posten zugeordnet"),
		).Render(ctx, w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	}

	slog.Debug("formD", "f", formData, "user", user.ID, "station", stationId, "err", err)

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

	slog.Debug("UPSERT")
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

	stationId, err := h.checkStationRole(user.ID, acl.Edit)
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
