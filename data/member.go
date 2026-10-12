package data

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

const tableAdminMember = "admin_member"

const (
	StatusNormal   int16 = 1
	StatusDisabled int16 = 2
)

var (
	ErrNotFound      = errors.New("管理员不存在")
	ErrInvalidDriver = errors.New("数据库驱动无效")
)

type AdminMember struct {
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

func (AdminMember) TableName() string { return tableAdminMember }

const postgresCreateTableSQL = `CREATE TABLE IF NOT EXISTS admin_member (
  id bigserial PRIMARY KEY,
  username varchar(128) NOT NULL,
  password_hash varchar(128) NOT NULL,
  name varchar(32) NOT NULL DEFAULT '',
  last_login_ip varchar(64) NOT NULL DEFAULT '',
  last_login_at timestamptz NULL,
  token_version bigint NOT NULL DEFAULT 0,
  status smallint NOT NULL DEFAULT 1,
  deleted_at timestamptz NULL,
  created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT uk_admin_member_username UNIQUE (username)
)`

const mysqlCreateTableSQL = "CREATE TABLE IF NOT EXISTS admin_member (" +
	"id bigint NOT NULL AUTO_INCREMENT," +
	"username varchar(128) NOT NULL," +
	"password_hash varchar(128) NOT NULL," +
	"name varchar(32) NOT NULL DEFAULT ''," +
	"last_login_ip varchar(64) NOT NULL DEFAULT ''," +
	"last_login_at datetime NULL," +
	"token_version bigint NOT NULL DEFAULT 0," +
	"status smallint NOT NULL DEFAULT 1," +
	"deleted_at datetime NULL," +
	"created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP," +
	"updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
	"PRIMARY KEY (id)," +
	"UNIQUE KEY uk_admin_member_username (username)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci"

func EnsureAdminMember(ctx context.Context, db *gorm.DB, driver string) error {
	if db.WithContext(ctx).Migrator().HasTable(tableAdminMember) {
		return nil
	}
	var sqlText string
	switch driver {
	case "postgres":
		sqlText = postgresCreateTableSQL
	case "mysql":
		sqlText = mysqlCreateTableSQL
	default:
		return ErrInvalidDriver
	}
	return db.WithContext(ctx).Exec(sqlText).Error
}

func CountAdminMember(ctx context.Context, db *gorm.DB) (int64, error) {
	var count int64
	err := db.WithContext(ctx).Model(&AdminMember{}).Count(&count).Error
	return count, err
}

func InsertAdminMember(ctx context.Context, db *gorm.DB, row *AdminMember) error {
	return db.WithContext(ctx).Create(row).Error
}

func TakeActiveAdminMember(ctx context.Context, db *gorm.DB, query string, arg any) (*AdminMember, error) {
	var row AdminMember
	err := db.WithContext(ctx).Where(query, arg).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func UpdateAdminMemberLogin(ctx context.Context, db *gorm.DB, id int64, ip string, at time.Time) error {
	return db.WithContext(ctx).Model(&AdminMember{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{"last_login_ip": ip, "last_login_at": at}).Error
}

func BumpAdminMemberTokenVersion(ctx context.Context, db *gorm.DB, id, version int64) (int64, error) {
	res := db.WithContext(ctx).Model(&AdminMember{}).
		Where("id = ? AND deleted_at IS NULL AND token_version = ?", id, version).
		Update("token_version", gorm.Expr("token_version + 1"))
	return res.RowsAffected, res.Error
}

func SoftDeleteAdminMember(ctx context.Context, db *gorm.DB, id int64, username string, at time.Time) (int64, error) {
	res := db.WithContext(ctx).Model(&AdminMember{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{"username": username, "deleted_at": at})
	return res.RowsAffected, res.Error
}
