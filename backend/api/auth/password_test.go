package auth

import "testing"

func TestHashPassword_VerifyRoundTrip(t *testing.T) {
	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if !VerifyPassword("correct horse battery staple", hash) {
		t.Error("expected correct password to verify")
	}

	if VerifyPassword("wrong password", hash) {
		t.Error("expected wrong password to fail verification")
	}
}

func TestHashPassword_UniqueSaltPerCall(t *testing.T) {
	// Regression test: HashPassword previously used a single fixed salt
	// (the JWT secret) for every user, so identical passwords produced
	// identical hashes. Each call must now use its own random salt.
	hash1, err := HashPassword("same-password")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	hash2, err := HashPassword("same-password")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hash1 == hash2 {
		t.Error("expected two hashes of the same password to differ (unique salt per call)")
	}

	if !VerifyPassword("same-password", hash1) || !VerifyPassword("same-password", hash2) {
		t.Error("both independently-salted hashes should still verify the same password")
	}
}

func TestVerifyPassword_MalformedHash(t *testing.T) {
	cases := []string{"", "no-dot-separator", "a.b.c", "not-base64!.also-not-base64!"}
	for _, c := range cases {
		if VerifyPassword("anything", c) {
			t.Errorf("expected malformed hash %q to fail verification, not panic or succeed", c)
		}
	}
}
