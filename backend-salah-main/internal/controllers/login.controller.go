package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/security"
	"encoding/json"
	"net/http"
	"strings"
)

type auth struct{}

var Auth auth

var dummyPasswordHash = func() string {
	hash, err := security.HashPassword("dummy-password-not-valid")
	if err != nil {
		panic(err)
	}
	return hash
}()

func (auth) AuthLoginWeb(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	var credentials struct {
		Usuario string `json:"usuario"`
		Contra  string `json:"contra"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&credentials); err != nil {
		http.Error(w, "Solicitud no valida", http.StatusBadRequest)
		return
	}
	credentials.Usuario = strings.TrimSpace(credentials.Usuario)
	if credentials.Usuario == "" || credentials.Contra == "" {
		http.Error(w, "Credenciales incorrectas", http.StatusUnauthorized)
		return
	}

	ip := requestClientIP(r)
	if !authAttempts.allowIP(ip) {
		http.Error(w, "Demasiados intentos; pruebe mas tarde", http.StatusTooManyRequests)
		return
	}

	var existing models.Usuario
	queryErr := db.GDB.Where("lower(usuario) = lower(?)", credentials.Usuario).First(&existing).Error
	if queryErr != nil {
		_ = security.CheckPassword(dummyPasswordHash, credentials.Contra)
		respondLoginFailure(w, ip, credentials.Usuario)
		return
	}
	if !existing.Estado || !security.CheckPassword(existing.Contra, credentials.Contra) {
		respondLoginFailure(w, ip, credentials.Usuario)
		return
	}
	if existing.SessionVersion == 0 {
		existing.SessionVersion = 1
		if err := db.GDB.Model(&existing).UpdateColumn("session_version", existing.SessionVersion).Error; err != nil {
			http.Error(w, "No se pudo iniciar sesion", http.StatusInternalServerError)
			return
		}
	}

	var role models.Rol
	if err := db.GDB.Where("id = ?", existing.IDRol).First(&role).Error; err != nil {
		http.Error(w, "No se pudo iniciar sesion", http.StatusInternalServerError)
		return
	}
	normalizedRole := strings.ToLower(strings.TrimSpace(role.Nombre))
	if normalizedRole != "admin" && normalizedRole != "encargado de ventas" && normalizedRole != "vendedor" && normalizedRole != "contador" {
		authAttempts.failure(ip, credentials.Usuario)
		http.Error(w, "Credenciales incorrectas", http.StatusUnauthorized)
		return
	}
	token, csrfToken, err := security.IssueToken(existing.ID, existing.SessionVersion)
	if err != nil {
		http.Error(w, "No se pudo iniciar sesion", http.StatusInternalServerError)
		return
	}
	authAttempts.success(ip, credentials.Usuario)
	security.SetSessionCookies(w, token, csrfToken)

	response := struct {
		ID                 uint   `json:"id"`
		Nombre             string `json:"nombre"`
		Rol                string `json:"rol"`
		MustChangePassword bool   `json:"must_change_password"`
	}{
		ID:                 existing.ID,
		Nombre:             existing.Nombre,
		Rol:                normalizedRole,
		MustChangePassword: existing.MustChangePassword,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func respondLoginFailure(w http.ResponseWriter, ip, username string) {
	if authAttempts.failure(ip, username) {
		http.Error(w, "Demasiados intentos; pruebe mas tarde", http.StatusTooManyRequests)
		return
	}
	http.Error(w, "Credenciales incorrectas", http.StatusUnauthorized)
}
