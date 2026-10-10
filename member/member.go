package member

import "time"

const tableAdminMember = "admin_member"

const (
	statusNormal   int16 = 1
	statusDisabled int16 = 2
)

type Member struct {
	ID       int64
	Username string
	Name     string
}

type adminMember struct {
	ID           int64      `gorm:"column:id;primaryKey"`
	Username     string     `gorm:"column:username"`
	PasswordHash string     `gorm:"column:password_hash"`
	Name         string     `gorm:"column:name"`
	LastLoginIP  string     `gorm:"column:last_login_ip"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at"`
	TokenVersion int64      `gorm:"column:token_version"`
	Status       int16      `gorm:"column:status"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
}

func (adminMember) TableName() string { return tableAdminMember }
