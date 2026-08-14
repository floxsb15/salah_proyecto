package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/security"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProveedorAutoDAO struct {
	ID            uint   `json:"id"`
	Nombre        string `json:"nombre"`
	CINIT         string `json:"ci_nit"`
	Telefono      string `json:"telefono"`
	Email         string `json:"email"`
	Direccion     string `json:"direccion"`
	Tipo          string `json:"tipo"`
	Observaciones string `json:"observaciones"`
}

type CompraAutoDAO struct {
	IDVehiculo       uint             `json:"id_vehiculo"`
	IDProveedor      uint             `json:"id_proveedor"`
	Proveedor        ProveedorAutoDAO `json:"proveedor"`
	FechaCompra      string           `json:"fecha_compra"`
	MonedaPrecio     string           `json:"moneda_precio"`
	PrecioCompra     float64          `json:"precio_compra"`
	TipoCambio       float64          `json:"tipo_cambio"`
	GastoImportacion float64          `json:"gasto_importacion"`
	GastoTransporte  float64          `json:"gasto_transporte"`
	GastoPapeleo     float64          `json:"gasto_papeleo"`
	MetodoPago       string           `json:"metodo_pago"`
	Pagos            []PagoVentaDAO   `json:"pagos"`
	Observacion      string           `json:"observacion"`
}

type CompraAutoHistorialDAO struct {
	ID                   uint            `json:"id"`
	IDVehiculo           uint            `json:"id_vehiculo"`
	IDProveedor          uint            `json:"id_proveedor"`
	IDUsuario            *uint           `json:"id_usuario"`
	FechaCompra          string          `json:"fecha_compra"`
	Vehiculo             string          `json:"vehiculo"`
	Proveedor            string          `json:"proveedor"`
	ProveedorCINIT       string          `json:"proveedor_ci_nit"`
	ProveedorTelefono    string          `json:"proveedor_telefono"`
	PrecioCompraUSD      float64         `json:"precio_compra_usd"`
	PrecioCompraBOB      float64         `json:"precio_compra_bob"`
	MonedaPrecio         string          `json:"moneda_precio"`
	GastoImportacion     float64         `json:"gasto_importacion"`
	GastoTransporte      float64         `json:"gasto_transporte"`
	GastoPapeleo         float64         `json:"gasto_papeleo"`
	GastosAdicionales    float64         `json:"gastos_adicionales"`
	GastosAdicionalesBOB float64         `json:"gastos_adicionales_bob"`
	CostoTotalUSD        float64         `json:"costo_total_usd"`
	CostoTotalBOB        float64         `json:"costo_total_bob"`
	MetodoPago           string          `json:"metodo_pago"`
	EstadoPago           string          `json:"estado_pago"`
	TipoCambioUsado      float64         `json:"tipo_cambio_usado"`
	Pagos                json.RawMessage `json:"pagos"`
	DetallePago          string          `json:"detalle_pago"`
	Comprador            string          `json:"comprador"`
	Observacion          string          `json:"observacion"`
}

type CompletarPagoCompraAutoDAO struct {
	Pagos []PagoVentaDAO `json:"pagos"`
}

func ObtenerComprasAutos(w http.ResponseWriter, r *http.Request) {
	compras := make([]CompraAutoHistorialDAO, 0)
	query := `
		select
			ca.id,
			ca.id_vehiculo,
			ca.id_proveedor,
			ca.id_usuario,
			to_char(ca.fecha_compra, 'YYYY-MM-DD') as fecha_compra,
			coalesce(nullif(concat_ws(' ', nullif(trim(v.marca), ''), nullif(trim(v.modelo), ''), nullif(v.anio::text, '0')), ''), v.nombre) as vehiculo,
			pa.nombre as proveedor,
			coalesce(pa.ci_nit, '') as proveedor_ci_nit,
			coalesce(pa.telefono, '') as proveedor_telefono,
			ca.precio_compra_usd,
			ca.precio_compra_bob,
			ca.moneda_precio,
			ca.gasto_importacion,
			ca.gasto_transporte,
			ca.gasto_papeleo,
			ca.gastos_adicionales,
			round((ca.gastos_adicionales * ca.tipo_cambio_usado)::numeric, 2) as gastos_adicionales_bob,
			ca.costo_total_usd,
			round((ca.costo_total_usd * ca.tipo_cambio_usado)::numeric, 2) as costo_total_bob,
			case
				when coalesce(pc.cantidad, 0) > 1 then 'Mixto'
				when coalesce(pc.cantidad, 0) = 1 then pc.unico_metodo
				else ca.metodo_pago
			end as metodo_pago,
			ca.estado_pago,
			ca.tipo_cambio_usado,
			coalesce(pc.pagos, '[]'::json) as pagos,
			coalesce(pc.detalle_pago, '') as detalle_pago,
			coalesce(concat_ws(' ', u.nombre, u.apellido), '') as comprador,
			coalesce(ca.observacion, '') as observacion
		from compras_autos ca
		inner join vehiculos v on v.id = ca.id_vehiculo
		inner join proveedores_autos pa on pa.id = ca.id_proveedor
		left join usuarios u on u.id = ca.id_usuario
		left join lateral (
			select
				json_agg(json_build_object('id', p.id, 'moneda', p.moneda, 'metodo', p.metodo, 'monto', p.monto) order by p.id) as pagos,
				string_agg(concat(p.moneda, ' ', p.metodo, ' ', p.monto::text), ' | ' order by p.id) as detalle_pago,
				count(*) as cantidad,
				max(p.metodo) as unico_metodo
			from pagos_compra_auto p
			where p.compra_id = ca.id
		) pc on true
		order by ca.fecha_compra desc, ca.id desc`

	if err := db.GDB.Raw(query).Scan(&compras).Error; err != nil {
		http.Error(w, "Error al consultar compras", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(compras)
}

func AgregarCompraAuto(w http.ResponseWriter, r *http.Request) {
	var payload CompraAutoDAO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Solicitud no valida", http.StatusBadRequest)
		return
	}

	principal, ok := security.PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "Autenticacion requerida", http.StatusUnauthorized)
		return
	}

	compra, pagos, proveedor, err := construirCompraAuto(payload, principal.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	var vehiculo models.Vehiculo
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", compra.IDVehiculo).First(&vehiculo).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Vehiculo no encontrado", http.StatusBadRequest)
		return
	}
	if proveedor.ID == 0 {
		if err := tx.Create(&proveedor).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al registrar proveedor", http.StatusInternalServerError)
			return
		}
	}
	compra.IDProveedor = proveedor.ID
	if err := tx.Create(&compra).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al registrar compra", http.StatusInternalServerError)
		return
	}
	if len(pagos) > 0 {
		for i := range pagos {
			pagos[i].CompraID = compra.ID
		}
		if err := tx.Create(&pagos).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Error al registrar pagos de la compra", http.StatusInternalServerError)
			return
		}
	}
	vehiculo.PrecioCompra = compra.CostoTotalUSD
	if err := tx.Save(&vehiculo).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al actualizar costo del vehiculo", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&compra)
}

func CompletarPagoCompraAuto(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var payload CompletarPagoCompraAutoDAO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Detalle de pago no valido", http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	var compra models.CompraAuto
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&compra).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Compra no encontrada", http.StatusNotFound)
		return
	}
	if compra.EstadoPago == "Pagado completo" {
		tx.Rollback()
		http.Error(w, "La compra ya esta pagada", http.StatusBadRequest)
		return
	}
	if compra.TipoCambioUsado <= 0 {
		tx.Rollback()
		http.Error(w, "Tipo de cambio de compra no valido", http.StatusBadRequest)
		return
	}

	pagosNormalizados, _, _, nuevoPagadoUSD, err := normalizarPagosVenta(payload.Pagos, compra.TipoCambioUsado)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pagadoAnteriorUSD, err := totalPagadoCompraUSD(tx, compra.ID, compra.TipoCambioUsado)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al calcular pagos de la compra", http.StatusInternalServerError)
		return
	}
	saldoPendienteUSD := roundMoney(compra.CostoTotalUSD - pagadoAnteriorUSD)
	if nuevoPagadoUSD > saldoPendienteUSD {
		tx.Rollback()
		http.Error(w, "Pago mayor al saldo pendiente", http.StatusBadRequest)
		return
	}

	pagos := make([]models.PagoCompraAuto, 0, len(pagosNormalizados))
	for _, pago := range pagosNormalizados {
		pagos = append(pagos, models.PagoCompraAuto{
			CompraID: compra.ID,
			Moneda:   pago.Moneda,
			Metodo:   pago.Metodo,
			Monto:    roundMoney(pago.Monto),
		})
	}
	if err := tx.Create(&pagos).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al registrar pagos de la compra", http.StatusInternalServerError)
		return
	}

	pagadoTotalUSD := roundMoney(pagadoAnteriorUSD + nuevoPagadoUSD)
	compra.EstadoPago = "Pendiente"
	if pagadoTotalUSD >= compra.CostoTotalUSD {
		compra.EstadoPago = "Pagado completo"
	}
	compra.MetodoPago = metodoPagoCompraActual(tx, compra.ID)
	if err := tx.Save(&compra).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al actualizar compra", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&compra)
}

func construirCompraAuto(payload CompraAutoDAO, idUsuario uint) (models.CompraAuto, []models.PagoCompraAuto, models.ProveedorAuto, error) {
	if payload.IDVehiculo == 0 {
		return models.CompraAuto{}, nil, models.ProveedorAuto{}, errors.New("Vehiculo requerido")
	}
	fecha, err := time.Parse("2006-01-02", payload.FechaCompra)
	if err != nil {
		return models.CompraAuto{}, nil, models.ProveedorAuto{}, errors.New("Fecha de compra no valida")
	}
	moneda, err := normalizarMonedaPago(payload.MonedaPrecio)
	if err != nil {
		return models.CompraAuto{}, nil, models.ProveedorAuto{}, err
	}
	tipoCambio := roundMoney(payload.TipoCambio)
	if tipoCambio <= 0 {
		return models.CompraAuto{}, nil, models.ProveedorAuto{}, errors.New("Tipo de cambio requerido")
	}
	precio := roundMoney(payload.PrecioCompra)
	if precio <= 0 {
		return models.CompraAuto{}, nil, models.ProveedorAuto{}, errors.New("Precio de compra requerido")
	}

	precioUSD := precio
	precioBOB := roundMoney(precio * tipoCambio)
	if moneda == "BOB" {
		precioBOB = precio
		precioUSD = roundMoney(precio / tipoCambio)
	}

	importacion := nonNegativeMoney(payload.GastoImportacion)
	transporte := nonNegativeMoney(payload.GastoTransporte)
	papeleo := nonNegativeMoney(payload.GastoPapeleo)
	gastos := roundMoney(importacion + transporte + papeleo)
	costoTotal := roundMoney(precioUSD + gastos)

	metodoPago, err := normalizarMetodoPagoCompra(payload.MetodoPago)
	if err != nil {
		return models.CompraAuto{}, nil, models.ProveedorAuto{}, err
	}
	pagosNormalizados, totalUSD, totalBOB, pagadoUSD, err := normalizarPagosVenta(payload.Pagos, tipoCambio)
	if err != nil {
		return models.CompraAuto{}, nil, models.ProveedorAuto{}, err
	}
	if pagadoUSD > costoTotal {
		return models.CompraAuto{}, nil, models.ProveedorAuto{}, errors.New("Pago mayor al costo total")
	}
	estadoPago := "Pendiente"
	if pagadoUSD >= costoTotal {
		estadoPago = "Pagado completo"
	}
	if metodoPago != "Mixto" {
		metodoPago = resumenMetodoPago(pagosNormalizados)
	}
	_ = totalUSD
	_ = totalBOB

	proveedor, err := proveedorCompra(payload)
	if err != nil {
		return models.CompraAuto{}, nil, models.ProveedorAuto{}, err
	}
	idUsuarioPtr := &idUsuario
	pagos := make([]models.PagoCompraAuto, 0, len(pagosNormalizados))
	for _, pago := range pagosNormalizados {
		pagos = append(pagos, models.PagoCompraAuto{
			Moneda: pago.Moneda,
			Metodo: pago.Metodo,
			Monto:  roundMoney(pago.Monto),
		})
	}

	return models.CompraAuto{
		IDVehiculo:        payload.IDVehiculo,
		IDProveedor:       proveedor.ID,
		IDUsuario:         idUsuarioPtr,
		FechaCompra:       fecha,
		MonedaPrecio:      moneda,
		PrecioCompraUSD:   precioUSD,
		PrecioCompraBOB:   precioBOB,
		TipoCambioUsado:   tipoCambio,
		GastoImportacion:  importacion,
		GastoTransporte:   transporte,
		GastoPapeleo:      papeleo,
		GastosAdicionales: gastos,
		CostoTotalUSD:     costoTotal,
		MetodoPago:        metodoPago,
		EstadoPago:        estadoPago,
		Observacion:       strings.TrimSpace(payload.Observacion),
	}, pagos, proveedor, nil
}

func proveedorCompra(payload CompraAutoDAO) (models.ProveedorAuto, error) {
	if payload.IDProveedor > 0 {
		var proveedor models.ProveedorAuto
		if err := db.GDB.Where("id = ?", payload.IDProveedor).First(&proveedor).Error; err != nil {
			return models.ProveedorAuto{}, errors.New("Proveedor no encontrado")
		}
		return proveedor, nil
	}
	nombre := strings.TrimSpace(payload.Proveedor.Nombre)
	if nombre == "" {
		return models.ProveedorAuto{}, errors.New("Proveedor requerido")
	}
	return models.ProveedorAuto{
		Nombre:        nombre,
		CINIT:         strings.TrimSpace(payload.Proveedor.CINIT),
		Telefono:      strings.TrimSpace(payload.Proveedor.Telefono),
		Email:         strings.TrimSpace(payload.Proveedor.Email),
		Direccion:     strings.TrimSpace(payload.Proveedor.Direccion),
		Tipo:          strings.TrimSpace(payload.Proveedor.Tipo),
		Observaciones: strings.TrimSpace(payload.Proveedor.Observaciones),
		Estado:        true,
	}, nil
}

func normalizarMetodoPagoCompra(metodo string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(metodo)) {
	case "", "contado":
		return "Contado", nil
	case "credito", "crédito":
		return "Credito", nil
	case "mixto":
		return "Mixto", nil
	default:
		return "", errors.New("Metodo de pago de compra no valido")
	}
}

func totalPagadoCompraUSD(tx *gorm.DB, compraID uint, tipoCambio float64) (float64, error) {
	var pagos []models.PagoCompraAuto
	if err := tx.Where("compra_id = ?", compraID).Find(&pagos).Error; err != nil {
		return 0, err
	}
	total := 0.0
	for _, pago := range pagos {
		if pago.Moneda == "USD" {
			total = roundMoney(total + pago.Monto)
			continue
		}
		total = roundMoney(total + (pago.Monto / tipoCambio))
	}
	return total, nil
}

func metodoPagoCompraActual(tx *gorm.DB, compraID uint) string {
	var pagos []models.PagoCompraAuto
	if err := tx.Where("compra_id = ?", compraID).Order("id asc").Find(&pagos).Error; err != nil || len(pagos) == 0 {
		return "Efectivo"
	}
	baseMetodo := pagos[0].Metodo
	baseMoneda := pagos[0].Moneda
	for _, pago := range pagos {
		if pago.Metodo != baseMetodo || pago.Moneda != baseMoneda {
			return "Mixto"
		}
	}
	return baseMetodo
}

func nonNegativeMoney(value float64) float64 {
	rounded := roundMoney(value)
	if rounded < 0 {
		return 0
	}
	return rounded
}
