package handler

import (
	"database/sql"
	"georgslauf/htmx"
	"georgslauf/internal/db"
	"georgslauf/internal/handler/templates"
	"georgslauf/internal/settings"
	"georgslauf/session"
	"io"
	"log/slog"
	"net/http"
	"strconv"

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
		slog.Error("ParseMulitpartForm", "err", err)
		return
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

	// FIXME resize
	// img, format, err := image.Decode(bytes.NewReader(data))
	// if err != nil {
	// 	slog.Error("Decode", "err", err)
	// }
	// slog.Info("test", "f", format, "img", img) // FIXME
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
		return // TODO
	}

	prev := h.settings.Get()
	prev.Login = set
	h.settings.Set(ctx, prev, user.ID)

	if err := templates.AlertSuccess("Gespeichert").Render(ctx, w); err != nil {
		slog.Warn("AlertSuccess", "err", err)
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
		slog.Warn("sqlc", "err", err)
	}

	templates.HostSettings(htmxRequest, user, &set, schedule, categories, csrf.Token(r)).Render(ctx, w)
}

func (h *Handler) PostTribeIcon(w http.ResponseWriter, r *http.Request) {
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
		slog.Error("ParseMultipartForm", "err", err)
		return
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

	w.WriteHeader(http.StatusOK)
	if err := templates.HostUsers(htmxRequest, user, csrf.Token(r)).Render(ctx, w); err != nil {
		slog.Warn("HostUsers", "err", err)
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
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := templates.HostTribes(htmxRequest, user, tribes, csrf.Token(r)).Render(ctx, w); err != nil {
		slog.Warn("HostTribes", "err", err)
		return
	}
}
