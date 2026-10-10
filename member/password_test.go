package member

import "testing"

func TestHashAndVerifyPassword(t *testing.T) {
	hash, err := HashPassword("admin-pass")
	if err != nil {
		t.Fatal(err)
	}
	if !VerifyPassword(hash, "admin-pass") {
		t.Fatal("expected password to match")
	}
	if VerifyPassword(hash, "other") {
		t.Fatal("expected password to miss")
	}
}

func TestRandomPassword(t *testing.T) {
	got, err := RandomPassword()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 16 {
		t.Fatalf("len=%d", len(got))
	}
	for _, r := range got {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			t.Fatalf("bad char %q", r)
		}
	}
}
