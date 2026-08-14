package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/functions"
	"backend-restaurant-delitto/internal/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

type ProveedorAutoDetalleDAO struct {
	ID               uint    `json:"id"`
	Nombre           string  `json:"nombre"`
	CINIT            string  `json:"ci_nit"`
	Telefono         string  `json:"telefono"`
	Email            string  `json:"email"`
	Direccion        string  `json:"direccion"`
	Tipo             string  `json:"tipo"`
	Observaciones    string  `json:"observaciones"`
	Estado           string  `json:"estado"`
	CantidadCompras  uint    `json:"cantidad_compras"`
	TotalCompradoUSD float64 `json:"total_comprado_usd"`
	TotalCompradoBOB float64 `json:"total_comprado_bob"`
}

type ProveedorCompraHistorialDAO struct {
	ID                   uint    `json:"id"`
	IDVehiculo           uint    `json:"id_vehiculo"`
	FechaCompra          string  `json:"fecha_compra"`
	Vehiculo             string  `json:"vehiculo"`
	PrecioCompraUSD      float64 `json:"precio_compra_usd"`
	PrecioCompraBOB      float64 `json:"precio_compra_bob"`
	GastosAdicionales    float64 `json:"gastos_adicionales"`
	GastosAdicionalesBOB float64 `json:"gastos_adicionales_bob"`
	CostoTotalUSD        float64 `json:"costo_total_usd"`
	CostoTotalBOB        float64 `json:"costo_total_bob"`
	EstadoPago           string  `json:"estado_pago"`
	TipoCambioUsado      float64 `json:"tipo_cambio_usado"`
}

func ObtenerProveedoresAutos(w http.ResponseWriter, r *http.Request) {
	proveedores := make([]ProveedorAutoDetalleDAO, 0)
	query := `
		select
			pa.id,
			pa.nombre,
			coalesce(pa.ci_nit, '') as ci_nit,
			coalesce(pa.telefono, '') as telefono,
			coalesce(pa.email, '') as email,
			coalesce(pa.direccion, '') as direccion,
			coalesce(pa.tipo, '') as tipo,
			coalesce(pa.observaciones, '') as observaciones,
			case when pa.estado then 'Activo' else 'Inactivo' end as estado,
			coalesce(count(ca.id), 0) as cantidad_compras,
			coalesce(sum(ca.costo_total_usd), 0) as total_comprado_usd,
			coalesce(sum(ca.costo_total_usd * nullif(ca.tipo_cambio_usado, 0)), 0) as total_comprado_bob
		from proveedores_autos pa
		left join compras_autos ca on ca.id_proveedor = pa.id
		group by pa.id
		order by pa.id asc`

	if err := db.GDB.Raw(query).Scan(&proveedores).Error; err != nil {
		http.Error(w, "Error al consultar proveedores", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proveedores)
}

func ObtenerProveedorAuto(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var proveedor ProveedorAutoDetalleDAO
	query := `
		select
			pa.id,
			pa.nombre,
			coalesce(pa.ci_nit, '') as ci_nit,
			coalesce(pa.telefono, '') as telefono,
			coalesce(pa.email, '') as email,
			coalesce(pa.direccion, '') as direccion,
			coalesce(pa.tipo, '') as tipo,
			coalesce(pa.observaciones, '') as observaciones,
			case when pa.estado then 'Activo' else 'Inactivo' end as estado,
			coalesce(count(ca.id), 0) as cantidad_compras,
			coalesce(sum(ca.costo_total_usd), 0) as total_comprado_usd,
			coalesce(sum(ca.costo_total_usd * nullif(ca.tipo_cambio_usado, 0)), 0) as total_comprado_bob
		from proveedores_autos pa
		left join compras_autos ca on ca.id_proveedor = pa.id
		where pa.id = ?
		group by pa.id
		limit 1`

	if err := db.GDB.Raw(query, id).Scan(&proveedor).Error; err != nil || proveedor.ID == 0 {
		http.Error(w, "Proveedor no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proveedor)
}

func AgregarProveedorAuto(w http.ResponseWriter, r *http.Request) {
	var payload ProveedorAutoDetalleDAO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Solicitud no valida", http.StatusBadRequest)
		return
	}
	proveedor, err := proveedorAutoFromPayload(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if duplicado, err := ciNitProveedorDuplicado(proveedor.CINIT, 0); err != nil {
		http.Error(w, "Error al validar CI/NIT", http.StatusInternalServerError)
		return
	} else if duplicado {
		http.Error(w, "Este CI/NIT ya existe", http.StatusConflict)
		return
	}

	if err := db.GDB.Create(&proveedor).Error; err != nil {
		http.Error(w, "Error al registrar proveedor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&proveedor)
}

func ModificarProveedorAuto(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var existente models.ProveedorAuto
	if err := db.GDB.Where("id = ?", id).First(&existente).Error; err != nil {
		http.Error(w, "Proveedor no encontrado", http.StatusNotFound)
		return
	}
	var payload ProveedorAutoDetalleDAO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Solicitud no valida", http.StatusBadRequest)
		return
	}
	actualizado, err := proveedorAutoFromPayload(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if duplicado, err := ciNitProveedorDuplicado(actualizado.CINIT, existente.ID); err != nil {
		http.Error(w, "Error al validar CI/NIT", http.StatusInternalServerError)
		return
	} else if duplicado {
		http.Error(w, "Este CI/NIT ya existe", http.StatusConflict)
		return
	}

	existente.Nombre = actualizado.Nombre
	existente.CINIT = actualizado.CINIT
	existente.Telefono = actualizado.Telefono
	existente.Email = actualizado.Email
	existente.Direccion = actualizado.Direccion
	existente.Tipo = actualizado.Tipo
	existente.Observaciones = actualizado.Observaciones
	existente.Estado = actualizado.Estado

	if err := db.GDB.Save(&existente).Error; err != nil {
		http.Error(w, "Error al actualizar proveedor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&existente)
}

func ObtenerHistorialComprasProveedor(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	historial := make([]ProveedorCompraHistorialDAO, 0)
	query := `
		select
			ca.id,
			ca.id_vehiculo,
			to_char(ca.fecha_compra, 'YYYY-MM-DD') as fecha_compra,
			coalesce(nullif(concat_ws(' ', nullif(trim(v.marca), ''), nullif(trim(v.modelo), ''), nullif(v.anio::text, '0')), ''), v.nombre) as vehiculo,
			ca.precio_compra_usd,
			ca.precio_compra_bob,
			ca.gastos_adicionales,
			round((ca.gastos_adicionales * ca.tipo_cambio_usado)::numeric, 2) as gastos_adicionales_bob,
			ca.costo_total_usd,
			round((ca.costo_total_usd * ca.tipo_cambio_usado)::numeric, 2) as costo_total_bob,
			ca.estado_pago,
			ca.tipo_cambio_usado
		from compras_autos ca
		inner join vehiculos v on v.id = ca.id_vehiculo
		where ca.id_proveedor = ?
		order by ca.fecha_compra desc, ca.id desc`

	if err := db.GDB.Raw(query, id).Scan(&historial).Error; err != nil {
		http.Error(w, "Error al consultar historial", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(historial)
}

func proveedorAutoFromPayload(payload ProveedorAutoDetalleDAO) (models.ProveedorAuto, error) {
	nombre := strings.TrimSpace(payload.Nombre)
	if nombre == "" {
		return models.ProveedorAuto{}, errors.New("Nombre/Razon social requerido")
	}
	ciNit := strings.TrimSpace(payload.CINIT)
	if len(ciNit) < 5 || len(ciNit) > 30 {
		return models.ProveedorAuto{}, errors.New("CI/NIT debe tener entre 5 y 30 caracteres")
	}
	telefono := strings.TrimSpace(payload.Telefono)
	if len(telefono) < 7 || len(telefono) > 12 || !regexp.MustCompile(`^[0-9+\-\s]+$`).MatchString(telefono) {
		return models.ProveedorAuto{}, errors.New("Telefono no valido")
	}
	email := strings.TrimSpace(payload.Email)
	if email != "" {
		if _, err := mail.ParseAddress(email); err != nil {
			return models.ProveedorAuto{}, errors.New("Email no valido")
		}
	}
	tipo := strings.TrimSpace(payload.Tipo)
	if tipo != "" && tipo != "Persona natural" && tipo != "Concesionaria" && tipo != "Importadora" {
		return models.ProveedorAuto{}, errors.New("Tipo de proveedor no valido")
	}
	direccion := strings.TrimSpace(payload.Direccion)
	if len(direccion) > 255 {
		return models.ProveedorAuto{}, errors.New("Direccion no valida")
	}
	observaciones := strings.TrimSpace(payload.Observaciones)
	if len(observaciones) > 500 {
		return models.ProveedorAuto{}, errors.New("Observaciones no valida")
	}
	estado := true
	if strings.TrimSpace(payload.Estado) != "" {
		parsed, err := functions.ActualizarEstado(payload.Estado)
		if err != nil {
			return models.ProveedorAuto{}, err
		}
		estado = parsed
	}
	return models.ProveedorAuto{
		Nombre:        nombre,
		CINIT:         ciNit,
		Telefono:      telefono,
		Email:         email,
		Direccion:     direccion,
		Tipo:          tipo,
		Observaciones: observaciones,
		Estado:        estado,
	}, nil
}

func ciNitProveedorDuplicado(ciNit string, excluirID uint) (bool, error) {
	ciNit = strings.TrimSpace(ciNit)
	if ciNit == "" {
		return false, nil
	}
	var count int64
	query := db.GDB.Model(&models.ProveedorAuto{}).Where("lower(trim(ci_nit)) = lower(trim(?))", ciNit)
	if excluirID > 0 {
		query = query.Where("id <> ?", excluirID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
