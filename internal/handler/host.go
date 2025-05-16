package handler

import (
	"bytes"
	"database/sql"
	"encoding/csv"
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

func (h *Handler) GetCreateStationRoleModal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stationId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.Warn("ParseInt", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("templ", "err", err)
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

	accounts, err := h.queries.GetUsersOrdered(ctx)
	if err != nil {
		slog.Error("sqlc", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		if err := templates.AlertError("Accounts nicht gefunden").Render(ctx, w); err != nil {
			slog.Error("templ", "err", err)
		}
		return
	}

	groupedAccounts := map[string][]db.GetUsersOrderedRow{}
	for _, entry := range accounts {
		if entry.TribeName.Valid {
			groupedAccounts[entry.TribeName.String] = append(groupedAccounts[entry.TribeName.String], entry)
		} else {
			groupedAccounts["Kein Stamm"] = append(groupedAccounts["Kein Stamm"], entry)
		}
	}

	if err := templates.CreateStationRoleModal(
		groupedAccounts,
		station,
		csrf.Token(r),
	).Render(ctx, w); err != nil {
		slog.Error("templ", "err", err)
	}
}

type postRole struct {
	UserID    int64 `schema:"user" validate:"gte=0"`
	StationID int64 `schema:"station" validate:"gte=0"`
	Role      int64 `schema:"role" validate:"gte=-1,lte=3"`
	TribeID   int64 `schema:"tribe" validate:"gte=0"` // used for dash and station
}

func (h *Handler) PostStationRole(w http.ResponseWriter, r *http.Request) {
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
	if err := h.formProcessor.ProcessForm(&data, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		return // TODO
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

func (h *Handler) DeleteStationRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.Warn("ParseInt", "err", err)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
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

	if err := templates.AlertSuccess("Gelöscht").Render(ctx, w); err != nil {
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

func (h *Handler) GetStationRoleModal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		slog.Warn("ParseInt", "err", err)
		return // TODO
	}

	stationRole, err := h.queries.GetStationRoleById(ctx, id)
	if err != nil {
		slog.Warn("GetStationRoleById", "err", err)
		return // TODO
	}

	if err := templates.StationRoleModal(stationRole, csrf.Token(r)).Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

func (h *Handler) GetTribeRoleModal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

func (h *Handler) PutStationRole(w http.ResponseWriter, r *http.Request) {
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

	var data putRole
	if err := h.formProcessor.ProcessForm(&data, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		return // TODO
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

type putRole struct {
	RoleID  int64 `schema:"id" validate:"gte=0"`
	Role    int64 `schema:"role" validate:"gte=-1,lte=3"`
	TribeID int64 `schema:"tribe" validate:"gte=0"` // used for dash
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

	var data putRole
	if err := h.formProcessor.ProcessForm(&data, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		return // TODO
	}

	accepted := sql.NullInt64{}
	if data.Role >= int64(acl.None) {
		accepted.Int64 = time.Now().Unix()
		accepted.Valid = true
	}

	if err := h.queries.UpdateTribeRole(ctx, db.UpdateTribeRoleParams{
		ID:         data.RoleID,
		TribeRole:  acl.ACL(data.Role),
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

func (h *Handler) GetStationsDownload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stations, err := h.queries.GetStationsDownload(ctx)
	if err != nil {
		slog.Error("sqlc", "err", err)
		return // TODO
	}

	csvStations := [][]string{
		{"ID", "Name", "Personen", "Vegan", "Beschreibung", "Material", "Position", "Stamm", "Kategorie"},
	}

	for _, entry := range stations {
		csvStations = append(csvStations, []string{
			strconv.FormatInt(entry.ID, 10),
			entry.Name,
			strconv.FormatInt(entry.Size, 10),
			strconv.FormatInt(entry.Vegan, 10),
			entry.Description.String,
			entry.Requirements.String,
			entry.Position.String,
			entry.Tribe.String,
			entry.Category.String,
		})
	}

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	if err := writer.WriteAll(csvStations); err != nil {
		slog.Error("csv error", "err", err)
		return
	}

	writer.Flush()

	w.Header().Set("Content-Disposition", "attachment; filename=posten.csv")
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Length", strconv.FormatInt(int64(buffer.Len()), 10))
	w.Write(buffer.Bytes())
}

type putHostPointForm struct {
	GroupId   int64 `schema:"group" validate:"gte=1"`   // sqlite
	StationId int64 `schema:"station" validate:"gte=1"` // sqlite
	Points    int64 `schema:"points" validate:"gte=0,lte=100"`
}

func (h *Handler) HostPutPointsToStation(w http.ResponseWriter, r *http.Request) {
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

	var formData putHostPointForm
	if err := h.formProcessor.ProcessForm(&formData, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if err := h.queries.UpsertPointToStation(ctx, db.UpsertPointToStationParams{
		CreatedBy: sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		},
		UpdatedBy: sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		},
		GroupID:   formData.GroupId,
		StationID: formData.StationId,
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

func (h *Handler) HostPutPointsToGroup(w http.ResponseWriter, r *http.Request) {
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

	var formData putHostPointForm
	if err := h.formProcessor.ProcessForm(&formData, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
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
		StationID: formData.StationId,
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

func (h *Handler) HostGetPointsDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	queries := r.URL.Query()

	if stationStr := queries.Get("station"); stationStr != "" {
		stationId, err := strconv.ParseInt(stationStr, 10, 64)
		if err != nil {
			slog.Warn("ParseInt", "err", err)
			return
		}

		station, err := h.queries.GetStationInfo(ctx, stationId)
		if err != nil {
			slog.Error("sqlc", "err", err)
			return // TODO
		}

		set := h.settings.Get()

		// points from station to group
		points, err := h.queries.GetPointsToGroups(ctx, stationId)
		if err != nil {
			slog.Error("sqlc", "err", err)
			return // TODO
		}

		if err := templates.HostPointsToGroup(
			station,
			set.Groups.ShowAbbr,
			csrf.Token(r),
			points,
		).Render(ctx, w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	}

	if groupStr := queries.Get("group"); groupStr != "" {
		groupId, err := strconv.ParseInt(groupStr, 10, 64)
		if err != nil {
			slog.Warn("ParseInt", "err", err)
			return
		}

		group, err := h.queries.GetGroupInfo(ctx, groupId)
		if err != nil {
			slog.Warn("sqlc", "err", err)
			return // TODO
		}

		set := h.settings.Get()

		// points from station to group
		points, err := h.queries.GetPointsToStations(ctx, groupId)
		if err != nil {
			slog.Error("sqlc", "err", err)
			return // TODO
		}

		if err := templates.HostPointsToStation(
			group,
			csrf.Token(r),
			points,
			set.Stations.ShowAbbr,
		).Render(ctx, w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	}
}

// NTH duplicate
func (h *Handler) HostGetResultsGroupsDownload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	results, err := h.queries.GetResultsGroups(ctx)
	if err != nil {
		slog.Warn("sqlc", "err", err)
		return
	}

	var rank int64
	var score int64 = 9223372036854775807
	groupingRank := [4]int64{}
	groupingScore := [4]int64{9223372036854775807, 9223372036854775807, 9223372036854775807, 9223372036854775807} // NTH grouping size dependend
	resultsCalc := make([]templates.GroupResult, len(results))
	for idx, row := range results {
		tmp := int64(row.Sum.Float64)
		shared := tmp == score
		if tmp < score {
			rank += 1
			score = tmp
		}

		// set grouping
		if int(row.Grouping) >= len(groupingScore) {
			slog.Warn("grouping unknown", "group", row.ID)
			templates.AlertError("Gruppe hat eine unbekannte Stufe").Render(ctx, w)
			return
		}
		groupShared := tmp == groupingScore[int(row.Grouping)]
		if tmp < groupingScore[int(row.Grouping)] {
			groupingRank[int(row.Grouping)] += 1
			groupingScore[int(row.Grouping)] = tmp
		}

		resultsCalc[idx] = templates.GroupResult{
			Rank:        rank,
			Shared:      shared,
			GroupRank:   groupingRank[int(row.Grouping)],
			GroupShared: groupShared,
			Row:         row,
		}
	}

	csvData := [][]string{
		{"ID", "Platzierung", "Stufenplatzierung", "Name", "Summe", "Stufe", "Stamm"},
	}

	for _, entry := range resultsCalc {
		csvData = append(csvData, []string{
			strconv.FormatInt(entry.Row.ID, 10),
			strconv.FormatInt(entry.Rank, 10),
			strconv.FormatInt(entry.GroupRank, 10),
			entry.Row.Name,
			strconv.FormatFloat(entry.Row.Sum.Float64, 'f', 0, 64),
			convertGrouping(entry.Row.Grouping),
			entry.Row.Tribe.String,
		})
	}

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	if err := writer.WriteAll(csvData); err != nil {
		slog.Error("csv error", "err", err)
		return
	}

	writer.Flush()

	w.Header().Set("Content-Disposition", "attachment; filename=auswertung_gruppen.csv")
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Length", strconv.FormatInt(int64(buffer.Len()), 10))
	w.Write(buffer.Bytes())
}

func (h *Handler) HostGetResultsGroups(w http.ResponseWriter, r *http.Request) {
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

	results, err := h.queries.GetResultsGroups(ctx)
	if err != nil {
		slog.Warn("sqlc", "err", err)
		return
	}

	var rank int64
	var score int64 = 9223372036854775807
	groupingRank := [4]int64{}
	groupingScore := [4]int64{9223372036854775807, 9223372036854775807, 9223372036854775807, 9223372036854775807} // NTH grouping size dependend
	resultsCalc := make([]templates.GroupResult, len(results))
	for idx, row := range results {
		tmp := int64(row.Sum.Float64)
		shared := tmp == score
		if tmp < score {
			rank += 1
			score = tmp
		}

		// set grouping
		if int(row.Grouping) >= len(groupingScore) {
			slog.Warn("grouping unknown", "group", row.ID)
			templates.AlertError("Gruppe hat eine unbekannte Stufe").Render(ctx, w)
			return
		}
		groupShared := tmp == groupingScore[int(row.Grouping)]
		if tmp < groupingScore[int(row.Grouping)] {
			groupingRank[int(row.Grouping)] += 1
			groupingScore[int(row.Grouping)] = tmp
		}

		resultsCalc[idx] = templates.GroupResult{
			Rank:        rank,
			Shared:      shared,
			GroupRank:   groupingRank[int(row.Grouping)],
			GroupShared: groupShared,
			Row:         row,
		}
	}

	if err := templates.HostResultsGroups(
		htmxRequest,
		user,
		csrf.Token(r),
		resultsCalc,
		set.Groups.ShowAbbr,
	).Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

// NTH duplicate
func (h *Handler) HostGetResultsStationsDownload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	results, err := h.queries.GetResultsStations(ctx)
	if err != nil {
		slog.Warn("sqlc", "err", err)
		return
	}

	var rank int64
	var score int64 = 9223372036854775807
	resultsCalc := make([]templates.StationResult, len(results))
	for idx, row := range results {
		tmp := int64(row.Sum.Float64)
		shared := tmp == score
		if tmp < score {
			rank += 1
			score = tmp
		}
		resultsCalc[idx] = templates.StationResult{
			Rank:   rank,
			Row:    row,
			Shared: shared,
		}
	}

	csvData := [][]string{
		{"ID", "Platz", "Name", "Punkte", "Position", "Stamm"},
	}

	for _, entry := range resultsCalc {
		csvData = append(csvData, []string{
			strconv.FormatInt(entry.Row.ID, 10),
			strconv.FormatInt(entry.Rank, 10),
			entry.Row.Name,
			strconv.FormatFloat(entry.Row.Sum.Float64, 'f', 0, 64),
			entry.Row.Position.String,
			entry.Row.Tribe.String,
		})
	}

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	if err := writer.WriteAll(csvData); err != nil {
		slog.Error("csv error", "err", err)
		return
	}

	writer.Flush()

	w.Header().Set("Content-Disposition", "attachment; filename=auswertung_posten.csv")
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Length", strconv.FormatInt(int64(buffer.Len()), 10))
	w.Write(buffer.Bytes())
}

func (h *Handler) HostGetResultsStations(w http.ResponseWriter, r *http.Request) {
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

	results, err := h.queries.GetResultsStations(ctx)
	if err != nil {
		slog.Warn("sqlc", "err", err)
		return
	}

	var rank int64
	var score int64 = 9223372036854775807
	resultsCalc := make([]templates.StationResult, len(results))
	for idx, row := range results {
		tmp := int64(row.Sum.Float64)
		shared := tmp == score
		if tmp < score {
			rank += 1
			score = tmp
		}
		resultsCalc[idx] = templates.StationResult{
			Rank:   rank,
			Row:    row,
			Shared: shared,
		}
	}

	if err := templates.HostResultsStations(
		htmxRequest,
		user,
		csrf.Token(r),
		resultsCalc,
	).Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

func (h *Handler) HostGetPoints(w http.ResponseWriter, r *http.Request) {
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

	stations, err := h.queries.GetStationsHost(ctx)
	if err != nil {
		slog.Error("sqlc", "err", err)
		return
	}

	groups, err := h.queries.GetGroupsHost(ctx)
	if err != nil {
		slog.Error("sqlc", "err", err)
	}

	w.WriteHeader(http.StatusOK)
	if err := templates.HostPoints(
		htmxRequest,
		user,
		csrf.Token(r),
		stations,
		groups,
	).Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
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

	roleRows, err := h.queries.GetStationRoles(ctx)
	if err != nil {
		slog.Error("sqlc", "err", err)
		return // TODO
	}
	roles := make(map[int64][]db.GetStationRolesRow)
	for _, row := range roleRows {
		roles[row.StationID] = append(roles[row.StationID], row)
	}

	set := h.settings.Get()

	w.WriteHeader(http.StatusOK)
	if err := templates.HostStations(
		htmxRequest,
		user,
		stations,
		csrf.Token(r),
		summary,
		set.Stations.EnableCategories,
		roles,
	).Render(ctx, w); err != nil {
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
			return // TODO
		}

		positionName = sql.NullString{
			Valid:  true,
			String: pos,
		}
	}

	roles, err := h.queries.GetStationRoles(ctx)
	if err != nil {
		slog.Error("sqlc", "err", err)
		return // TODO
	}

	templates.HostStationUpdate(
		db.GetStationsDetailsRow{
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
		},
		set.Stations.EnableCategories,
		roles,
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

func convertGrouping(grouping int64) string {
	switch grouping {
	case 0:
		return "Wölflinge"
	case 1:
		return "Jupfis"
	case 2:
		return "Pfadis"
	case 3:
		return "Rover"
	default:
		return "Unbekannt"
	}
}

func (h *Handler) GetGroupsDownload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	groups, err := h.queries.GetGroupsDownload(ctx)
	if err != nil {
		slog.Error("sqlc", "err", err)
		return // TODO
	}

	csvGroups := [][]string{
		{"ID", "Name", "Laufgruppe", "Personen", "Vegan", "Kommentar", "Stufe", "Stamm"},
	}

	for _, entry := range groups {
		csvGroups = append(csvGroups, []string{
			strconv.FormatInt(entry.ID, 10),
			entry.Name,
			entry.Abbr.String,
			strconv.FormatInt(entry.Size.Int64, 10),
			strconv.FormatInt(entry.Vegan, 10),
			entry.Comment.String,
			convertGrouping(entry.Grouping),
			entry.Tribe.String,
		})
	}

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	if err := writer.WriteAll(csvGroups); err != nil {
		slog.Error("csv error", "err", err)
		return
	}

	writer.Flush()

	w.Header().Set("Content-Disposition", "attachment; filename=gruppen.csv")
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Length", strconv.FormatInt(int64(buffer.Len()), 10))
	w.Write(buffer.Bytes())
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
