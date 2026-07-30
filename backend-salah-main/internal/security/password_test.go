package security

import "testing"

func TestPasswordHashing(t *testing.T) {
	password := "frase-segura-de-prueba"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == password || !IsPasswordHash(hash) {
		t.Fatal("password was not stored as a bcrypt hash")
	}
	if !CheckPassword(hash, password) {
		t.Fatal("correct password was rejected")
	}
	if CheckPassword(hash, "otra-contrasena") {
		t.Fatal("incorrect password was accepted")
	}
}

func TestShortPasswordIsRejected(t *testing.T) {
	if _, err := HashPassword("muy-corta"); err == nil {
		t.Fatal("short password was accepted")
	}
}
