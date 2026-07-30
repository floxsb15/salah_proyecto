package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/functions"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/querys"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type ClienteDAO struct {
	ID        uint   `json:"id"`
	Nombre    string `json:"nombre"`
	Apellido  string `json:"apellido"`
	CI        string `json:"ci"`
	Celular   string `json:"celular"`
	Direccion string `json:"direccion"`
	Estado    string `json:"estado"`
}

type ClienteCompraHistorialDAO struct {
	ID             uint    `json:"id"`
	IDVehiculo     uint    `json:"id_vehiculo"`
	Fecha          string  `json:"fecha"`
	Tipo           string  `json:"tipo"`
	MontoTotal     float64 `json:"monto_total"`
	MontoPagado    float64 `json:"monto_pagado"`
	SaldoPendiente float64 `json:"saldo_pendiente"`
	Estado         string  `json:"estado"`
	Vehiculo       string  `json:"vehiculo"`
	Categoria      string  `json:"categoria"`
	Segmento       string  `json:"segmento"`
	Observacion    string  `json:"observacion"`
}

type ClienteCuotaPendienteDAO struct {
	VentaID        uint    `json:"venta_id"`
	Vehiculo       string  `json:"vehiculo"`
	SaldoPendiente float64 `json:"saldo_pendiente"`
	Estado         string  `json:"estado"`
}

type ClienteHistorialComprasDAO struct {
	ClienteID          uint                        `json:"cliente_id"`
	Compras            []ClienteCompraHistorialDAO `json:"compras"`
	CuotasPendientes   []ClienteCuotaPendienteDAO  `json:"cuotas_pendientes"`
	TotalGastado       float64                     `json:"total_gastado"`
	TotalPendiente     float64                     `json:"total_pendiente"`
	TieneCreditoActivo bool                        `json:"tiene_credito_activo"`
}

func ObtenerClientes(w http.ResponseWriter, r *http.Request) {
	clientes := make([]ClienteDAO, 0)

	err := db.GDB.Raw(querys.Clientes).Scan(&clientes).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientes)
}

func ObtenerCliente(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var cliente ClienteDAO

	err := db.GDB.Raw(querys.Cliente, id).Scan(&cliente).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cliente)
}

func ObtenerHistorialComprasCliente(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var cliente models.Cliente
	if err := db.GDB.Where("id = ?", id).First(&cliente).Error; err != nil {
		http.Error(w, "Cliente no encontrado", http.StatusNotFound)
		return
	}

	compras := make([]ClienteCompraHistorialDAO, 0)
	query := `
		select
			vv.id,
			vv.id_vehiculo,
			to_char(vv.fecha, 'YYYY-MM-DD') as fecha,
			vv.tipo_venta as tipo,
			vv.precio_total as monto_total,
			case
				when lower(vv.estado_pago) = 'pagado completo' then vv.precio_total
				else 0
			end as monto_pagado,
			case
				when lower(vv.estado_pago) = 'pagado completo' then 0
				else vv.precio_total
			end as saldo_pendiente,
			vv.estado_venta as estado,
			coalesce(nullif(concat_ws(' ', nullif(trim(v.marca), ''), nullif(trim(v.modelo), ''), nullif(v.anio::text, '0')), ''), v.nombre) as vehiculo,
			coalesce(cat.nombre, '') as categoria,
			coalesce(seg.nombre, '') as segmento,
			coalesce(vv.observacion, '') as observacion
		from ventas_vehiculos vv
		inner join vehiculos v on v.id = vv.id_vehiculo
		left join categoria_vehiculo cat on cat.id = v.id_categoria
		left join segmento_vehiculo seg on seg.id = v.id_segmento
		where vv.id_cliente = ?
		order by vv.fecha desc, vv.id desc`

	if err := db.GDB.Raw(query, id).Scan(&compras).Error; err != nil {
		http.Error(w, "Error al consultar historial de ventas", http.StatusInternalServerError)
		return
	}

	cuotasPendientes := make([]ClienteCuotaPendienteDAO, 0)
	totalGastado := 0.0
	totalPendiente := 0.0
	tieneCreditoActivo := false

	for _, compra := range compras {
		totalGastado += compra.MontoTotal
		totalPendiente += compra.SaldoPendiente

		if compra.SaldoPendiente > 0 {
			tieneCreditoActivo = true
			cuotasPendientes = append(cuotasPendientes, ClienteCuotaPendienteDAO{
				VentaID:        compra.ID,
				Vehiculo:       compra.Vehiculo,
				SaldoPendiente: compra.SaldoPendiente,
				Estado:         compra.Estado,
			})
		}
	}

	res := ClienteHistorialComprasDAO{
		ClienteID:          cliente.ID,
		Compras:            compras,
		CuotasPendientes:   cuotasPendientes,
		TotalGastado:       totalGastado,
		TotalPendiente:     totalPendiente,
		TieneCreditoActivo: tieneCreditoActivo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func AgregarCliente(w http.ResponseWriter, r *http.Request) {
	var cliente ClienteDAO

	if err := json.NewDecoder(r.Body).Decode(&cliente); err != nil {
		http.Error(w, "Solicitud no valida", http.StatusBadRequest)
		return
	}

	cliente.CI = strings.TrimSpace(cliente.CI)
	if cliente.CI == "" {
		http.Error(w, "CI/NIT requerido", http.StatusBadRequest)
		return
	}

	ciDuplicado, err := ciClienteDuplicado(cliente.CI, 0)
	if err != nil {
		http.Error(w, "Error al validar CI/NIT", http.StatusInternalServerError)
		return
	}
	if ciDuplicado {
		http.Error(w, "Este CI/NIT ya existe", http.StatusConflict)
		return
	}

	nuevoEstado, err := functions.ActualizarEstado(cliente.Estado)
	if err != nil {
		http.Error(w, "Estado no valido", http.StatusBadRequest)
		return
	}

	nuevoCliente := models.Cliente{
		Nombre:    cliente.Nombre,
		Apellido:  cliente.Apellido,
		CI:        cliente.CI,
		Celular:   cliente.Celular,
		Direccion: cliente.Direccion,
		Estado:    nuevoEstado,
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&nuevoCliente).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al agregar Cliente", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&nuevoCliente)
}

func ModificarCliente(w http.ResponseWriter, r *http.Request) {
	idCliente := mux.Vars(r)["id"]
	var clienteExistente models.Cliente

	err := db.GDB.Where("id = ?", idCliente).First(&clienteExistente).Error
	if err != nil {
		http.Error(w, "Cliente no encontrado", http.StatusNotFound)
		return
	}

	var clienteActualizado ClienteDAO
	if err := json.NewDecoder(r.Body).Decode(&clienteActualizado); err != nil {
		http.Error(w, "Solicitud no valida", http.StatusBadRequest)
		return
	}

	clienteActualizado.CI = strings.TrimSpace(clienteActualizado.CI)
	if clienteActualizado.CI == "" {
		http.Error(w, "CI/NIT requerido", http.StatusBadRequest)
		return
	}

	ciDuplicado, err := ciClienteDuplicado(clienteActualizado.CI, clienteExistente.ID)
	if err != nil {
		http.Error(w, "Error al validar CI/NIT", http.StatusInternalServerError)
		return
	}
	if ciDuplicado {
		http.Error(w, "Este CI/NIT ya existe", http.StatusConflict)
		return
	}

	nuevoEstado, err := functions.ActualizarEstado(clienteActualizado.Estado)
	if err != nil {
		http.Error(w, "Estado no valido", http.StatusBadRequest)
		return
	}

	clienteExistente.Nombre = clienteActualizado.Nombre
	clienteExistente.Apellido = clienteActualizado.Apellido
	clienteExistente.CI = clienteActualizado.CI
	clienteExistente.Celular = clienteActualizado.Celular
	clienteExistente.Direccion = clienteActualizado.Direccion
	clienteExistente.Estado = nuevoEstado

	if err := db.GDB.Save(&clienteExistente).Error; err != nil {
		http.Error(w, "Error al actualizar Cliente", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&clienteExistente)
}

func ciClienteDuplicado(ci string, excluirID uint) (bool, error) {
	var count int64
	query := db.GDB.Model(&models.Cliente{}).
		Where("lower(trim(ci)) = lower(trim(?))", ci)

	if excluirID > 0 {
		query = query.Where("id <> ?", excluirID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
