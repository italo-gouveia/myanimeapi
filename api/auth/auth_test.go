package auth

import (
	"testing"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "E2Epassword123!"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}
	t.Logf("Hash: %s (len=%d)", hash, len(hash))

	valid, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash error: %v", err)
	}
	if !valid {
		t.Fatalf("CheckPasswordHash returned false for correct password")
	}
}

func TestMigrateHashArgon2(t *testing.T) {
	password := "E2Epassword123!"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}

	returned, err := MigrateHash(password, hash)
	if err != nil {
		t.Fatalf("MigrateHash error for argon2 hash: %v", err)
	}
	if returned != hash {
		t.Fatalf("MigrateHash returned different hash: got %q want %q", returned, hash)
	}
}

func TestMigrateHashWrongPassword(t *testing.T) {
	hash, err := HashPassword("correct_password")
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}

	_, err = MigrateHash("wrong_password", hash)
	if err == nil {
		t.Fatal("MigrateHash should have returned error for wrong password")
	}
}
