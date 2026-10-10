package member

import (
	"context"
	"errors"
	"testing"
)

func TestInitRejectsInvalidConfig(t *testing.T) {
	ctx := context.Background()
	cases := []InitConfig{
		{},
		{Driver: "sqlite", DSN: "x"},
		{Driver: "postgres", DSN: ""},
		{Driver: "mysql", DSN: ""},
	}
	for _, cfg := range cases {
		_, _, _, err := Init(ctx, cfg)
		if !errors.Is(err, ErrInvalidConfig) {
			t.Fatalf("cfg=%+v err=%v", cfg, err)
		}
	}
}
