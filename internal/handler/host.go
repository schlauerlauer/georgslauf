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

	setMd := h.md.Get()

	templates.HostSettings(htmxRequest, user, &set, schedule, categories, csrf.Token(r), setMd).Render(ctx, w)
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
