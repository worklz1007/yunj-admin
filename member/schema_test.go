package member

import (
	"strings"
	"testing"
)

func TestCreateTableSQL(t *testing.T) {
	pg, err := createTableSQL("postgres")
	if err != nil {
		t.Fatal(err)
	}
	for _, part := range []string{"bigserial", "timestamptz", "uk_admin_member_username", "deleted_at"} {
		if !strings.Contains(pg, part) {
			t.Fatalf("postgres sql missing %s", part)
		}
	}
	my, err := createTableSQL("mysql")
	if err != nil {
		t.Fatal(err)
	}
	for _, part := range []string{"AUTO_INCREMENT", "utf8mb4", "ON UPDATE CURRENT_TIMESTAMP", "uk_admin_member_username"} {
		if !strings.Contains(my, part) {
			t.Fatalf("mysql sql missing %s", part)
		}
	}
	if _, err = createTableSQL("sqlite"); err != ErrInvalidConfig {
		t.Fatalf("err=%v", err)
	}
}
