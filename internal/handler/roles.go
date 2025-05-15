package handler

import (
	"context"
	"errors"
	"georgslauf/acl"
	"georgslauf/internal/db"
)

var (
	ErrorTribeRoleNotSufficient   = errors.New("Stammes Berechtigung nicht ausreichend")
	ErrorStationRoleNotSufficient = errors.New("Posten Berechtigung nicht ausreichend")
)

func (h *Handler) checkTribeRole(userId int64, tribeId int64, requiredRole acl.ACL) error {
	tribeRole, err := h.queries.GetTribeRoleByTribe(context.Background(), db.GetTribeRoleByTribeParams{
		UserID:  userId,
		TribeID: tribeId,
	})
	if err != nil {
		return err
	}

	if tribeRole.TribeRole < requiredRole || !tribeRole.AcceptedAt.Valid {
		return ErrorTribeRoleNotSufficient
	}

	return nil
}

func (h *Handler) checkStationRoleIsTribe(roleId int64, tribeId int64) error {
	_, err := h.queries.GetCheckStationRoleIsTribe(context.Background(), db.GetCheckStationRoleIsTribeParams{
		RoleID:  roleId,
		TribeID: tribeId,
	})
	return err
}

// single role only
func (h *Handler) checkStationRole(userId int64, requireRole acl.ACL) (int64, acl.ACL, error) {
	stationRole, err := h.queries.GetStationRoleByUser(context.Background(), userId)
	if err != nil {
		return 0, acl.None, err
	}

	if stationRole.StationRole < requireRole {
		return 0, acl.None, ErrorStationRoleNotSufficient
	}

	return stationRole.StationID, stationRole.StationRole, nil
}
