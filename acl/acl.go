package acl

type ACL int64

const (
	Denied ACL = iota - 1
	None
	View
	Edit
	Admin
)

type RoleFunc func(role ACL) bool

var acl = []string{"Denied", "None", "View", "Edit", "Admin"}

func (a ACL) String() string {
	if int64(a) >= int64(len(acl)) {
		return "Invalid"
	}
	return acl[a]
}

func ACLViewUp(role ACL) bool {
	return role >= View
}

func ACLEditUp(role ACL) bool {
	return role >= Edit
}
