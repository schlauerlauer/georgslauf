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
		slog.Warn("user is empty")
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

func (h *Handler) DeleteStation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.Warn("ParseInt", "err", err)
		return // TODO
	}

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

	if !set.Stations.AllowDelete {
		if err := templates.AlertError("Entfernen ist ausgeschaltet").Render(ctx, w); err != nil {
			slog.Error("templ", "err", err)
		}
		return
	}

	tribeRole, err := h.queries.GetTribeRoleByTribe(ctx, db.GetTribeRoleByTribeParams{
		UserID:  user.ID,
		TribeID: tribeId,
	})
	if err != nil {
		slog.Error("GetTribeRoleByTribe", "err", err)
		if err := templates.AlertError("Keine Berechtigung").Render(ctx, w); err != nil {
			slog.Error("templ", "err", err)
		}
		return
	}

	if tribeRole.TribeRole < acl.Edit || !tribeRole.AcceptedAt.Valid {
		if err := templates.AlertError("Keine Berechtigung").Render(ctx, w); err != nil {
			slog.Error("templ", "err", err)
		}
		return
	}

	if err := h.queries.DeleteStation(ctx, id); err != nil {
		slog.Error("DeleteStation", "err", err)
		if err := templates.AlertError("Entfernen fehlgeschlagen").Render(ctx, w); err != nil {
			slog.Error("templ", "err", err)
		}
		return
	}

	if err := templates.AlertSuccess("Posten entfernt").Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

func (h *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.Warn("ParseInt", "err", err)
		return // TODO
	}

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

	if !set.Groups.AllowDelete {
		if err := templates.AlertError("Entfernen ist ausgeschaltet").Render(ctx, w); err != nil {
			slog.Error("templ", "err", err)
		}
		return
	}

	tribeRole, err := h.queries.GetTribeRoleByTribe(ctx, db.GetTribeRoleByTribeParams{
		UserID:  user.ID,
		TribeID: tribeId,
	})
	if err != nil {
		slog.Error("GetTribeRoleByTribe", "err", err)
		if err := templates.AlertError("Keine Berechtigung").Render(ctx, w); err != nil {
			slog.Error("templ", "err", err)
		}
		return
	}

	if tribeRole.TribeRole < acl.Edit || !tribeRole.AcceptedAt.Valid {
		if err := templates.AlertError("Keine Berechtigung").Render(ctx, w); err != nil {
			slog.Error("templ", "err", err)
		}
		return
	}

	if err := h.queries.DeleteGroup(ctx, id); err != nil {
		slog.Error("DeleteGroup", "err", err)
		if err := templates.AlertError("Entfernen fehlgeschlagen").Render(ctx, w); err != nil {
			slog.Error("templ", "err", err)
		}
		return
	}

	if err := templates.AlertSuccess("Gruppe entfernt").Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

type postStation struct {
	TribeId      int64  `schema:"tribe" validate:"gte=0"`
	Name         string `schema:"name" validate:"required,min=3,max=30" mod:"trim,sanitize"`
	Size         int64  `schema:"size" validate:"gte=0"`
	Vegan        int64  `schema:"vegan" validate:"gte=0"`
	Category     int64  `schema:"category"`
	Description  string `schema:"description" validate:"max=1024" mod:"trim,sanitize"`
	Requirements string `schema:"requirements" validate:"max=1024" mod:"trim,sanitize"`
	PositionId   int64  `schema:"position" validate:"gte=0"`
}

func (h *Handler) PostStation(w http.ResponseWriter, r *http.Request) {
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

	if !set.Stations.AllowCreate {
		if err := templates.AlertError("Anmeldung ist ausgestellt").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	var data postStation
	if err := h.formProcessor.ProcessForm(&data, r); err != nil {
		slog.Error("ProcessForm", "err", err)
		if err := templates.AlertError("Falsche Eingabe").Render(ctx, w); err != nil {
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

	if tribeRole.TribeRole < acl.Edit || !tribeRole.AcceptedAt.Valid {
		if err := templates.AlertError("Keine Berechtigung").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	categories, err := h.queries.GetStationCategories(ctx)
	if err != nil {
		slog.Error("GetStationCategories", "err", err)
		return // TODO
	}

	category := db.StationCategory{}
	categoryValid := false
	if set.Stations.EnableCategories {
		cat, err := h.queries.GetStationCategory(ctx, data.Category)
		if err != nil {
			slog.Warn("GetStationCategory", "err", err)
			templates.AlertError("Posten Kategorie nicht gefunden").Render(ctx, w)
			return
		}

		count, err := h.queries.GetStationCategoryCount(ctx, sql.NullInt64{
			Int64: cat.ID,
			Valid: true,
		})
		if err != nil {
			slog.Warn("GetStationCategoryCount", "err", err)
			templates.AlertError("Kategorie Postenanzahl nicht verfügbar").Render(ctx, w)
			return
		}

		if count >= cat.Max && cat.Max != 0 {
			slog.Warn("Category full", "id", cat.ID)
			templates.AlertError("Die Posten Kategorie ist voll").Render(ctx, w)
			return
		}

		category = db.StationCategory{
			ID:   cat.ID,
			Name: cat.Name,
			// Max:  cat.Max,
		}
		categoryValid = true
	}

	var positionId sql.NullInt64
	var positionName sql.NullString
	// sqlite
	if data.PositionId > 0 {
		pos, err := h.queries.GetStationPosition(ctx, data.PositionId)
		if err != nil {
			slog.Warn("GetStationPosition", "err", err)
		} else {
			positionId = sql.NullInt64{
				Valid: true,
				Int64: data.PositionId,
			}
			positionName = sql.NullString{
				Valid:  true,
				String: pos,
			}
		}
	}

	timestamp := time.Now().Unix()
	id, err := h.queries.InsertStation(ctx, db.InsertStationParams{
		Name:    data.Name,
		Size:    data.Size,
		TribeID: data.TribeId,
		CategoryID: sql.NullInt64{
			Int64: category.ID,
			Valid: categoryValid,
		},
		Description: sql.NullString{
			String: data.Description,
			Valid:  data.Description != "",
		},
		Requirements: sql.NullString{
			String: data.Requirements,
			Valid:  data.Requirements != "",
		},
		CreatedBy: sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		},
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
		UpdatedBy: sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		},
		Vegan:      data.Vegan,
		PositionID: positionId,
	})
	if err != nil {
		slog.Error("InsertStation", "err", err)
		if err := templates.AlertError("Anmelden fehlgeschlagen").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	positions, err := h.queries.GetStationPositionsOpen(ctx)
	if err != nil {
		slog.Error("sqlc", "err", err)
		return // TODO
	}

	if err := templates.DashStation(db.GetStationsByTribeRow{
		ID:           id,
		CreatedAt:    timestamp,
		UpdatedAt:    timestamp,
		Name:         data.Name,
		PositionID:   positionId,
		PositionName: positionName,
		Size:         data.Size,
		Description: sql.NullString{
			String: data.Description,
			Valid:  data.Description != "",
		},
		Requirements: sql.NullString{
			String: data.Requirements,
			Valid:  data.Requirements != "",
		},
		CategoryID: sql.NullInt64{
			Int64: category.ID,
			Valid: category.ID > 0, // sqlite
		},
		CategoryName: sql.NullString{
			String: category.Name,
			Valid:  category.Name != "",
		},
		Firstname: sql.NullString{
			String: user.Firstname,
			Valid:  true,
		},
		UserImage: []byte{}, // NTH
		Vegan:     data.Vegan,
	}, csrf.Token(r), data.TribeId, set.Stations, categories, true, user.HasPicture, positions, true).Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
	}
}

type putStation struct {
	TribeId      int64  `schema:"tribe" validate:"gte=0"`
	StationId    int64  `schema:"station" validate:"gte=0"`
	Name         string `schema:"name" validate:"required,min=3,max=30" mod:"trim,sanitize"`
	Size         int64  `schema:"size" validate:"gte=0"`
	Vegan        int64  `schema:"vegan" validate:"gte=0"`
	Category     int64  `schema:"category"`
	Description  string `schema:"description" validate:"max=1024" mod:"trim,sanitize"`
	Requirements string `schema:"requirements" validate:"max=1024" mod:"trim,sanitize"`
	PositionId   int64  `schema:"position" validate:"gte=0"`
}

func (h *Handler) PutStation(w http.ResponseWriter, r *http.Request) {
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

	if !set.Stations.AllowUpdate {
		if err := templates.AlertError("Bearbeitung ist ausgeschaltet").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	category := sql.NullInt64{}
	categoryName := sql.NullString{}
	if set.Stations.EnableCategories {
		existing, err := h.queries.GetCategoryOfStation(ctx, data.StationId)
		if err != nil {
			slog.Warn("sqlc", "err", err)
			return // TODO
		}

		// category stays the same
		if existing.CategoryID.Valid && existing.CategoryID.Int64 == data.Category {
			category = sql.NullInt64{
				Valid: true,
				Int64: existing.CategoryID.Int64,
			}
			categoryName = existing.Name
		} else {
			// category has changed
			cat, err := h.queries.GetStationCategory(ctx, data.Category)
			if err != nil {
				slog.Warn("GetStationCategory", "err", err)
				templates.AlertError("Posten Kategorie nicht gefunden").Render(ctx, w)
				return
			}

			count, err := h.queries.GetStationCategoryCount(ctx, sql.NullInt64{
				Int64: cat.ID,
				Valid: true,
			})
			if err != nil {
				slog.Warn("GetStationCategoryCount", "err", err)
				templates.AlertError("Kategorie Postenanzahl nicht verfügbar").Render(ctx, w)
				return
			}

			if count >= cat.Max && cat.Max != 0 {
				templates.AlertError("Die Posten Kategorie ist voll").Render(ctx, w)
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

	slog.Debug("test", "tribeRole", data)

	if tribeRole.TribeRole < acl.Edit || !tribeRole.AcceptedAt.Valid {
		slog.Debug("not enough rights")
		return // TODO
	}

	updatedAt := time.Now()
	if err := h.queries.UpdateStation(ctx, db.UpdateStationParams{
		ID:        data.StationId,
		TribeID:   data.TribeId, // just to be sure
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
		PositionID: sql.NullInt64{
			Int64: data.PositionId,
			Valid: data.PositionId > 0, // sqlite
		},
	}); err != nil {
		slog.Error("UpdateStation", "err", err)
		if err := templates.AlertError("Speichern fehlgeschlagen").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if err := templates.PutStation(updatedAt, data.StationId, data.Name, user.Firstname, user.HasPicture, categoryName, set.Stations.EnableCategories).Render(ctx, w); err != nil {
		slog.Error("PutStation", "err", err)
		return
	}
}

type putGroup struct {
	TribeId  int64  `schema:"tribe" validate:"gte=0"`
	GroupId  int64  `schema:"group" validate:"gte=0"`
	Name     string `schema:"name" validate:"required,min=3,max=30" mod:"trim,sanitize"`
	Size     int64  `schema:"size" validate:"gte=0"`
	Vegan    int64  `schema:"vegan" validate:"gte=0"`
	Comment  string `schema:"comment" validate:"max=1024" mod:"trim,sanitize"`
	Grouping int64  `schema:"grouping" validate:"gte=0,lte=3"`
	Abbr     string `schema:"abbr" validate:"max=3"` // only for host
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

	if tribeRole.TribeRole < acl.Edit || !tribeRole.AcceptedAt.Valid {
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
		Vegan:    data.Vegan,
	}); err != nil {
		slog.Error("UpdateGroup", "err", err)
		if err := templates.AlertError("Speichern fehlgeschlagen").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
	}

	if err := templates.PutGroup(updatedAt, data.GroupId, data.Name, data.Grouping, user.Firstname, user.HasPicture).Render(ctx, w); err != nil {
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
		// TODO render & swap
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
	Vegan    int64  `schema:"vegan" validate:"gte=0"`
	Comment  string `schema:"comment" validate:"max=1024" mod:"trim,sanitize"`
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
		if err := templates.AlertError("Anmeldung ist ausgestellt").Render(ctx, w); err != nil {
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

	if tribeRole.TribeRole < acl.Edit || !tribeRole.AcceptedAt.Valid {
		if err := templates.AlertError("Keine Berechtigung").Render(ctx, w); err != nil {
			slog.Error("AlertError", "err", err)
		}
		return
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
		Vegan:    data.Vegan,
	})
	if err != nil {
		slog.Error("InsertGroup", "err", err)
		if err := templates.AlertError("Anmelden fehlgeschlagen").Render(ctx, w); err != nil {
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
		// UserImage: ,
		Vegan: data.Vegan,
	}, csrf.Token(r), data.TribeId, set.Groups, true, user.HasPicture).Render(ctx, w); err != nil {
		slog.Warn("DashGroup", "err", err)
	}
}

func (h *Handler) GetNewStation(w http.ResponseWriter, r *http.Request) {
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

	set := h.settings.Get()

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

	if err := templates.DashNewStation(csrf.Token(r), tribeId, set.Stations, categories, positions).Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
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

	set := h.settings.Get()

	if err := templates.DashNewGroup(csrf.Token(r), tribeId, set.Groups).Render(ctx, w); err != nil {
		slog.Warn("templ", "err", err)
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

	if tribeRole.TribeRole < acl.Edit || !tribeRole.AcceptedAt.Valid {
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
		slog.Warn("DashGroups", "err", err)
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

	if tribeRole.TribeRole < acl.Edit || !tribeRole.AcceptedAt.Valid {
		return // TODO
	}

	stations, err := h.queries.GetStationsByTribe(ctx, tribeId)
	if err != nil {
		slog.Error("GetStationsByTribe", "err", err)
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

	if err := templates.DashStations(
		tribeId,
		stations,
		set.Stations,
		csrf.Token(r),
		categories,
		positions,
	).Render(ctx, w); err != nil {
		slog.Warn("DashStations", "err", err)
	}

	// TODO admin assign users to station -> automatically grant access // role state unclaimed / claimed?
}

func (h *Handler) Dash(w http.ResponseWriter, r *http.Request) {
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
		return // TODO
	}

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
	var wasAccepted bool
	var hasIcon bool

	if !tribeId.Valid {
		if tmpRole, err := h.queries.GetTribeRoleWithIcon(ctx, user.ID); err != nil {
			slog.Debug("GetTribeRoleWithIcon", "err", err)
			if tribeName, ok := h.createTribeRequest(ctx, user.ID, user.Email); ok {
				w.WriteHeader(http.StatusCreated)
				// TODO tribe icon größer in der mitte
				// Show reload button instead of to home page
				if err := templates.ErrorPage(
					htmxRequest,
					user,
					tribeId.Int64,
					hasIcon,
					http.StatusCreated,
					fmt.Sprintf("Du wurdest zu %s zugeordnet!", tribeName),
				).Render(ctx, w); err != nil {
					slog.Warn("templ", "err", err)
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
					slog.Warn("templ", "err", err)
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
			wasAccepted = tmpRole.AcceptedAt.Valid
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
				slog.Warn("templ", "err", err)
			}
			return
		} else {
			tribeRole = tmpRole.TribeRole
			hasIcon = tmpRole.IconID.Valid
			wasAccepted = tmpRole.AcceptedAt.Valid
		}
	}

	if !tribeId.Valid {
		slog.Error("tribeId invalid")
		return
	}

	var isAdmin, isEdit bool

	slog.Debug("role", "acl", int64(tribeRole), "test", tribeRole, "accepted", wasAccepted)

	switch tribeRole {
	case acl.Denied:
		w.WriteHeader(http.StatusUnauthorized)
		if err := templates.ErrorPage(
			htmxRequest,
			user,
			tribeId.Int64,
			hasIcon,
			http.StatusUnauthorized,
			"Deine Berechtigung ist abgelaufen",
		).Render(ctx, w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	case acl.None:
		if !wasAccepted {
			w.WriteHeader(http.StatusUnauthorized)
			// TODO
			if err := templates.ErrorPage(
				htmxRequest,
				user,
				tribeId.Int64,
				hasIcon,
				http.StatusUnauthorized,
				"Dein Account muss von uns noch bestätigt werden",
			).Render(ctx, w); err != nil {
				slog.Warn("templ", "err", err)
			}
			return
		}
	case acl.View:
		// showUsers = false
		if err := templates.ErrorPage(
			htmxRequest,
			user,
			tribeId.Int64,
			hasIcon,
			http.StatusOK,
			"TODO Leserechte",
		).Render(ctx, w); err != nil {
			slog.Warn("templ", "err", err)
		}
		return
	case acl.Edit:
		isEdit = true
	case acl.Admin:
		isEdit = true
		isAdmin = true
	default:
		slog.Error("tribe role out of range", "role", int64(tribeRole))
		return // TODO
	}

	if err := templates.Dash(
		htmxRequest,
		user,
		tribeId.Int64,
		hasIcon,
		isEdit,
		isAdmin,
	).Render(ctx, w); err != nil {
		slog.Error("templ", "err", err)
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
		TribeRole: acl.None,
		AcceptedAt: sql.NullInt64{
			Int64: time.Now().Unix(),
			Valid: true,
		},
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
