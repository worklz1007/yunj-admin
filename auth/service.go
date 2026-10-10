package auth

import (
	"context"
	"errors"
	"time"

	"github.com/worklz1007/yunj-admin/member"
)

const defaultTokenTTL = 24 * time.Hour

type Config struct {
	JWTSecret string
	TokenTTL  time.Duration
}

type Identity struct {
	ID           int64
	Username     string
	Name         string
	TokenVersion int64
}

func New(members *member.Service, cfg Config) (*Service, error) {
	if members == nil || cfg.JWTSecret == "" || cfg.TokenTTL < 0 {
		return nil, ErrInvalidConfig
	}
	ttl := cfg.TokenTTL
	if ttl == 0 {
		ttl = defaultTokenTTL
	}
	return &Service{members: members, secret: cfg.JWTSecret, ttl: ttl}, nil
}

func (s *Service) Login(ctx context.Context, username, password, ip string) (string, int64, error) {
	user, version, err := s.members.Authenticate(ctx, username, password)
	if err != nil {
		return "", 0, err
	}
	token, exp, err := s.sign(user.ID, version, time.Now())
	if err != nil {
		return "", 0, err
	}
	if err = s.members.RecordLogin(ctx, user.ID, ip); err != nil {
		return "", 0, err
	}
	return token, exp, nil
}

func (s *Service) Current(ctx context.Context, accessToken string) (*Identity, error) {
	return s.identityFromToken(ctx, accessToken)
}

func (s *Service) Logout(ctx context.Context, id, tokenVersion int64) error {
	if id <= 0 {
		return ErrUnauthenticated
	}
	err := s.members.BumpTokenVersion(ctx, id, tokenVersion)
	if errors.Is(err, member.ErrNotFound) {
		return ErrTokenRevoked
	}
	return err
}

func (s *Service) identityFromToken(ctx context.Context, accessToken string) (*Identity, error) {
	id, version, err := s.parse(accessToken)
	if err != nil {
		return nil, err
	}
	user, storedVersion, err := s.members.Active(ctx, id)
	if errors.Is(err, member.ErrNotFound) {
		return nil, ErrTokenRevoked
	}
	if err != nil {
		return nil, err
	}
	if storedVersion != version {
		return nil, ErrTokenRevoked
	}
	return &Identity{
		ID:           user.ID,
		Username:     user.Username,
		Name:         user.Name,
		TokenVersion: version,
	}, nil
}
