package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/security"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const (
	estadoPedidoRegistrado = "Pedido registrado"
	estadoPedidoTransito   = "En tránsito"
	estadoPedidoRecibido   = "Recibido"
	estadoPedidoCompletado = "Completado"
)

type PedidoVehiculoDAO struct {
	IDCliente            uint           `json:"id_cliente"`
	IDProveedor          uint           `json:"id_proveedor"`
	Fecha                string         `json:"fecha"`
	Marca                string         `json:"marca"`
	Modelo               string         `json:"modelo"`
	Anio                 uint           `json:"anio"`
	Color                string         `json:"color"`
	Version              string         `json:"version"`
	PaisOrigen           string         `json:"pais_origen"`
	PrecioEstimadoUSD    float64        `json:"precio_estimado_usd"`
	TipoCambio           float64        `json:"tipo_cambio"`
	FechaLlegadaEstimada string         `json:"fecha_llegada_estimada"`
	AdelantoRequeridoUSD float64        `json:"adelanto_requerido_usd"`
	AdelantoPorcentaje   float64        `json:"adelanto_porcentaje"`
	ValidezDias          uint           `json:"validez_dias"`
	Pagos                []PagoVentaDAO `json:"pagos"`
	Observacion          string         `json:"observacion"`
}

type PedidoVehiculoHistorialDAO struct {
	ID                   uint            `json:"id"`
	IDCliente            uint            `json:"id_cliente"`
	IDProveedor          *uint           `json:"id_proveedor"`
	IDVehiculo           *uint           `json:"id_vehiculo"`
	IDUsuario            *uint           `json:"id_usuario"`
	Fecha                string          `json:"fecha"`
	Cliente              string          `json:"cliente"`
	CICliente            string          `json:"ci_cliente"`
	Proveedor            string          `json:"proveedor"`
	Vehiculo             string          `json:"vehiculo"`
	Marca                string          `json:"marca"`
	Modelo               string          `json:"modelo"`
	Anio                 uint            `json:"anio"`
	Color                string          `json:"color"`
	Version              string          `json:"version"`
	PaisOrigen           string          `json:"pais_origen"`
	PrecioEstimadoUSD    float64         `json:"precio_estimado_usd"`
	TipoCambioUsado      float64         `json:"tipo_cambio_usado"`
	FechaLlegadaEstimada string          `json:"fecha_llegada_estimada"`
	AdelantoRequeridoUSD float64         `json:"adelanto_requerido_usd"`
	AdelantoPorcentaje   float64         `json:"adelanto_porcentaje"`
	AdelantoPagadoUSD    float64         `json:"adelanto_pagado_usd"`
	AdelantoPagadoBOB    float64         `json:"adelanto_pagado_bob"`
	SaldoPendienteUSD    float64         `json:"saldo_pendiente_usd"`
	Estado               string          `json:"estado"`
	FechaVencimiento     string          `json:"fecha_vencimiento"`
	FechaRecepcion       string          `json:"fecha_recepcion"`
	FechaCompletado      string          `json:"fecha_completado"`
	Observacion          string          `json:"observacion"`
	Pagos                json.RawMessage `json:"pagos"`
	DetallePago          string          `json:"detalle_pago"`
}

type RecibirPedidoVehiculoDAO struct {
	IDVehiculo  uint   `json:"id_vehiculo"`
	Observacion string `json:"observacion"`
}

type CompletarPedidoVehiculoDAO struct {
	TipoCambio  float64        `json:"tipo_cambio"`
	Pagos       []PagoVentaDAO `json:"pagos"`
	Observacion string         `json:"observacion"`
}

func ObtenerPedidosVehiculos(w http.ResponseWriter, r *http.Request) {
	idUsuario := r.URL.Query().Get("id_usuario")
	pedidos := make([]PedidoVehiculoHistorialDAO, 0)
	query := pedidosVehiculosQuery()
	args := []interface{}{}

	principal, ok := security.PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "Autenticacion requerida", http.StatusUnauthorized)
		return
	}
	if !security.CurrentUserHasRole(r, "admin", "encargado de ventas") {
		idUsuario = strconv.FormatUint(uint64(principal.ID), 10)
	}
	if idUsuario != "" {
		if _, err := strconv.Atoi(idUsuario); err != nil {
			http.Error(w, "Usuario no valido", http.StatusBadRequest)
			return
		}
		query += " where p.id_usuario = ?"
		args = append(args, idUsuario)
	}
	query += " order by p.fecha desc, p.id desc"

	if err := db.GDB.Raw(query, args...).Scan(&pedidos).Error; err != nil {
		http.Error(w, "Error al consultar pedidos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pedidos)
}

func ObtenerPedidoVehiculo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if !puedeAccederPedido(r, id) {
		http.Error(w, "Pedido no encontrado", http.StatusNotFound)
		return
	}
	pedidos := make([]PedidoVehiculoHistorialDAO, 0)
	if err := db.GDB.Raw(pedidosVehiculosQuery()+" where p.id = ? limit 1", id).Scan(&pedidos).Error; err != nil || len(pedidos) == 0 {
		http.Error(w, "Pedido no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pedidos[0])
}

func AgregarPedidoVehiculo(w http.ResponseWriter, r *http.Request) {
	var payload PedidoVehiculoDAO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Solicitud no valida", http.StatusBadRequest)
		return
	}
	principal, ok := security.PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "Autenticacion requerida", http.StatusUnauthorized)
		return
	}

	pedido, pagos, err := construirPedidoVehiculo(payload, principal.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&pedido).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al registrar pedido", http.StatusInternalServerError)
		return
	}
	for index := range pagos {
		pagos[index].PedidoID = pedido.ID
	}
	if len(pagos) > 0 {
		if err := tx.Create(&pagos).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al registrar pagos del pedido", http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&pedido)
}

func RecibirPedidoVehiculo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var payload RecibirPedidoVehiculoDAO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Solicitud no valida", http.StatusBadRequest)
		return
	}
	principal, ok := security.PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "Autenticacion requerida", http.StatusUnauthorized)
		return
	}
	if payload.IDVehiculo == 0 {
		http.Error(w, "Vehiculo requerido", http.StatusBadRequest)
		return
	}
	var vehiculo models.Vehiculo
	if err := db.GDB.Where("id = ?", payload.IDVehiculo).First(&vehiculo).Error; err != nil {
		http.Error(w, "Vehiculo no encontrado", http.StatusBadRequest)
		return
	}

	var pedido models.PedidoVehiculo
	if err := db.GDB.Where("id = ?", id).First(&pedido).Error; err != nil {
		http.Error(w, "Pedido no encontrado", http.StatusNotFound)
		return
	}
	if pedido.Estado == estadoPedidoCompletado {
		http.Error(w, "El pedido ya fue completado", http.StatusBadRequest)
		return
	}

	now := time.Now()
	pedido.IDVehiculo = &payload.IDVehiculo
	pedido.Estado = estadoPedidoRecibido
	pedido.FechaRecepcion = &now
	pedido.IDUsuarioRecibe = &principal.ID
	pedido.Observacion = appendObservacion(pedido.Observacion, payload.Observacion)

	if err := db.GDB.Save(&pedido).Error; err != nil {
		http.Error(w, "Error al recibir pedido", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&pedido)
}

func MarcarPedidoVehiculoEnTransito(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	principal, ok := security.PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "Autenticacion requerida", http.StatusUnauthorized)
		return
	}

	var pedido models.PedidoVehiculo
	if err := db.GDB.Where("id = ?", id).First(&pedido).Error; err != nil {
		http.Error(w, "Pedido no encontrado", http.StatusNotFound)
		return
	}
	if pedido.Estado != estadoPedidoRegistrado {
		http.Error(w, "Solo los pedidos registrados pueden pasar a transito", http.StatusBadRequest)
		return
	}

	pedido.Estado = estadoPedidoTransito
	pedido.Observacion = appendObservacion(pedido.Observacion, "Marcado en transito por usuario "+strconv.Itoa(int(principal.ID)))
	if err := db.GDB.Save(&pedido).Error; err != nil {
		http.Error(w, "Error al actualizar estado del pedido", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&pedido)
}

func CompletarPedidoVehiculo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var payload CompletarPedidoVehiculoDAO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Solicitud no valida", http.StatusBadRequest)
		return
	}
	principal, ok := security.PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "Autenticacion requerida", http.StatusUnauthorized)
		return
	}

	tx := db.GDB.Begin()
	var pedido models.PedidoVehiculo
	if err := tx.Where("id = ?", id).First(&pedido).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Pedido no encontrado", http.StatusNotFound)
		return
	}
	if pedido.Estado != estadoPedidoRecibido {
		tx.Rollback()
		http.Error(w, "El pedido debe estar recibido antes de completarse", http.StatusBadRequest)
		return
	}
	if pedido.IDVehiculo == nil {
		tx.Rollback()
		http.Error(w, "Debe vincular el vehiculo recibido", http.StatusBadRequest)
		return
	}

	pagosNormalizados, _, _, pagadoUSD, err := normalizarPagosVenta(payload.Pagos, roundMoney(payload.TipoCambio))
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if pagadoUSD != roundMoney(pedido.SaldoPendienteUSD) {
		tx.Rollback()
		http.Error(w, "El pago final debe completar el saldo pendiente", http.StatusBadRequest)
		return
	}

	pagos := make([]models.PagoPedidoVehiculo, 0, len(pagosNormalizados))
	for _, pago := range pagosNormalizados {
		pagos = append(pagos, models.PagoPedidoVehiculo{
			PedidoID: pedido.ID,
			Etapa:    "final",
			Moneda:   pago.Moneda,
			Metodo:   pago.Metodo,
			Monto:    roundMoney(pago.Monto),
		})
	}
	if err := tx.Create(&pagos).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al registrar pago final", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	pedido.Estado = estadoPedidoCompletado
	pedido.SaldoPendienteUSD = 0
	pedido.FechaCompletado = &now
	pedido.IDUsuarioCompleta = &principal.ID
	pedido.Observacion = appendObservacion(pedido.Observacion, payload.Observacion)
	if err := tx.Save(&pedido).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al completar pedido", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&pedido)
}

func construirPedidoVehiculo(payload PedidoVehiculoDAO, userID uint) (models.PedidoVehiculo, []models.PagoPedidoVehiculo, error) {
	if payload.IDCliente == 0 {
		return models.PedidoVehiculo{}, nil, errors.New("Cliente requerido")
	}
	if strings.TrimSpace(payload.Marca) == "" || strings.TrimSpace(payload.Modelo) == "" || payload.Anio == 0 || strings.TrimSpace(payload.PaisOrigen) == "" {
		return models.PedidoVehiculo{}, nil, errors.New("Datos del vehiculo solicitados incompletos")
	}
	if payload.PrecioEstimadoUSD <= 0 {
		return models.PedidoVehiculo{}, nil, errors.New("Precio estimado requerido")
	}
	tipoCambio := roundMoney(payload.TipoCambio)
	if tipoCambio <= 0 {
		return models.PedidoVehiculo{}, nil, errors.New("Tipo de cambio requerido")
	}

	fecha, err := time.Parse("2006-01-02", payload.Fecha)
	if err != nil {
		return models.PedidoVehiculo{}, nil, errors.New("Fecha no valida")
	}
	llegada, err := time.Parse("2006-01-02", payload.FechaLlegadaEstimada)
	if err != nil {
		return models.PedidoVehiculo{}, nil, errors.New("Fecha estimada de llegada no valida")
	}

	adelantoRequerido := roundMoney(payload.AdelantoRequeridoUSD)
	porcentaje := roundMoney(payload.AdelantoPorcentaje)
	if porcentaje > 0 {
		adelantoRequerido = roundMoney(payload.PrecioEstimadoUSD * porcentaje / 100)
	}
	if adelantoRequerido <= 0 {
		return models.PedidoVehiculo{}, nil, errors.New("Adelanto requerido")
	}
	if adelantoRequerido > payload.PrecioEstimadoUSD {
		return models.PedidoVehiculo{}, nil, errors.New("Adelanto mayor al precio estimado")
	}

	pagosNormalizados, totalUSD, totalBOB, pagadoUSD, err := normalizarPagosVenta(payload.Pagos, tipoCambio)
	if err != nil {
		return models.PedidoVehiculo{}, nil, err
	}
	if pagadoUSD > adelantoRequerido {
		return models.PedidoVehiculo{}, nil, errors.New("Pago mayor al adelanto requerido")
	}

	validez := payload.ValidezDias
	if validez == 0 {
		validez = 15
	}
	var idProveedor *uint
	if payload.IDProveedor != 0 {
		idProveedor = &payload.IDProveedor
	}
	idUsuario := userID

	pagos := make([]models.PagoPedidoVehiculo, 0, len(pagosNormalizados))
	for _, pago := range pagosNormalizados {
		pagos = append(pagos, models.PagoPedidoVehiculo{
			Etapa:  "adelanto",
			Moneda: pago.Moneda,
			Metodo: pago.Metodo,
			Monto:  roundMoney(pago.Monto),
		})
	}

	return models.PedidoVehiculo{
		IDCliente:            payload.IDCliente,
		IDProveedor:          idProveedor,
		IDUsuario:            &idUsuario,
		Fecha:                fecha,
		Marca:                strings.TrimSpace(payload.Marca),
		Modelo:               strings.TrimSpace(payload.Modelo),
		Anio:                 payload.Anio,
		Color:                strings.TrimSpace(payload.Color),
		Version:              strings.TrimSpace(payload.Version),
		PaisOrigen:           strings.TrimSpace(payload.PaisOrigen),
		PrecioEstimadoUSD:    roundMoney(payload.PrecioEstimadoUSD),
		TipoCambioUsado:      tipoCambio,
		FechaLlegadaEstimada: llegada,
		AdelantoRequeridoUSD: adelantoRequerido,
		AdelantoPorcentaje:   porcentaje,
		AdelantoPagadoUSD:    totalUSD,
		AdelantoPagadoBOB:    totalBOB,
		SaldoPendienteUSD:    roundMoney(payload.PrecioEstimadoUSD - pagadoUSD),
		Estado:               estadoPedidoRegistrado,
		FechaVencimiento:     fecha.AddDate(0, 0, int(validez)),
		Observacion:          strings.TrimSpace(payload.Observacion),
	}, pagos, nil
}

func puedeAccederPedido(r *http.Request, id interface{}) bool {
	principal, ok := security.PrincipalFromContext(r.Context())
	if !ok {
		return false
	}
	query := db.GDB.Model(&models.PedidoVehiculo{}).Where("id = ?", id)
	if !security.CurrentUserHasRole(r, "admin", "encargado de ventas") {
		query = query.Where("id_usuario = ?", principal.ID)
	}
	var count int64
	return query.Count(&count).Error == nil && count == 1
}

func appendObservacion(actual string, nueva string) string {
	nueva = strings.TrimSpace(nueva)
	if nueva == "" {
		return strings.TrimSpace(actual)
	}
	actual = strings.TrimSpace(actual)
	if actual == "" {
		return nueva
	}
	return actual + " | " + nueva
}

func pedidosVehiculosQuery() string {
	return `
		select
			p.id,
			p.id_cliente,
			p.id_proveedor,
			p.id_vehiculo,
			p.id_usuario,
			to_char(p.fecha, 'YYYY-MM-DD') as fecha,
			concat_ws(' ', c.nombre, c.apellido) as cliente,
			coalesce(c.ci, '') as ci_cliente,
			coalesce(pa.nombre, '') as proveedor,
			coalesce(nullif(concat_ws(' ', nullif(trim(v.marca), ''), nullif(trim(v.modelo), ''), nullif(v.anio::text, '0')), ''), v.nombre, '') as vehiculo,
			p.marca,
			p.modelo,
			p.anio,
			coalesce(p.color, '') as color,
			coalesce(p.version, '') as version,
			p.pais_origen,
			p.precio_estimado_usd,
			p.tipo_cambio_usado,
			to_char(p.fecha_llegada_estimada, 'YYYY-MM-DD') as fecha_llegada_estimada,
			p.adelanto_requerido_usd,
			p.adelanto_porcentaje,
			p.adelanto_pagado_usd,
			p.adelanto_pagado_bob,
			p.saldo_pendiente_usd,
			p.estado,
			to_char(p.fecha_vencimiento, 'YYYY-MM-DD') as fecha_vencimiento,
			coalesce(to_char(p.fecha_recepcion, 'YYYY-MM-DD'), '') as fecha_recepcion,
			coalesce(to_char(p.fecha_completado, 'YYYY-MM-DD'), '') as fecha_completado,
			coalesce(p.observacion, '') as observacion,
			coalesce(pp.pagos, '[]'::json) as pagos,
			coalesce(pp.detalle_pago, '') as detalle_pago
		from pedidos_vehiculos p
		inner join clientes c on c.id = p.id_cliente
		left join proveedores_autos pa on pa.id = p.id_proveedor
		left join vehiculos v on v.id = p.id_vehiculo
		left join lateral (
			select
				json_agg(json_build_object('id', pago.id, 'etapa', pago.etapa, 'moneda', pago.moneda, 'metodo', pago.metodo, 'monto', pago.monto) order by pago.id) as pagos,
				string_agg(concat(pago.etapa, ': ', pago.moneda, ' ', pago.metodo, ' ', pago.monto::text), ' | ' order by pago.id) as detalle_pago
			from pagos_pedido_vehiculo pago
			where pago.pedido_id = p.id
		) pp on true`
}
