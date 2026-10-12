package member

import (
	"context"
	"errors"
	"time"
	"unicode/utf8"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/worklz1007/yunj-admin/data"
)

const maxUsernameLen = 64

type Member struct {
	ID       int64
	Username string
	Name     string
}

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Seed(ctx context.Context) (bool, string, string, error) {
	count, err := data.CountAdminMember(ctx, s.db)
	if err != nil {
		return false, "", "", ErrInternal
	}
	if count > 0 {
		return false, "", "", nil
	}
	raw, err := RandomPassword()
	if err != nil {
		return false, "", "", ErrInternal
	}
	hash, err := HashPassword(raw)
	if err != nil {
		return false, "", "", ErrInternal
	}
	row := &data.AdminMember{
		Username:     "admin",
		PasswordHash: hash,
		Name:         "管理员",
		Status:       data.StatusNormal,
	}
	if err = data.InsertAdminMember(ctx, s.db, row); err != nil {
		return false, "", "", ErrInternal
	}
	return true, "admin", raw, nil
}

func (s *Service) Authenticate(ctx context.Context, username, password string) (*Member, int64, error) {
	if username == "" || password == "" {
		return nil, 0, ErrEmptyCredentials
	}
	if utf8.RuneCountInString(username) > maxUsernameLen {
		return nil, 0, ErrBadCredentials
	}
	row, err := data.TakeActiveAdminMember(ctx, s.db, "username = ? AND deleted_at IS NULL", username)
	if errors.Is(err, data.ErrNotFound) {
		return nil, 0, ErrBadCredentials
	}
	if err != nil {
		return nil, 0, ErrInternal
	}
	if !VerifyPassword(row.PasswordHash, password) {
		return nil, 0, ErrBadCredentials
	}
	if row.Status != data.StatusNormal {
		return nil, 0, ErrDisabled
	}
	return toMember(row), row.TokenVersion, nil
}

func (s *Service) Active(ctx context.Context, id int64) (*Member, int64, error) {
	row, err := data.TakeActiveAdminMember(ctx, s.db, "id = ? AND deleted_at IS NULL", id)
	if errors.Is(err, data.ErrNotFound) {
		return nil, 0, ErrNotFound
	}
	if err != nil {
		return nil, 0, ErrInternal
	}
	if row.Status != data.StatusNormal {
		return nil, 0, ErrDisabled
	}
	return toMember(row), row.TokenVersion, nil
}

func toMember(row *data.AdminMember) *Member {
	return &Member{ID: row.ID, Username: row.Username, Name: row.Name}
}

func (s *Service) RecordLogin(ctx context.Context, id int64, ip string) error {
	if err := data.UpdateAdminMemberLogin(ctx, s.db, id, ip, time.Now()); err != nil {
		return ErrInternal
	}
	return nil
}

func (s *Service) BumpTokenVersion(ctx context.Context, id, version int64) error {
	rows, err := data.BumpAdminMemberTokenVersion(ctx, s.db, id, version)
	if err != nil {
		return ErrInternal
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	current, err := data.TakeActiveAdminMember(ctx, s.db, "id = ? AND deleted_at IS NULL", id)
	if errors.Is(err, data.ErrNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return ErrInternal
	}
	for attempt := 0; attempt < 4; attempt++ {
		now := time.Now()
		rows, err := data.SoftDeleteAdminMember(ctx, s.db, id, DeletedUsername(current.Username, now), now)
		if err == nil {
			if rows == 0 {
				return ErrNotFound
			}
			return nil
		}
		if attempt < 3 && isUniqueViolation(err) {
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
