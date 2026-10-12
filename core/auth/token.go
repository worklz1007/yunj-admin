package auth

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/worklz1007/yunj-admin/core/member"
)

const issuer = "yunj-admin"

type Service struct {
	members *member.Service
	secret  string
	ttl     time.Duration
}

type jwtClaims struct {
	ID           int64 `json:"id"`
	TokenVersion int64 `json:"token_version"`
	jwt.RegisteredClaims
}

func (s *Service) sign(id, version int64, now time.Time) (string, int64, error) {
	exp := now.Add(s.ttl)
	claims := jwtClaims{
		ID:           id,
		TokenVersion: version,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   strconv.FormatInt(id, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			ID:        uuid.NewString(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", 0, ErrInternal
	}
	return signed, exp.Unix(), nil
}

func (s *Service) parse(tokenStr string) (int64, int64, error) {
	if tokenStr == "" {
		return 0, 0, ErrUnauthenticated
	}
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	token, err := parser.ParseWithClaims(tokenStr, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil || !token.Valid {
		return 0, 0, ErrUnauthenticated
	}
	claims, ok := token.Claims.(*jwtClaims)
	if !ok || claims.Issuer != issuer {
		return 0, 0, ErrUnauthenticated
	}
	return claims.ID, claims.TokenVersion, nil
}
