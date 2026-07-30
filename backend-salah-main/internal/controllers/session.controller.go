package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/security"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

func ObtenerPerfilActual(w http.ResponseWriter, r *http.Request) {
	principal, ok := security.PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "Autenticacion requerida", http.StatusUnauthorized)
		return
	}
	var profile UsuarioDAO
	if err := db.GDB.Raw(`
		SELECT u.id, u.nombre, u.apellido, u.ci, u.celular, u.direccion, u.foto,
		       u.usuario,
		       CASE WHEN u.estado THEN 'Activo' ELSE 'Inactivo' END AS estado,
		       r.nombre AS rol
		FROM usuarios u
		JOIN roles r ON r.id = u.id_rol
		WHERE u.id = ?
		LIMIT 1`, principal.ID).Scan(&profile).Error; err != nil {
		http.Error(w, "No se pudo obtener el perfil", http.StatusInternalServerError)
		return
	}
	if profile.Foto != "" && profile.Foto != "N/A" {
		if encoded, err := encodeImageToBase64(profile.Foto); err == nil {
			profile.Foto = encoded
		} else {
			profile.Foto = "N/A"
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(profile)
}

func CerrarSesion(w http.ResponseWriter, r *http.Request) {
	principal, ok := security.PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "Autenticacion requerida", http.StatusUnauthorized)
		return
	}
	if err := db.GDB.Model(&models.Usuario{}).Where("id = ?", principal.ID).
		UpdateColumn("session_version", gorm.Expr("session_version + 1")).Error; err != nil {
		http.Error(w, "No se pudo cerrar sesion", http.StatusInternalServerError)
		return
	}
	security.ClearSessionCookies(w)
	w.WriteHeader(http.StatusNoContent)
}

func CambiarContrasenaActual(w http.ResponseWriter, r *http.Request) {
	principal, ok := security.PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "Autenticacion requerida", http.StatusUnauthorized)
		return
	}
	var payload struct {
		Actual string `json:"actual"`
		Nueva  string `json:"nueva"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "Solicitud no valida", http.StatusBadRequest)
		return
	}
	var user models.Usuario
	if err := db.GDB.First(&user, principal.ID).Error; err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}
	if !security.CheckPassword(user.Contra, payload.Actual) {
		http.Error(w, "Contrasena actual incorrecta", http.StatusUnauthorized)
		return
	}
	hash, err := security.HashPassword(payload.Nueva)
	if err != nil {
		http.Error(w, "La contrasena debe contener al menos 12 caracteres", http.StatusBadRequest)
		return
	}
	if err := db.GDB.Model(&user).Updates(map[string]interface{}{
		"contra":               hash,
		"must_change_password": false,
		"session_version":      gorm.Expr("session_version + 1"),
	}).Error; err != nil {
		http.Error(w, "No se pudo cambiar la contrasena", http.StatusInternalServerError)
		return
	}
	security.ClearSessionCookies(w)
	w.WriteHeader(http.StatusNoContent)
}
