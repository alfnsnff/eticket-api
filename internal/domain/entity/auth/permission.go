package entity

type Permission struct {
	ID     uint
	RoleID uint
	Object string // e.g. "document"
	Action string // e.g. "read"
}

func (p *Permission) TableName() string {
	return "permission"
}
