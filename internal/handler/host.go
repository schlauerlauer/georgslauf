package handler

import (
	"database/sql"
	"georgslauf/acl"
	"georgslauf/htmx"
	"georgslauf/internal/db"
	"georgslauf/internal/handler/templates"
	"georgslauf/internal/settings"
	"georgslauf/md"
	"georgslauf/session"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/csrf"
)

func (h *Handler) PutTribeIcon(w http.ResponseWriter, r *http.Request) {
	// TODO allow for tribes

	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.Warn("ParseInt", "err", err)
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

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		slog.Error("ParseMulitpartForm", "err", err)
		return // TODO
	}

	file, _, err := r.FormFile("icon")
	if err != nil {
		slog.Warn("FormFile", "err", err)
		return // TODO
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		slog.Warn("ReadAll", "err", err)
		return // TODO
	}

	// TODO resize
	// img, format, err := image.Decode(bytes.NewReader(data))
	// if err != nil {
	// 	slog.Error("Decode", "err", err)
	// }
	// slog.Info("test", "f", format, "img", img)
	// // png.

	if err := h.queries.UpdateTribeIcon(ctx, db.UpdateTribeIconParams{
		ID: id,
		CreatedBy: sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		},
		Image: data,
	}); err != nil {
		slog.Error("UpdateTribeIcon", "err", err)
	}

	if err := templates.TribeIcon(id, csrf.Token(r)).Render(ctx, w); err != nil {
		slog.Error("TribeIcon", "err", err)
	}
}

func (h *Handler) PutSettingsGroups(w http.ResponseWriter, r *http.Request) {
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

	var set settings.Groups
	if err := h.formProcessor.ProcessForm(&set, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return // TODO
	}

	prev := h.settings.Get()
	prev.Groups = set
	h.settings.Set(ctx, prev, user.ID)

	if err := templates.AlertSuccess("Gespeichert").Render(ctx, w); err != nil {
		slog.Warn("AlertSuccess", "err", err)
	}
}

func (h *Handler) PutSettingsStations(w http.ResponseWriter, r *http.Request) {
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

	var set settings.Stations
	if err := h.formProcessor.ProcessForm(&set, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return // TODO
	}

	prev := h.settings.Get()
	prev.Stations = set
	h.settings.Set(ctx, prev, user.ID)

	if err := templates.AlertSuccess("Gespeichert").Render(ctx, w); err != nil {
		slog.Warn("AlertSuccess", "err", err)
	}
}

func (h *Handler) GetStationCategoryNew(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := templates.HostStationCategoryNew(csrf.Token(r)).Render(ctx, w); err != nil {
		slog.Warn("HostStationCategory", "err", err)
	}
}

type upsertStationCategory struct {
	Name string `schema:"name" validate:"required,min=3,max=30" mod:"trim,sanitize"`
	Max  int64  `schema:"max" validate:"gte=0"`
}

func (h *Handler) PostStationCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var data upsertStationCategory
	if err := h.formProcessor.ProcessForm(&data, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if id, err := h.queries.InsertStationCateogy(ctx, db.InsertStationCateogyParams{
		Name: data.Name,
		Max:  data.Max,
	}); err != nil {
		slog.Error("sqlc", "err", err)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	} else {
		if err := templates.HostStationCategory(csrf.Token(r), db.GetStationCategoriesRow{
			ID:   id,
			Name: data.Name,
			Max:  data.Max,
		}).Render(ctx, w); err != nil {
			slog.Warn("templ", "err", err)
		}
	}
}

func (h *Handler) PutStationCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.Warn("ParseInt", "err", err)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	var data upsertStationCategory
	if err := h.formProcessor.ProcessForm(&data, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if err := h.queries.UpdateStationCategory(ctx, db.UpdateStationCategoryParams{
		ID:   id,
		Name: data.Name,
		Max:  data.Max,
	}); err != nil {
		slog.Error("sqlc", "err", err)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if err := templates.HostStationCategory(csrf.Token(r), db.GetStationCategoriesRow{
		ID:   id,
		Name: data.Name,
		Max:  data.Max,
	}).Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

func (h *Handler) DeleteStationCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.Warn("ParseInt", "err", err)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if err := h.queries.DeleteStationCategory(ctx, id); err != nil {
		slog.Warn("sqlc", "err", err)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if err := templates.AlertSuccess("Gespeichert").Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

type putUserRole struct {
	UserID   int64 `schema:"id" validate:"gte=0"`
	UserRole int64 `schema:"role" validate:"gte=-1,lte=2"`
}

func (h *Handler) PutUserRole(w http.ResponseWriter, r *http.Request) {
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

	var data putUserRole
	if err := h.formProcessor.ProcessForm(&data, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		return // TODO
	}

	if current, err := h.queries.GetUserRole(ctx, data.UserID); err != nil {
		slog.Error("GetUserRole", "err", err)
		return
	} else {
		if current == acl.Admin {
			return // TODO
		}
	}

	if err := h.queries.UpdateUserRole(ctx, db.UpdateUserRoleParams{
		ID:   data.UserID,
		Role: acl.ACL(data.UserRole),
	}); err != nil {
		slog.Error("UpdateUserRole", "err", err)
		return // TODO
	}

	slog.Debug("PutUserRole", "data", data)
}

func (h *Handler) GetTribeRoleModal(w http.ResponseWriter, r *http.Request) {
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

	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		slog.Warn("ParseInt", "err", err)
		return // TODO
	}

	tribeRole, err := h.queries.GetTribeRoleById(ctx, id)
	if err != nil {
		slog.Warn("GetTribeRoleById", "err", err)
		return // TODO
	}

	if err := templates.TribeRoleModal(tribeRole, csrf.Token(r)).Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

type putTribeRole struct {
	TribeRoleID int64 `schema:"id" validate:"gte=0"`
	TribeRole   int64 `schema:"role" validate:"gte=-1,lte=3"`
}

func (h *Handler) PutTribeRole(w http.ResponseWriter, r *http.Request) {
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

	var data putTribeRole
	if err := h.formProcessor.ProcessForm(&data, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		return // TODO
	}

	accepted := sql.NullInt64{}
	if data.TribeRole >= int64(acl.None) {
		accepted.Int64 = time.Now().Unix()
		accepted.Valid = true
	}

	if err := h.queries.UpdateTribeRole(ctx, db.UpdateTribeRoleParams{
		ID:         data.TribeRoleID,
		TribeRole:  acl.ACL(data.TribeRole),
		AcceptedAt: accepted,
		UpdatedBy: sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		},
		UpdatedAt: time.Now().Unix(),
	}); err != nil {
		slog.Error("UpdateTribeRole", "err", err)
		return // TODO
	}

	slog.Debug("PutTribeRole", "data", data)
	if err := templates.AlertSuccess("Berechtigung gesetzt").Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

func (h *Handler) PutSettingsHome(w http.ResponseWriter, r *http.Request) {
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

	var set md.Input
	if err := h.formProcessor.ProcessForm(&set, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		if err := templates.AlertError("Eingabe nicht richtig: Hilfe").Render(ctx, w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	}

	res, err := h.md.Update(set)
	if err != nil {
		slog.Warn("Update", "err", err)
		return // TODO
	}

	prev := h.settings.Get()
	prev.Home = set
	h.settings.Set(ctx, prev, user.ID)

	if err := templates.Md(res).Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

func (h *Handler) PutSettingsHelp(w http.ResponseWriter, r *http.Request) {
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

	var set settings.Help
	if err := h.formProcessor.ProcessForm(&set, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		if err := templates.AlertError("Eingabe nicht richtig: Hilfe").Render(ctx, w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	}

	prev := h.settings.Get()
	prev.Help = set
	h.settings.Set(ctx, prev, user.ID)

	templates.SetHelp(set.Footer)

	if err := templates.AlertSuccess("Gespeichert").Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

func (h *Handler) PutSettingsLogin(w http.ResponseWriter, r *http.Request) {
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

	var set settings.Login
	if err := h.formProcessor.ProcessForm(&set, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		if err := templates.AlertError("Eingabe nicht richtig: Anmeldung").Render(ctx, w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	}

	prev := h.settings.Get()
	prev.Login = set
	h.settings.Set(ctx, prev, user.ID)

	if err := templates.AlertSuccess("Gespeichert").Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) {
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

	set := h.settings.Get()

	schedule, err := h.queries.GetSchedule(ctx)
	if err != nil {
		slog.Warn("GetSchedule", "err", err)
		return
	}

	categories, err := h.queries.GetStationCategories(ctx)
	if err != nil {
		slog.Warn("GetStationCategories", "err", err)
		return // TODO
	}

	positions, err := h.queries.GetStationPositionsStation(ctx)

	setMd := h.md.Get()

	templates.HostSettings(htmxRequest, user, &set, schedule, categories, csrf.Token(r), setMd, positions).Render(ctx, w)
}

func (h *Handler) PostTribeIcon(w http.ResponseWriter, r *http.Request) {
	// TODO allow for tribes

	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.Warn("ParseInt", "err", err)
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

	// TODO resize instead
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		slog.Error("ParseMultipartForm", "err", err)
		return // TODO
	}

	// TODO resize image

	file, _, err := r.FormFile("icon")
	if err != nil {
		slog.Warn("FormFile", "err", err)
		return // TODO
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		slog.Warn("ReadAll", "err", err)
		return // TODO
	}

	if err := h.queries.CreateTribeIcon(ctx, db.CreateTribeIconParams{
		ID: id,
		CreatedBy: sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		},
		Image: data,
	}); err != nil {
		slog.Error("CreateTribeIcon", "err", err)
	}

	slog.Info("icon created")

	if err := templates.TribeIcon(id, csrf.Token(r)).Render(ctx, w); err != nil {
		slog.Error("TribeIcon", "err", err)
	}
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
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

	users, err := h.queries.GetUsersRoleLargerNone(ctx)
	if err != nil {
		slog.Error("GetUsersRoleLargerNone", "err", err)
	}

	requests, err := h.queries.GetUsersRoleNone(ctx)
	if err != nil {
		slog.Error("GetUsersRoleNone", "err", err)
	}

	w.WriteHeader(http.StatusOK)
	if err := templates.HostUsers(htmxRequest, user, csrf.Token(r), users, requests).Render(ctx, w); err != nil {
		slog.Warn("HostUsers", "err", err)
	}
}

func (h *Handler) GetStations(w http.ResponseWriter, r *http.Request) {
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

	stations, err := h.queries.GetStationsDetails(ctx)
	if err != nil {
		slog.Error("sqlc", "err", err)
		return // TODO
	}

	summary, err := h.queries.GetStationOverview(ctx)
	if err != nil {
		slog.Error("sqlc", "err", err)
		return // TODO
	}

	set := h.settings.Get()

	w.WriteHeader(http.StatusOK)
	if err := templates.HostStations(htmxRequest, user, stations, csrf.Token(r), summary, set.Stations.EnableCategories).Render(ctx, w); err != nil {
		slog.Warn("HostStations", "err", err)
	}
}

func (h *Handler) HostDeleteStation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.Warn("ParseInt", "err", err)
		return // TODO
	}

	if err := h.queries.DeleteStation(ctx, id); err != nil {
		slog.Error("DeleteStation", "err", err)
		if err := templates.AlertError("Entfernen fehlgeschlagen").Render(ctx, w); err != nil {
			slog.Error("templ", "err", err)
		}
		return
	}

	if err := templates.HostDeleteCloseModal("Posten entfernt").Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

func (h *Handler) HostDeleteGroup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.Warn("ParseInt", "err", err)
		return // TODO
	}

	if err := h.queries.DeleteGroup(ctx, id); err != nil {
		slog.Error("DeleteGroup", "err", err)
		if err := templates.AlertError("Entfernen fehlgeschlagen").Render(ctx, w); err != nil {
			slog.Error("templ", "err", err)
		}
		return
	}

	if err := templates.HostDeleteCloseModal("Gruppe entfernt").Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

func (h *Handler) HostPutStation(w http.ResponseWriter, r *http.Request) {
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

	var data putStation
	if err := h.formProcessor.ProcessForm(&data, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	set := h.settings.Get()

	category := sql.NullInt64{}
	categoryName := sql.NullString{}
	if set.Stations.EnableCategories {
		cat, err := h.queries.GetStationCategory(ctx, data.Category)
		if err != nil {
			slog.Warn("GetStationCategory", "err", err)
			templates.AlertError("Posten Kategorie nicht gefunden").Render(ctx, w)
			return
		}

		category = sql.NullInt64{
			Valid: true,
			Int64: cat.ID,
		}
		categoryName = sql.NullString{
			String: cat.Name,
			Valid:  true,
		}
	}

	updatedAt := time.Now()
	if err := h.queries.UpdateStationHost(ctx, db.UpdateStationHostParams{
		ID:      data.StationId,
		TribeID: data.TribeId,
		PositionID: sql.NullInt64{
			Valid: data.PositionId > 0, // sqlite
			Int64: data.PositionId,
		},
		UpdatedAt: updatedAt.Unix(),
		UpdatedBy: sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		},
		Name:       data.Name,
		Size:       data.Size,
		CategoryID: category,
		Description: sql.NullString{
			String: data.Description,
			Valid:  data.Description != "",
		},
		Requirements: sql.NullString{
			String: data.Requirements,
			Valid:  data.Requirements != "",
		},
		Vegan: data.Vegan,
	}); err != nil {
		slog.Error("UpdateStation", "err", err)
		if err := templates.AlertError("Speichern fehlgeschlagen").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	tribe, err := h.queries.GetTribeNameIcon(ctx, data.TribeId)
	if err != nil {
		slog.Info("sqlc", "err", err)
	}

	var positionName sql.NullString
	// sqlite
	if data.PositionId > 0 {
		pos, err := h.queries.GetStationPosition(ctx, data.PositionId)
		if err != nil {
			slog.Error("sqlc", "err", err)
			return
		}

		positionName = sql.NullString{
			Valid:  true,
			String: pos,
		}
	}

	templates.HostStationUpdate(db.GetStationsDetailsRow{
		ID:           data.StationId,
		Name:         data.Name,
		CategoryName: categoryName,
		Tribe: sql.NullString{
			String: tribe.Name,
			Valid:  true,
		},
		Size:         data.Size,
		TribeIcon:    tribe.TribeIcon,
		PositionName: positionName,
	}, set.Stations.EnableCategories,
	).Render(ctx, w)
}

func (h *Handler) HostPutGroup(w http.ResponseWriter, r *http.Request) {
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

	updatedAt := time.Now()
	if err := h.queries.UpdateGroupHost(ctx, db.UpdateGroupHostParams{
		TribeID: data.TribeId,
		Abbr: sql.NullString{
			String: data.Abbr,
			Valid:  data.Abbr != "",
		},
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
		Grouping: data.Grouping,
		Vegan:    data.Vegan,
	}); err != nil {
		slog.Error("UpdateGroup", "err", err)
		if err := templates.AlertError("Speichern fehlgeschlagen").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	tribe, err := h.queries.GetTribeNameIcon(ctx, data.TribeId)
	if err != nil {
		slog.Info("sqlc", "err", err)
	}

	templates.HostGroupUpdate(db.GetGroupsDetailsRow{
		ID:       data.GroupId,
		Name:     data.Name,
		Grouping: data.Grouping,
		Tribe: sql.NullString{
			String: tribe.Name,
			Valid:  true,
		},
		Size: sql.NullInt64{
			Int64: data.Size,
			Valid: true,
		},
		TribeIcon: tribe.TribeIcon,
		Abbr: sql.NullString{
			String: data.Abbr,
			Valid:  data.Abbr != "",
		},
	}).Render(ctx, w)
}

func (h *Handler) GetStation(w http.ResponseWriter, r *http.Request) {
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

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return // TODO
	}

	station, err := h.queries.GetStation(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
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
		return
	}

	set := h.settings.Get()

	self := station.UpdatedBy.Valid && station.UpdatedBy.Int64 == user.ID

	tribes, err := h.queries.GetTribesName(ctx)
	if err != nil {
		slog.Error("GetTribesName", "err", err)
		return // TODO
	}

	w.WriteHeader(http.StatusOK)
	if err := templates.HostStationModal(
		station,
		csrf.Token(r),
		set.Stations.EnableCategories,
		self,
		user.HasPicture,
		categories,
		tribes,
		positions,
	).Render(ctx, w); err != nil {
		slog.Error("templ", "err", err)
	}
}

func (h *Handler) GetGroup(w http.ResponseWriter, r *http.Request) {
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

	groupId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return // TODO
	}

	group, err := h.queries.GetGroup(ctx, groupId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return // TODO
	}

	set := h.settings.Get()

	self := group.UpdatedBy.Valid && group.UpdatedBy.Int64 == user.ID

	tribes, err := h.queries.GetTribesName(ctx)
	if err != nil {
		slog.Error("GetTribesName", "err", err)
		return // TODO
	}

	w.WriteHeader(http.StatusOK)
	if err := templates.HostGroupModal(
		group,
		set.Groups,
		csrf.Token(r),
		self,
		user.HasPicture,
		tribes,
	).Render(ctx, w); err != nil {
		slog.Error("templ", "err", err)
	}
}

func (h *Handler) GetGroupsAbbr(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	abbrs, err := h.queries.GetGroupsAbbr(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("sqlc", "err", err)
		return // TODO
	}

	x := firstMissing(abbrs)

	templates.FirstValid(x).Render(ctx, w)
}

func firstMissing(ordered []int64) int64 {
	for idx, val := range ordered {
		if val != int64(idx)+1 {
			return int64(idx) + 1
		}
	}

	return int64(len(ordered) + 1)
}

func (h *Handler) GetGroups(w http.ResponseWriter, r *http.Request) {
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

	groups, err := h.queries.GetGroupsDetails(ctx)
	if err != nil {
		slog.Error("sqlc", "err", err)
		return // TODO
	}

	summary, err := h.queries.GetGroupOverview(ctx)
	if err != nil {
		slog.Error("sqlc", "err", err)
		return // TODO
	}

	w.WriteHeader(http.StatusOK)
	if err := templates.HostGroups(htmxRequest, user, groups, csrf.Token(r), summary).Render(ctx, w); err != nil {
		slog.Warn("HostGroups", "err", err)
	}
}

func (h *Handler) GetTribes(w http.ResponseWriter, r *http.Request) {
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

	tribes, err := h.queries.GetTribes(ctx)
	if err != nil {
		slog.Warn("GetTribes", "err", err)
		return // TODO
	}

	tribeRoles, err := h.queries.GetTribeRolesOpen(ctx)
	if err != nil {
		slog.Error("GetTribeRolesOpen", "err", err)
		return // TODO
	}

	accountsRows, err := h.queries.GetTribeRolesAssigned(ctx)
	if err != nil {
		slog.Error("GetTribeRolesAssigned", "err", err)
		return // TODO
	}
	accounts := make(map[int64][]db.GetTribeRolesAssignedRow)
	for _, row := range accountsRows {
		accounts[row.TribeID] = append(accounts[row.TribeID], row)
	}

	groupRows, err := h.queries.GetGroupsHost(ctx)
	if err != nil {
		slog.Error("GetGroupsHost", "err", err)
		return // TODO
	}
	groups := make(map[int64][]db.GetGroupsHostRow)
	for _, row := range groupRows {
		groups[row.TribeID] = append(groups[row.TribeID], row)
	}

	stationRows, err := h.queries.GetStationsHost(ctx)
	if err != nil {
		slog.Error("GetStationsHost", "err", err)
		return // TODO
	}
	stations := make(map[int64][]db.GetStationsHostRow)
	for _, row := range stationRows {
		stations[row.TribeID] = append(stations[row.TribeID], row)
	}

	w.WriteHeader(http.StatusOK)
	if err := templates.HostTribes(
		htmxRequest,
		user,
		tribes,
		csrf.Token(r),
		tribeRoles,
		accounts,
		groups,
		stations,
	).Render(ctx, w); err != nil {
		slog.Warn("HostTribes", "err", err)
		return
	}
}
