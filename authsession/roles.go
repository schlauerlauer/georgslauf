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

func RoleAtLeastUser(userRole session.Role) bool {
	return userRole >= RoleUser
}

func RoleAtLeastStation(userRole session.Role) bool {
	return userRole >= RoleStation
}

func RoleAtLeastTribe(userRole session.Role) bool {
	return userRole >= RoleTribe
}

func RoleAtLeastHost(userRole session.Role) bool {
	return userRole >= RoleHost
}

func RoleAtLeastAdmin(userRole session.Role) bool {
	return userRole >= RoleAdmin
}
