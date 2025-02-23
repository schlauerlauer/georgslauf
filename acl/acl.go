package acl

type ACL int64

const (
	Denied ACL = iota - 1
	None
	View
	Edit
	Admin
)
