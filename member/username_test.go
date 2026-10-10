package member

import (
	"testing"
	"time"
)

func TestDeletedUsername(t *testing.T) {
	at := time.UnixMilli(1760080000123)
	got := DeletedUsername("admin", at)
	if got != "admin_del_1760080000123" {
		t.Fatalf("got %s", got)
	}
}
