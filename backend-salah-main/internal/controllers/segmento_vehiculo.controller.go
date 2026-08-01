package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/functions"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/querys"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type SegmentoVehiculoDAO struct {
	ID          uint   `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Estado      string `json:"estado"`
	IDCategoria uint   `json:"id_categoria"`
	Categoria   string `json:"categoria"`
}

type SegmentoVehiculoMOD struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Estado      string `json:"estado"`
	IDCategoria uint   `json:"id_categoria"`
}

func ObtenerSegmentosVehiculos(w http.ResponseWriter, r *http.Request) {
	segmentos := make([]SegmentoVehiculoDAO, 0)

	query := querys.SegmentosVehiculos
	args := []interface{}{}
	if idCategoria := r.URL.Query().Get("id_categoria"); idCategoria != "" {
		if _, err := strconv.Atoi(idCategoria); err != nil {
			http.Error(w, "Categoria no valida", http.StatusBadRequest)
			return
		}
		query += " where s.id_categoria = ?"
		args = append(args, idCategoria)
	}
	query += " order by c.nombre asc, s.nombre asc;"

	if err := db.GDB.Raw(query, args...).Scan(&segmentos).Error; err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(segmentos)
}

func ObtenerSegmentoVehiculo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var segmento SegmentoVehiculoDAO

	if err := db.GDB.Raw(querys.SegmentoVehiculo, id).Scan(&segmento).Error; err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(segmento)
}

func AgregarSegmentoVehiculo(w http.ResponseWriter, r *http.Request) {
	var segmento SegmentoVehiculoMOD

	if err := json.NewDecoder(r.Body).Decode(&segmento); err != nil {
		http.Error(w, "Solicitud no valida", http.StatusBadRequest)
		return
	}

	nuevoEstado, err := functions.ActualizarEstado(segmento.Estado)
	if err != nil {
		http.Error(w, "Estado no valido", http.StatusBadRequest)
		return
	}

	if err := validarCategoriaVehiculo(segmento.IDCategoria); err != nil {
		http.Error(w, "Categoria no encontrada", http.StatusBadRequest)
		return
	}

	nuevoSegmento := models.SegmentoVehiculo{
		Nombre:      segmento.Nombre,
		Descripcion: segmento.Descripcion,
		Estado:      nuevoEstado,
		IDCategoria: segmento.IDCategoria,
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&nuevoSegmento).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Segmento", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nuevoSegmento)
}

func ModificarSegmentoVehiculo(w http.ResponseWriter, r *http.Request) {
	idSegmento := mux.Vars(r)["id"]
	var segmentoExistente models.SegmentoVehiculo

	if err := db.GDB.Where("id = ?", idSegmento).First(&segmentoExistente).Error; err != nil {
		http.Error(w, "Segmento no encontrado", http.StatusNotFound)
		return
	}

	var segmentoActualizado SegmentoVehiculoMOD
	if err := json.NewDecoder(r.Body).Decode(&segmentoActualizado); err != nil {
		http.Error(w, "Solicitud no valida", http.StatusBadRequest)
		return
	}

	nuevoEstado, err := functions.ActualizarEstado(segmentoActualizado.Estado)
	if err != nil {
		http.Error(w, "Estado no valido", http.StatusBadRequest)
		return
	}

	if err := validarCategoriaVehiculo(segmentoActualizado.IDCategoria); err != nil {
		http.Error(w, "Categoria no encontrada", http.StatusBadRequest)
		return
	}

	segmentoExistente.Nombre = segmentoActualizado.Nombre
	segmentoExistente.Descripcion = segmentoActualizado.Descripcion
	segmentoExistente.Estado = nuevoEstado
	segmentoExistente.IDCategoria = segmentoActualizado.IDCategoria

	if err := db.GDB.Save(&segmentoExistente).Error; err != nil {
		http.Error(w, "Error al actualizar Segmento", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&segmentoExistente)
}

func validarCategoriaVehiculo(idCategoria uint) error {
	var categoria models.CategoriaVehiculo
	return db.GDB.Where("id = ?", idCategoria).First(&categoria).Error
}
