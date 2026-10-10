package auth

import (
	"context"
	"testing"
	"time"
)

func TestSignAndParse(t *testing.T) {
	svc := &Service{secret: "test-secret", ttl: time.Hour}
	now := time.Now().Truncate(time.Second)
	token, exp, err := svc.sign(7, 3, now)
	if err != nil {
		t.Fatal(err)
	}
	if exp != now.Add(time.Hour).Unix() {
		t.Fatalf("exp=%d", exp)
	}
	id, version, err := svc.parse(token)
	if err != nil {
		t.Fatal(err)
	}
	if id != 7 || version != 3 {
		t.Fatalf("id=%d version=%d", id, version)
	}
	if _, _, err = svc.parse(token + "x"); err != ErrUnauthenticated {
		t.Fatalf("err=%v", err)
	}
}

func TestNewRejectsInvalidConfig(t *testing.T) {
	if _, err := New(nil, Config{JWTSecret: "s"}); err != ErrInvalidConfig {
		t.Fatalf("nil members err=%v", err)
	}
}

func TestLogoutRejectsNonPositiveID(t *testing.T) {
	svc := &Service{}
	if err := svc.Logout(context.Background(), 0, 1); err != ErrUnauthenticated {
		t.Fatalf("err=%v", err)
	}
}
