package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/security"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type ProformaVehicularDAO struct {
	ID                uint    `json:"id"`
	IDCliente         *uint   `json:"id_cliente"`
	IDVehiculo        uint    `json:"id_vehiculo"`
	IDUsuario         *uint   `json:"id_usuario"`
	Fecha             string  `json:"fecha"`
	ClienteNombre     string  `json:"cliente_nombre"`
	ClienteDireccion  string  `json:"cliente_direccion"`
	ClienteTelefono   string  `json:"cliente_telefono"`
	Modalidad         string  `json:"modalidad"`
	PrecioUnidad      float64 `json:"precio_unidad"`
	Cantidad          uint    `json:"cantidad"`
	PrecioTotal       float64 `json:"precio_total"`
	CuotaInicial      float64 `json:"cuota_inicial"`
	Saldo             float64 `json:"saldo"`
	ValidezDias       uint    `json:"validez_dias"`
	FechaVencimiento  string  `json:"fecha_vencimiento"`
	PrecioCatalogoRef float64 `json:"precio_catalogo_ref"`
	Vehiculo          string  `json:"vehiculo"`
	Vendedor          string  `json:"vendedor"`
}

type ProformaVehicularPayload struct {
	IDCliente        *uint   `json:"id_cliente"`
	IDVehiculo       uint    `json:"id_vehiculo"`
	ClienteNombre    string  `json:"cliente_nombre"`
	ClienteDireccion string  `json:"cliente_direccion"`
	ClienteTelefono  string  `json:"cliente_telefono"`
	Modalidad        string  `json:"modalidad"`
	PrecioUnidad     float64 `json:"precio_unidad"`
	Cantidad         uint    `json:"cantidad"`
	CuotaInicial     float64 `json:"cuota_inicial"`
	ValidezDias      uint    `json:"validez_dias"`
}

func ObtenerProformasVehiculares(w http.ResponseWriter, r *http.Request) {
	principal, ok := security.PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "Autenticacion requerida", http.StatusUnauthorized)
		return
	}

	proformas := make([]ProformaVehicularDAO, 0)
	query := `
		select
			p.id,
			p.id_cliente,
			p.id_vehiculo,
			p.id_usuario,
			to_char(p.fecha, 'YYYY-MM-DD HH24:MI:SS') as fecha,
			p.cliente_nombre,
			p.cliente_direccion,
			p.cliente_telefono,
			p.modalidad,
			p.precio_unidad,
			p.cantidad,
			p.precio_total,
			p.cuota_inicial,
			p.saldo,
			p.validez_dias,
			to_char(p.fecha_vencimiento, 'YYYY-MM-DD') as fecha_vencimiento,
			p.precio_catalogo_ref,
			coalesce(nullif(concat_ws(' ', nullif(trim(v.marca), ''), nullif(trim(v.modelo), ''), nullif(v.anio::text, '0')), ''), v.nombre) as vehiculo,
			coalesce(concat_ws(' ', u.nombre, u.apellido), '') as vendedor
		from proformas_vehiculares p
		inner join vehiculos v on v.id = p.id_vehiculo
		left join usuarios u on u.id = p.id_usuario`

	args := []any{}
	if !security.CurrentUserHasRole(r, "admin", "encargado de ventas") {
		query += " where p.id_usuario = ?"
		args = append(args, principal.ID)
	}
	query += " order by p.fecha desc, p.id desc"

	if err := db.GDB.Raw(query, args...).Scan(&proformas).Error; err != nil {
		http.Error(w, "Error al consultar proformas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proformas)
}

func AgregarProformaVehicular(w http.ResponseWriter, r *http.Request) {
	principal, ok := security.PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "Autenticacion requerida", http.StatusUnauthorized)
		return
	}

	var payload ProformaVehicularPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Solicitud no valida", http.StatusBadRequest)
		return
	}

	var vehiculo models.Vehiculo
	if err := db.GDB.Where("id = ?", payload.IDVehiculo).First(&vehiculo).Error; err != nil {
		http.Error(w, "Vehiculo no encontrado", http.StatusNotFound)
		return
	}
	if !vehiculo.Estado || vehiculo.CantidadDisponible == 0 {
		http.Error(w, "Vehiculo no disponible para proforma", http.StatusBadRequest)
		return
	}

	clienteNombre := strings.TrimSpace(payload.ClienteNombre)
	clienteDireccion := strings.TrimSpace(payload.ClienteDireccion)
	clienteTelefono := strings.TrimSpace(payload.ClienteTelefono)
	if payload.IDCliente != nil && *payload.IDCliente > 0 {
		var cliente models.Cliente
		if err := db.GDB.Where("id = ?", *payload.IDCliente).First(&cliente).Error; err != nil {
			http.Error(w, "Cliente no encontrado", http.StatusNotFound)
			return
		}
		if clienteNombre == "" {
			clienteNombre = strings.TrimSpace(strings.Join([]string{cliente.Nombre, cliente.Apellido}, " "))
		}
		if clienteDireccion == "" {
			clienteDireccion = cliente.Direccion
		}
		if clienteTelefono == "" {
			clienteTelefono = cliente.Celular
		}
	}
	if clienteNombre == "" {
		http.Error(w, "Nombre del cliente requerido", http.StatusBadRequest)
		return
	}

	modalidad := strings.TrimSpace(payload.Modalidad)
	if modalidad == "" {
		modalidad = "Almacen"
	}
	if payload.Cantidad == 0 {
		payload.Cantidad = 1
	}
	if payload.ValidezDias == 0 {
		payload.ValidezDias = 10
	}
	if payload.PrecioUnidad < vehiculo.Precio {
		http.Error(w, "El precio ofertado no puede ser menor al precio del catalogo", http.StatusBadRequest)
		return
	}

	precioTotal := payload.PrecioUnidad * float64(payload.Cantidad)
	if payload.CuotaInicial < 0 || payload.CuotaInicial > precioTotal {
		http.Error(w, "Cuota inicial no valida", http.StatusBadRequest)
		return
	}

	fecha := time.Now()
	idUsuario := principal.ID
	proforma := models.ProformaVehicular{
		IDCliente:         payload.IDCliente,
		IDVehiculo:        payload.IDVehiculo,
		IDUsuario:         &idUsuario,
		Fecha:             fecha,
		ClienteNombre:     clienteNombre,
		ClienteDireccion:  clienteDireccion,
		ClienteTelefono:   clienteTelefono,
		Modalidad:         modalidad,
		PrecioUnidad:      payload.PrecioUnidad,
		Cantidad:          payload.Cantidad,
		PrecioTotal:       precioTotal,
		CuotaInicial:      payload.CuotaInicial,
		Saldo:             precioTotal - payload.CuotaInicial,
		ValidezDias:       payload.ValidezDias,
		FechaVencimiento:  fecha.AddDate(0, 0, int(payload.ValidezDias)),
		PrecioCatalogoRef: vehiculo.Precio,
	}

	if err := db.GDB.Create(&proforma).Error; err != nil {
		http.Error(w, "Error al guardar la proforma", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proforma)
}
