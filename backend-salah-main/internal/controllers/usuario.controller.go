package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/functions"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/querys"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UsuarioDAO struct {
	ID        uint   `json:"id"`
	Nombre    string `json:"nombre"`
	Apellido  string `json:"apellido"`
	CI        string `json:"ci"`
	Celular   string `json:"celular"`
	Direccion string `json:"direccion"`
	Foto      string `json:"foto"`
	Usuario   string `json:"usuario"`
	Contra    string `json:"contra"`
	Estado    string `json:"estado"`
	Rol       string `json:"rol"`
}

type UsuarioModificado struct {
	Nombre    string `json:"nombre"`
	Apellido  string `json:"apellido"`
	CI        string `json:"ci"`
	Celular   string `json:"celular"`
	Direccion string `json:"direccion"`
	Foto      string `json:"foto"`
	Usuario   string `json:"usuario"`
	Contra    string `json:"contra"`
	Estado    string `json:"estado"`
	Rol       string `json:"rol"`
}

// validatePassword checks if the password meets the criteria
func validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	if !regexp.MustCompile(`[a-zA-Z]`).MatchString(password) {
		return false
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return false
	}
	return true
}

func ObtenerUsuarios(w http.ResponseWriter, r *http.Request) {

	var usuarios []UsuarioDAO

	err := db.GDB.Raw(querys.Usuarios).Scan(&usuarios).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	for i, usuario := range usuarios {
		if usuario.Foto != "" && usuario.Foto != "N/A" {
			encodeImagen, err := encodeImageToBase64(usuario.Foto)
			if err == nil {
				usuarios[i].Foto = encodeImagen
			} else {
				usuarios[i].Foto = "N/A"
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarios)
}

func ObtenerUsuario(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var usuario UsuarioDAO

	err := db.GDB.Raw(querys.Usuario, id).Scan(&usuario).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	if usuario.Foto != "" && usuario.Foto != "N/A" {
		encodeImagen, err := encodeImageToBase64(usuario.Foto)
		if err == nil {
			usuario.Foto = encodeImagen
		} else {
			usuario.Foto = "N/A"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

func AgregarUsuario(w http.ResponseWriter, r *http.Request) {
	var usuario UsuarioDAO

	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Error al parsear el formulario", http.StatusInternalServerError)
			return
		}

		usuario.Nombre = r.FormValue("nombre")
		usuario.Apellido = r.FormValue("apellido")
		usuario.CI = r.FormValue("ci")
		usuario.Celular = r.FormValue("celular")
		usuario.Direccion = r.FormValue("direccion")
		usuario.Usuario = r.FormValue("usuario")
		usuario.Contra = r.FormValue("contra")
		usuario.Estado = r.FormValue("estado")
		usuario.Rol = r.FormValue("rol")

		direccionImagen := "N/A"
		file, handler, err := r.FormFile("foto")
		if err == nil {
			defer file.Close()

			if err := os.MkdirAll("internal/images/usuarios", os.ModePerm); err != nil {
				http.Error(w, "Error al preparar la carpeta de fotos", http.StatusInternalServerError)
				return
			}

			nombreImagen := fmt.Sprintf("usuario-%s%s", uuid.New().String(), filepath.Ext(handler.Filename))
			rutaImagen := "internal/images/usuarios/" + nombreImagen

			outFile, err := os.Create(rutaImagen)
			if err != nil {
				http.Error(w, "Error al guardar la foto", http.StatusInternalServerError)
				return
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, file); err != nil {
				http.Error(w, "Error al escribir la foto", http.StatusInternalServerError)
				return
			}

			direccionImagen = rutaImagen
		} else if err != http.ErrMissingFile {
			http.Error(w, "Error al obtener la foto: "+err.Error(), http.StatusInternalServerError)
			return
		}
		usuario.Foto = direccionImagen
	} else {
		if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Validate password
	if !validatePassword(usuario.Contra) {
		http.Error(w, "La contraseña debe tener al menos 8 caracteres, incluyendo letras y números.", http.StatusBadRequest)
		return
	}

	nuevoEstado, err := functions.ActualizarEstado(usuario.Estado)
	if err != nil {
		http.Error(w, "Estado no valido", http.StatusBadRequest)
		return
	}
	nuevoRol, err := functions.ActualizarRol(usuario.Rol)
	if err != nil {
		http.Error(w, "Rol no valido:"+err.Error(), http.StatusBadRequest)
		return
	}

	nuevoUsuario := models.Usuario{
		Nombre:    usuario.Nombre,
		Apellido:  usuario.Apellido,
		CI:        usuario.CI,
		Celular:   usuario.Celular,
		Direccion: usuario.Direccion,
		Foto:      usuario.Foto,
		Usuario:   usuario.Usuario,
		Contra:    usuario.Contra,
		Estado:    nuevoEstado,
		IDRol:     nuevoRol,
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&nuevoUsuario).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Usuario", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&nuevoUsuario)
}

func ModificarUsuario(w http.ResponseWriter, r *http.Request) {
	id_usuario := mux.Vars(r)["id"]
	var usuarioExistente models.Usuario

	err := db.GDB.Where("id = ?", id_usuario).First(&usuarioExistente).Error
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	var usuarioActualizado UsuarioModificado
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Error al parsear el formulario: "+err.Error(), http.StatusBadRequest)
			return
		}

		usuarioActualizado.Nombre = r.FormValue("nombre")
		usuarioActualizado.Apellido = r.FormValue("apellido")
		usuarioActualizado.CI = r.FormValue("ci")
		usuarioActualizado.Celular = r.FormValue("celular")
		usuarioActualizado.Direccion = r.FormValue("direccion")
		usuarioActualizado.Usuario = r.FormValue("usuario")
		usuarioActualizado.Contra = r.FormValue("contra")
		usuarioActualizado.Estado = r.FormValue("estado")
		usuarioActualizado.Rol = r.FormValue("rol")

		file, handler, err := r.FormFile("foto")
		if err == nil {
			defer file.Close()

			if err := os.MkdirAll("internal/images/usuarios", os.ModePerm); err != nil {
				http.Error(w, "Error al preparar la carpeta de fotos", http.StatusInternalServerError)
				return
			}

			if usuarioExistente.Foto != "" && usuarioExistente.Foto != "N/A" {
				_ = os.Remove(usuarioExistente.Foto)
			}

			nombreImagen := fmt.Sprintf("usuario-%s%s", uuid.New().String(), filepath.Ext(handler.Filename))
			rutaImagen := "internal/images/usuarios/" + nombreImagen

			outFile, err := os.Create(rutaImagen)
			if err != nil {
				http.Error(w, "Error al guardar la nueva foto", http.StatusInternalServerError)
				return
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, file); err != nil {
				http.Error(w, "Error al escribir la nueva foto", http.StatusInternalServerError)
				return
			}

			usuarioExistente.Foto = rutaImagen
		} else if err != http.ErrMissingFile {
			http.Error(w, "Error al procesar la foto: "+err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&usuarioActualizado); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Validate password if provided
	if usuarioActualizado.Contra != "" {
		if !validatePassword(usuarioActualizado.Contra) {
			http.Error(w, "La contraseña debe tener al menos 8 caracteres, incluyendo letras y números.", http.StatusBadRequest)
			return
		}
		usuarioExistente.Contra = usuarioActualizado.Contra
	}

	// Cambios
	usuarioExistente.Nombre = usuarioActualizado.Nombre
	usuarioExistente.Apellido = usuarioActualizado.Apellido
	usuarioExistente.CI = usuarioActualizado.CI
	usuarioExistente.Celular = usuarioActualizado.Celular
	if usuarioActualizado.Direccion != "" {
		usuarioExistente.Direccion = usuarioActualizado.Direccion
	}
	usuarioExistente.Usuario = usuarioActualizado.Usuario
	nuevoEstado, err := functions.ActualizarEstado(usuarioActualizado.Estado)
	if err != nil {
		http.Error(w, "Estado no valido", http.StatusBadRequest)
		return
	}
	usuarioExistente.Estado = nuevoEstado
	nuevoRol, err := functions.ActualizarRol(usuarioActualizado.Rol)
	if err != nil {
		http.Error(w, "Rol no valido", http.StatusBadRequest)
		return
	}
	usuarioExistente.IDRol = nuevoRol

	if err := db.GDB.Save(&usuarioExistente).Error; err != nil {
		http.Error(w, "Error al actualizar usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&usuarioExistente)
}
