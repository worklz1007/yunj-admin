package boot

import (
	"context"
	"time"

	"github.com/worklz1007/yunj-admin/core/auth"
	"github.com/worklz1007/yunj-admin/core/member"
	"github.com/worklz1007/yunj-admin/data"
	"github.com/worklz1007/yunj-admin/pkg/util"
)

type Config struct {
	Driver    string
	DSN       string
	JWTSecret string
	TokenTTL  time.Duration
}

type Result struct {
	AdminCreated bool
	Username     string
	Password     string
	Auth         *auth.Service
	Member       *member.Service
}

func Init(ctx context.Context, cfg Config) (*Result, func(), error) {
	if (cfg.Driver != "mysql" && cfg.Driver != "postgres") || cfg.DSN == "" || cfg.JWTSecret == "" || cfg.TokenTTL < 0 {
		return nil, nil, member.ErrInvalidConfig
	}
	db, cleanup, err := util.Open(ctx, cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, nil, err
	}
	if err = data.EnsureAdminMember(ctx, db, cfg.Driver); err != nil {
		cleanup()
		return nil, nil, err
	}
	members := member.New(db)
	created, username, password, err := members.Seed(ctx)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	authSvc, err := auth.New(members, auth.Config{JWTSecret: cfg.JWTSecret, TokenTTL: cfg.TokenTTL})
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return &Result{
		AdminCreated: created,
		Username:     username,
		Password:     password,
		Auth:         authSvc,
		Member:       members,
	}, cleanup, nil
}
