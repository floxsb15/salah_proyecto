package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/functions"
	"backend-restaurant-delitto/internal/models"
	"encoding/json"
	"net/http"
	"strings"
)

type MarcaVehiculoDAO struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Estado string `json:"estado"`
}

type MarcaVehiculoMOD struct {
	Nombre string `json:"nombre"`
	Estado string `json:"estado"`
}

type AnioVehiculoDAO struct {
	ID     uint   `json:"id"`
	Valor  uint   `json:"valor"`
	Estado string `json:"estado"`
}

type AnioVehiculoMOD struct {
	Valor  uint   `json:"valor"`
	Estado string `json:"estado"`
}

func ObtenerMarcasVehiculos(w http.ResponseWriter, r *http.Request) {
	marcas := make([]MarcaVehiculoDAO, 0)
	query := `
		select id, nombre,
			case when estado then 'Activo' else 'Inactivo' end as estado
		from marcas_vehiculos
		order by nombre asc`

	if err := db.GDB.Raw(query).Scan(&marcas).Error; err != nil {
		http.Error(w, "Error al consultar marcas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(marcas)
}

func AgregarMarcaVehiculo(w http.ResponseWriter, r *http.Request) {
	var marca MarcaVehiculoMOD
	if err := json.NewDecoder(r.Body).Decode(&marca); err != nil {
		http.Error(w, "Solicitud no valida", http.StatusBadRequest)
		return
	}

	nombre := strings.TrimSpace(marca.Nombre)
	if nombre == "" {
		http.Error(w, "Marca requerida", http.StatusBadRequest)
		return
	}

	estado, err := functions.ActualizarEstado(marca.Estado)
	if err != nil {
		http.Error(w, "Estado no valido", http.StatusBadRequest)
		return
	}

	nuevaMarca := models.MarcaVehiculo{
		Nombre: nombre,
		Estado: estado,
	}

	if err := db.GDB.Create(&nuevaMarca).Error; err != nil {
		http.Error(w, "Error al guardar marca", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nuevaMarca)
}

func ObtenerAniosVehiculos(w http.ResponseWriter, r *http.Request) {
	anios := make([]AnioVehiculoDAO, 0)
	query := `
		select id, valor,
			case when estado then 'Activo' else 'Inactivo' end as estado
		from anios_vehiculos
		order by valor desc`

	if err := db.GDB.Raw(query).Scan(&anios).Error; err != nil {
		http.Error(w, "Error al consultar anios", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(anios)
}

func AgregarAnioVehiculo(w http.ResponseWriter, r *http.Request) {
	var anio AnioVehiculoMOD
	if err := json.NewDecoder(r.Body).Decode(&anio); err != nil {
		http.Error(w, "Solicitud no valida", http.StatusBadRequest)
		return
	}

	if anio.Valor < 1900 || anio.Valor > 2100 {
		http.Error(w, "Anio no valido", http.StatusBadRequest)
		return
	}

	estado, err := functions.ActualizarEstado(anio.Estado)
	if err != nil {
		http.Error(w, "Estado no valido", http.StatusBadRequest)
		return
	}

	nuevoAnio := models.AnioVehiculo{
		Valor:  anio.Valor,
		Estado: estado,
	}

	if err := db.GDB.Create(&nuevoAnio).Error; err != nil {
		http.Error(w, "Error al guardar anio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nuevoAnio)
}
