package security

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

func HashPassword(password string) (string, error) {
	if len(password) < 12 {
		return "", errors.New("la contrasena debe contener al menos 12 caracteres")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func IsPasswordHash(value string) bool {
	return strings.HasPrefix(value, "$2a$") || strings.HasPrefix(value, "$2b$") || strings.HasPrefix(value, "$2y$")
}

// MigratePlaintextPasswords converts legacy plaintext passwords after the
// schema has been migrated. It is intentionally idempotent.
func MigratePlaintextPasswords() error {
	var users []models.Usuario
	if err := db.GDB.Select("id", "contra").Find(&users).Error; err != nil {
		return err
	}
	for _, user := range users {
		if IsPasswordHash(user.Contra) {
			continue
		}
		hashBytes, err := bcrypt.GenerateFromPassword([]byte(user.Contra), bcryptCost)
		if err != nil {
			return err
		}
		hash := string(hashBytes)
		if err := db.GDB.Model(&models.Usuario{}).
			Where("id = ?", user.ID).
			Updates(map[string]interface{}{"contra": hash, "session_version": 1, "must_change_password": true}).Error; err != nil {
			return err
		}
	}
	return nil
}
