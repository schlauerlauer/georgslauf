package authsession

import "georgslauf/session"

const (
	RoleDefault session.Role = iota
	RoleUser
	RoleStation
	RoleTribe
	RoleHost
	RoleAdmin
)

func RoleAtLeastStation(userRole session.Role) bool {
	return userRole >= RoleStation
}

func RoleAtLeastHost(userRole session.Role) bool {
	return userRole >= RoleHost
}
