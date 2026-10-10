package member

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const maxUsernameLen = 64

type InitConfig struct {
	Driver string
	DSN    string
}

type InitResult struct {
	AdminCreated bool
	Username     string
	Password     string
}

type Service struct {
	db     *gorm.DB
	driver string
}

func Init(ctx context.Context, cfg InitConfig) (*Service, InitResult, func(), error) {
	if err := validateInit(cfg); err != nil {
		return nil, InitResult{}, nil, err
	}
	db, err := openDB(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, InitResult{}, nil, ErrInternal
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, InitResult{}, nil, ErrInternal
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err = sqlDB.PingContext(pingCtx); err != nil {
		_ = sqlDB.Close()
		return nil, InitResult{}, nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	svc := &Service{db: db, driver: cfg.Driver}
	result, err := svc.ensureSchema(ctx)
	if err != nil {
		_ = sqlDB.Close()
		return nil, InitResult{}, nil, err
	}
	return svc, result, func() { _ = sqlDB.Close() }, nil
}

func validateInit(cfg InitConfig) error {
	if cfg.Driver != "mysql" && cfg.Driver != "postgres" {
		return ErrInvalidConfig
	}
	if cfg.DSN == "" {
		return ErrInvalidConfig
	}
	return nil
}

func openDB(driver, dsn string) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch driver {
	case "mysql":
		dialector = gormmysql.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	default:
		return nil, ErrInvalidConfig
	}
	return gorm.Open(dialector, &gorm.Config{})
}

func (s *Service) ensureSchema(ctx context.Context) (InitResult, error) {
	has := s.db.WithContext(ctx).Migrator().HasTable(tableAdminMember)
	if !has {
		sqlText, err := createTableSQL(s.driver)
		if err != nil {
			return InitResult{}, err
		}
		if err = s.db.WithContext(ctx).Exec(sqlText).Error; err != nil {
			return InitResult{}, ErrInternal
		}
	}
	var count int64
	if err := s.db.WithContext(ctx).Model(&adminMember{}).Count(&count).Error; err != nil {
		return InitResult{}, ErrInternal
	}
	if count > 0 {
		return InitResult{}, nil
	}
	raw, err := RandomPassword()
	if err != nil {
		return InitResult{}, ErrInternal
	}
	hash, err := HashPassword(raw)
	if err != nil {
		return InitResult{}, ErrInternal
	}
	row := adminMember{
		Username:     "admin",
		PasswordHash: hash,
		Name:         "管理员",
		Status:       statusNormal,
	}
	if err = s.db.WithContext(ctx).Create(&row).Error; err != nil {
		return InitResult{}, ErrInternal
	}
	return InitResult{AdminCreated: true, Username: "admin", Password: raw}, nil
}

func (s *Service) Authenticate(ctx context.Context, username, password string) (*Member, int64, error) {
	if username == "" || password == "" {
		return nil, 0, ErrEmptyCredentials
	}
	if utf8.RuneCountInString(username) > maxUsernameLen {
		return nil, 0, ErrBadCredentials
	}
	row, err := s.activeBy(ctx, "username = ? AND deleted_at IS NULL", username)
	if errors.Is(err, ErrNotFound) {
		return nil, 0, ErrBadCredentials
	}
	if err != nil {
		return nil, 0, err
	}
	if !VerifyPassword(row.PasswordHash, password) {
		return nil, 0, ErrBadCredentials
	}
	if row.Status != statusNormal {
		return nil, 0, ErrDisabled
	}
	return toMember(row), row.TokenVersion, nil
}

func (s *Service) Active(ctx context.Context, id int64) (*Member, int64, error) {
	row, err := s.activeBy(ctx, "id = ? AND deleted_at IS NULL", id)
	if err != nil {
		return nil, 0, err
	}
	if row.Status != statusNormal {
		return nil, 0, ErrDisabled
	}
	return toMember(row), row.TokenVersion, nil
}

func (s *Service) activeBy(ctx context.Context, query string, arg any) (*adminMember, error) {
	var row adminMember
	err := s.db.WithContext(ctx).Where(query, arg).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, ErrInternal
	}
	return &row, nil
}

func toMember(row *adminMember) *Member {
	return &Member{ID: row.ID, Username: row.Username, Name: row.Name}
}

func (s *Service) RecordLogin(ctx context.Context, id int64, ip string) error {
	now := time.Now()
	err := s.db.WithContext(ctx).Model(&adminMember{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{"last_login_ip": ip, "last_login_at": now}).Error
	if err != nil {
		return ErrInternal
	}
	return nil
}

func (s *Service) BumpTokenVersion(ctx context.Context, id, version int64) error {
	res := s.db.WithContext(ctx).Model(&adminMember{}).
		Where("id = ? AND deleted_at IS NULL AND token_version = ?", id, version).
		Update("token_version", gorm.Expr("token_version + 1"))
	if res.Error != nil {
		return ErrInternal
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	var current adminMember
	err := s.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		Take(&current).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return ErrInternal
	}
	for attempt := 0; attempt < 4; attempt++ {
		now := time.Now()
		res := s.db.WithContext(ctx).Model(&adminMember{}).
			Where("id = ? AND deleted_at IS NULL", id).
			Updates(map[string]any{
				"username":   DeletedUsername(current.Username, now),
				"deleted_at": now,
			})
		if res.Error == nil {
			if res.RowsAffected == 0 {
				return ErrNotFound
			}
			return nil
		}
		if attempt < 3 && isUniqueViolation(res.Error) {
			continue
		}
		return ErrInternal
	}
	return ErrInternal
}

func isUniqueViolation(err error) bool {
	var pg *pgconn.PgError
	if errors.As(err, &pg) && pg.Code == "23505" {
		return true
	}
	var my *mysqldriver.MySQLError
	if errors.As(err, &my) && my.Number == 1062 {
		return true
	}
	return false
}
