package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

const (
	estadoVentaRegistrada    = "Registrada"
	estadoVentaCompletada    = "Completada"
	estadoVentaAnulada       = "Anulada"
	estadoVentaEnCredito     = "en_credito"
	estadoVentaPagado        = "pagado_completo"
	tipoVentaContado         = "Contado"
	tipoVentaCredito         = "Credito"
	tipoVentaCreditoDirecto  = "credito_directo"
	tipoVentaCreditoBancario = "credito_bancario"
	tipoVentaReserva         = "Reserva"
)

type VentaVehiculoDAO struct {
	IDCliente           uint           `json:"id_cliente"`
	IDVehiculo          uint           `json:"id_vehiculo"`
	IDUsuario           uint           `json:"id_usuario"`
	Fecha               string         `json:"fecha"`
	TipoVenta           string         `json:"tipo_venta"`
	Cantidad            uint           `json:"cantidad"`
	TipoCambio          float64        `json:"tipo_cambio"`
	PagoUSD             float64        `json:"pago_usd"`
	PagoBOB             float64        `json:"pago_bob"`
	Pagos               []PagoVentaDAO `json:"pagos"`
	ValidezProformaDias uint           `json:"validez_proforma_dias"`
	EstadoVenta         string         `json:"estado_venta"`
	EstadoPago          string         `json:"estado_pago"`
	MetodoPago          string         `json:"metodo_pago"`
	EstadoEntrega       string         `json:"estado_entrega"`
	FechaEntrega        string         `json:"fecha_entrega"`
	ReferenciaBancaria  string         `json:"referencia_bancaria"`
	EstadoDesembolso    string         `json:"estado_desembolso"`
	MontoReserva        float64        `json:"monto_reserva"`
	MontoInicial        float64        `json:"monto_inicial"`
	NumeroCuotas        uint           `json:"numero_cuotas"`
	FechaInicioCredito  string         `json:"fecha_inicio_credito"`
	FrecuenciaPago      string         `json:"frecuencia_pago"`
	TieneRespaldo       bool           `json:"tiene_respaldo"`
	TipoGarantia        string         `json:"tipo_garantia"`
	DocumentoGarantia   string         `json:"documento_garantia"`
	DatosGarante        string         `json:"datos_garante"`
	Observacion         string         `json:"observacion"`
}

type PagoVentaDAO struct {
	Moneda string  `json:"moneda"`
	Metodo string  `json:"metodo"`
	Monto  float64 `json:"monto"`
}

type VentaVehiculoEstadoDAO struct {
	EstadoVenta   string `json:"estado_venta"`
	EstadoPago    string `json:"estado_pago"`
	MetodoPago    string `json:"metodo_pago"`
	EstadoEntrega string `json:"estado_entrega"`
	FechaEntrega  string `json:"fecha_entrega"`
	Observacion   string `json:"observacion"`
}

type CompletarReservaDAO struct {
	MetodoPago    string  `json:"metodo_pago"`
	EstadoEntrega string  `json:"estado_entrega"`
	FechaEntrega  string  `json:"fecha_entrega"`
	MontoPago     float64 `json:"monto_pago"`
	IDUsuarioPago uint    `json:"id_usuario_pago"`
	Observacion   string  `json:"observacion"`
}

type VentaVehiculoHistorialDAO struct {
	ID                       uint            `json:"id"`
	IDCliente                uint            `json:"id_cliente"`
	IDVehiculo               uint            `json:"id_vehiculo"`
	IDUsuario                *uint           `json:"id_usuario"`
	Fecha                    string          `json:"fecha"`
	FechaVenta               string          `json:"fecha_venta"`
	TipoVenta                string          `json:"tipo_venta"`
	Cantidad                 uint            `json:"cantidad"`
	PrecioUnidad             float64         `json:"precio_unidad"`
	PrecioTotal              float64         `json:"precio_total"`
	PrecioUSD                float64         `json:"precio_usd"`
	TipoCambioUsado          float64         `json:"tipo_cambio_usado"`
	MontoBOBCalculado        float64         `json:"monto_bob_calculado"`
	PagoUSD                  float64         `json:"pago_usd"`
	PagoBOB                  float64         `json:"pago_bob"`
	Pagos                    json.RawMessage `json:"pagos"`
	DetallePago              string          `json:"detalle_pago"`
	SaldoBOB                 float64         `json:"saldo_bob"`
	CuotaInicial             float64         `json:"cuota_inicial"`
	Saldo                    float64         `json:"saldo"`
	ValidezProformaDias      uint            `json:"validez_proforma_dias"`
	FechaVencimientoProforma string          `json:"fecha_vencimiento_proforma"`
	EstadoVenta              string          `json:"estado_venta"`
	EstadoPago               string          `json:"estado_pago"`
	MetodoPago               string          `json:"metodo_pago"`
	EstadoEntrega            string          `json:"estado_entrega"`
	FechaEntrega             string          `json:"fecha_entrega"`
	ReferenciaBancaria       string          `json:"referencia_bancaria"`
	EstadoDesembolso         string          `json:"estado_desembolso"`
	MontoFinanciado          float64         `json:"monto_financiado"`
	NumeroCuotas             uint            `json:"numero_cuotas"`
	MontoCuota               float64         `json:"monto_cuota"`
	FechaInicioCredito       string          `json:"fecha_inicio_credito"`
	FrecuenciaPago           string          `json:"frecuencia_pago"`
	TieneRespaldo            bool            `json:"tiene_respaldo"`
	TipoGarantia             string          `json:"tipo_garantia"`
	DocumentoGarantia        string          `json:"documento_garantia"`
	DatosGarante             string          `json:"datos_garante"`
	Observacion              string          `json:"observacion"`
	ProformaVencida          bool            `json:"proforma_vencida"`
	Cliente                  string          `json:"cliente"`
	CICliente                string          `json:"ci_cliente"`
	Vehiculo                 string          `json:"vehiculo"`
	Categoria                string          `json:"categoria"`
	Segmento                 string          `json:"segmento"`
	Vendedor                 string          `json:"vendedor"`
	UsuarioPagoReserva       string          `json:"usuario_pago_reserva"`
}

type CuotaCreditoDAO struct {
	ID               uint    `json:"id"`
	IDVentaVehiculo  uint    `json:"id_venta_vehiculo"`
	Numero           uint    `json:"numero"`
	FechaVencimiento string  `json:"fecha_vencimiento"`
	Monto            float64 `json:"monto"`
	TipoCambioPago   float64 `json:"tipo_cambio_pago"`
	MontoBOBPagado   float64 `json:"monto_bob_pagado"`
	Estado           string  `json:"estado"`
	FechaPago        string  `json:"fecha_pago"`
	UsuarioPago      string  `json:"usuario_pago"`
}

type PagoCuotaCreditoDAO struct {
	TipoCambioPago float64 `json:"tipo_cambio_pago"`
	IDUsuarioPago  uint    `json:"id_usuario_pago"`
}

func AgregarVentaVehiculo(w http.ResponseWriter, r *http.Request) {
	venta, documentoGuardado, err := leerVentaVehiculoRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nuevaVenta, err := construirVentaVehiculo(venta)
	if err != nil {
		eliminarArchivoVenta(documentoGuardado)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(venta.Pagos) > 0 {
		pagosNormalizados, _, _, _, err := normalizarPagosVenta(venta.Pagos, venta.TipoCambio)
		if err != nil {
			eliminarArchivoVenta(documentoGuardado)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		venta.Pagos = pagosNormalizados
	}

	tx := db.GDB.Begin()
	if err := validarDisponibilidadVenta(tx, nuevaVenta.IDVehiculo, nuevaVenta.Cantidad, 0); err != nil {
		tx.Rollback()
		eliminarArchivoVenta(documentoGuardado)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ventaDescuentaStock(nuevaVenta.EstadoVenta) {
		if err := descontarCantidadVehiculo(tx, nuevaVenta.IDVehiculo, nuevaVenta.Cantidad); err != nil {
			tx.Rollback()
			eliminarArchivoVenta(documentoGuardado)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if err := tx.Create(&nuevaVenta).Error; err != nil {
		tx.Rollback()
		eliminarArchivoVenta(documentoGuardado)
		http.Error(w, "Error al registrar venta", http.StatusInternalServerError)
		return
	}
	if len(venta.Pagos) > 0 {
		pagos := make([]models.PagoVenta, 0, len(venta.Pagos))
		for _, pago := range venta.Pagos {
			pagos = append(pagos, models.PagoVenta{
				VentaID: nuevaVenta.ID,
				Moneda:  pago.Moneda,
				Metodo:  pago.Metodo,
				Monto:   roundMoney(pago.Monto),
			})
		}
		if err := tx.Create(&pagos).Error; err != nil {
			tx.Rollback()
			eliminarArchivoVenta(documentoGuardado)
			http.Error(w, "Error al registrar pagos de la venta", http.StatusInternalServerError)
			return
		}
	}
	if esVentaCredito(nuevaVenta.TipoVenta) {
		cuotas := generarCuotasCredito(nuevaVenta)
		if err := tx.Create(&cuotas).Error; err != nil {
			tx.Rollback()
			eliminarArchivoVenta(documentoGuardado)
			http.Error(w, "Error al generar cuotas de credito", http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&nuevaVenta)
}

func leerVentaVehiculoRequest(r *http.Request) (VentaVehiculoDAO, string, error) {
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if !strings.Contains(contentType, "multipart/form-data") {
		var venta VentaVehiculoDAO
		if err := json.NewDecoder(r.Body).Decode(&venta); err != nil {
			return VentaVehiculoDAO{}, "", err
		}
		return venta, "", nil
	}

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		return VentaVehiculoDAO{}, "", errors.New("Error al parsear el formulario")
	}

	documentoGarantia, err := guardarDocumentoGarantiaVenta(r)
	if err != nil {
		return VentaVehiculoDAO{}, "", err
	}

	venta, err := ventaVehiculoDAOFromForm(r)
	if err != nil {
		eliminarArchivoVenta(documentoGarantia)
		return VentaVehiculoDAO{}, "", err
	}
	venta.DocumentoGarantia = documentoGarantia
	return venta, documentoGarantia, nil
}

func ventaVehiculoDAOFromForm(r *http.Request) (VentaVehiculoDAO, error) {
	idCliente, err := parseUintFormValue(r.FormValue("id_cliente"))
	if err != nil {
		return VentaVehiculoDAO{}, errors.New("Cliente no valido")
	}
	idVehiculo, err := parseUintFormValue(r.FormValue("id_vehiculo"))
	if err != nil {
		return VentaVehiculoDAO{}, errors.New("Vehiculo no valido")
	}
	idUsuario, err := parseUintFormValueWithZero(r.FormValue("id_usuario"))
	if err != nil {
		return VentaVehiculoDAO{}, errors.New("Usuario no valido")
	}
	cantidad, err := parseUintFormValue(r.FormValue("cantidad"))
	if err != nil {
		return VentaVehiculoDAO{}, errors.New("Cantidad no valida")
	}
	tipoCambio, err := parseFloatFormValue(r.FormValue("tipo_cambio"))
	if err != nil {
		return VentaVehiculoDAO{}, errors.New("Tipo de cambio no valido")
	}
	pagoUSD, err := parseFloatFormValueWithZero(r.FormValue("pago_usd"))
	if err != nil {
		return VentaVehiculoDAO{}, errors.New("Pago USD no valido")
	}
	pagoBOB, err := parseFloatFormValueWithZero(r.FormValue("pago_bob"))
	if err != nil {
		return VentaVehiculoDAO{}, errors.New("Pago BOB no valido")
	}
	pagos, err := parsePagosVentaFormValue(r.FormValue("pagos"))
	if err != nil {
		return VentaVehiculoDAO{}, err
	}
	validez, err := parseUintFormValueWithZero(r.FormValue("validez_proforma_dias"))
	if err != nil {
		return VentaVehiculoDAO{}, errors.New("Validez de proforma no valida")
	}
	montoReserva, err := parseFloatFormValueWithZero(r.FormValue("monto_reserva"))
	if err != nil {
		return VentaVehiculoDAO{}, errors.New("Monto de reserva no valido")
	}
	montoInicial, err := parseFloatFormValueWithZero(r.FormValue("monto_inicial"))
	if err != nil {
		return VentaVehiculoDAO{}, errors.New("Monto inicial no valido")
	}
	numeroCuotas, err := parseUintFormValueWithZero(r.FormValue("numero_cuotas"))
	if err != nil {
		return VentaVehiculoDAO{}, errors.New("Numero de cuotas no valido")
	}

	return VentaVehiculoDAO{
		IDCliente:           idCliente,
		IDVehiculo:          idVehiculo,
		IDUsuario:           idUsuario,
		Fecha:               r.FormValue("fecha"),
		TipoVenta:           r.FormValue("tipo_venta"),
		Cantidad:            cantidad,
		TipoCambio:          tipoCambio,
		PagoUSD:             pagoUSD,
		PagoBOB:             pagoBOB,
		Pagos:               pagos,
		ValidezProformaDias: validez,
		EstadoVenta:         r.FormValue("estado_venta"),
		EstadoPago:          r.FormValue("estado_pago"),
		MetodoPago:          r.FormValue("metodo_pago"),
		EstadoEntrega:       r.FormValue("estado_entrega"),
		FechaEntrega:        r.FormValue("fecha_entrega"),
		ReferenciaBancaria:  r.FormValue("referencia_bancaria"),
		EstadoDesembolso:    r.FormValue("estado_desembolso"),
		MontoReserva:        montoReserva,
		MontoInicial:        montoInicial,
		NumeroCuotas:        numeroCuotas,
		FechaInicioCredito:  r.FormValue("fecha_inicio_credito"),
		FrecuenciaPago:      r.FormValue("frecuencia_pago"),
		TieneRespaldo:       parseBoolFormValue(r.FormValue("tiene_respaldo")),
		TipoGarantia:        r.FormValue("tipo_garantia"),
		DatosGarante:        r.FormValue("datos_garante"),
		Observacion:         r.FormValue("observacion"),
	}, nil
}

func guardarDocumentoGarantiaVenta(r *http.Request) (string, error) {
	if r.MultipartForm == nil {
		return "", nil
	}
	fileHeaders := r.MultipartForm.File["documento_garantia"]
	if len(fileHeaders) == 0 {
		return "", nil
	}
	if len(fileHeaders) > 1 {
		return "", errors.New("Solo se permite un documento de garantia")
	}

	fileHeader := fileHeaders[0]
	if fileHeader.Size > 10*1024*1024 {
		return "", errors.New("El documento de garantia no debe superar los 10MB")
	}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !extensionDocumentoGarantiaPermitida(ext) {
		return "", errors.New("Documento de garantia no valido")
	}
	if err := os.MkdirAll("internal/images/garantias", 0755); err != nil {
		return "", errors.New("Error al preparar la carpeta de garantias")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", errors.New("Error al obtener el documento de garantia")
	}
	defer file.Close()

	rutaDocumento := fmt.Sprintf("internal/images/garantias/garantia-%s%s", uuid.New().String(), ext)
	outFile, err := os.Create(rutaDocumento)
	if err != nil {
		return "", errors.New("Error al guardar el documento de garantia")
	}

	_, copyErr := io.Copy(outFile, file)
	closeErr := outFile.Close()
	if copyErr != nil || closeErr != nil {
		_ = os.Remove(rutaDocumento)
		return "", errors.New("Error al escribir el documento de garantia")
	}
	return rutaDocumento, nil
}

func extensionDocumentoGarantiaPermitida(ext string) bool {
	switch ext {
	case ".pdf", ".png", ".jpg", ".jpeg", ".webp":
		return true
	default:
		return false
	}
}

func eliminarArchivoVenta(path string) {
	if strings.TrimSpace(path) != "" {
		_ = os.Remove(path)
	}
}

func parseUintFormValue(value string) (uint, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, errors.New("valor requerido")
	}
	return parseUintFormValueWithZero(value)
}

func parseUintFormValueWithZero(value string) (uint, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
}

func parseFloatFormValue(value string) (float64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, errors.New("valor requerido")
	}
	return parseFloatFormValueWithZero(value)
}

func parseFloatFormValueWithZero(value string) (float64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, nil
	}
	return strconv.ParseFloat(value, 64)
}

func parseBoolFormValue(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "1", "si", "sí", "on":
		return true
	default:
		return false
	}
}

func parsePagosVentaFormValue(value string) ([]PagoVentaDAO, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	var pagos []PagoVentaDAO
	if err := json.Unmarshal([]byte(value), &pagos); err != nil {
		return nil, errors.New("Detalle de pago no valido")
	}
	return pagos, nil
}

func ObtenerCuotasCreditoVenta(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	cuotas := make([]CuotaCreditoDAO, 0)
	query := `
		select cc.id, cc.id_venta_vehiculo, cc.numero, to_char(cc.fecha_vencimiento, 'YYYY-MM-DD') as fecha_vencimiento,
			cc.monto,
			coalesce(cc.tipo_cambio_pago, 0) as tipo_cambio_pago,
			coalesce(cc.monto_bob_pagado, 0) as monto_bob_pagado,
			case when cc.estado = 'pendiente' and cc.fecha_vencimiento < current_date then 'atrasada' else cc.estado end as estado,
			coalesce(to_char(cc.fecha_pago, 'YYYY-MM-DD'), '') as fecha_pago,
			coalesce(concat_ws(' ', up.nombre, up.apellido), '') as usuario_pago
		from cuotas_credito cc
		left join usuarios up on up.id = cc.id_usuario_pago
		where cc.id_venta_vehiculo = ?
		order by cc.numero asc`

	if err := db.GDB.Raw(query, id).Scan(&cuotas).Error; err != nil {
		http.Error(w, "Error al consultar cuotas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cuotas)
}

func PagarCuotaCredito(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var payload PagoCuotaCreditoDAO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Tipo de cambio requerido", http.StatusBadRequest)
		return
	}
	tipoCambioPago := roundMoney(payload.TipoCambioPago)
	if tipoCambioPago <= 0 {
		http.Error(w, "Tipo de cambio requerido", http.StatusBadRequest)
		return
	}
	if payload.IDUsuarioPago == 0 {
		http.Error(w, "Usuario que acepta el pago requerido", http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	var cuota models.CuotaCredito
	if err := tx.Where("id = ?", id).First(&cuota).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Cuota no encontrada", http.StatusNotFound)
		return
	}
	if cuota.Estado == "pagada" {
		tx.Rollback()
		http.Error(w, "La cuota ya esta pagada", http.StatusBadRequest)
		return
	}

	now := time.Now()
	cuota.Estado = "pagada"
	cuota.FechaPago = &now
	cuota.TipoCambioPago = tipoCambioPago
	cuota.MontoBOBPagado = roundMoney(cuota.Monto * tipoCambioPago)
	cuota.IDUsuarioPago = &payload.IDUsuarioPago
	if err := tx.Save(&cuota).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al pagar cuota", http.StatusInternalServerError)
		return
	}
	if err := actualizarEstadoVentaPorCuotas(tx, cuota.IDVentaVehiculo); err != nil {
		tx.Rollback()
		http.Error(w, "Error al actualizar estado de venta", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&cuota)
}

func ObtenerVentasVehiculos(w http.ResponseWriter, r *http.Request) {
	idUsuario := r.URL.Query().Get("id_usuario")
	ventas := make([]VentaVehiculoHistorialDAO, 0)

	query := ventasVehiculosQuery()
	args := []interface{}{}
	if idUsuario != "" {
		if _, err := strconv.Atoi(idUsuario); err != nil {
			http.Error(w, "Usuario no valido", http.StatusBadRequest)
			return
		}
		query += " where vv.id_usuario = ?"
		args = append(args, idUsuario)
	}
	query += " order by coalesce(vv.fecha_venta, vv.fecha) desc, vv.id desc"

	if err := db.GDB.Raw(query, args...).Scan(&ventas).Error; err != nil {
		http.Error(w, "Error al consultar ventas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ventas)
}

func ObtenerVentaVehiculo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ventas := make([]VentaVehiculoHistorialDAO, 0)
	query := ventasVehiculosQuery() + " where vv.id = ? limit 1"

	if err := db.GDB.Raw(query, id).Scan(&ventas).Error; err != nil || len(ventas) == 0 {
		http.Error(w, "Venta no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ventas[0])
}

func ActualizarEstadoVentaVehiculo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var payload VentaVehiculoEstadoDAO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	var venta models.VentaVehiculo
	if err := tx.Where("id = ?", id).First(&venta).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Venta no encontrada", http.StatusNotFound)
		return
	}

	estadoAnterior := venta.EstadoVenta
	if payload.EstadoVenta != "" {
		estadoVenta, err := normalizarEstadoVenta(payload.EstadoVenta)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		venta.EstadoVenta = estadoVenta
	}
	if payload.EstadoPago != "" {
		venta.EstadoPago = payload.EstadoPago
	}
	if payload.MetodoPago != "" {
		metodoPago, err := normalizarMetodoPago(payload.MetodoPago)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		venta.MetodoPago = metodoPago
	}
	if payload.EstadoEntrega != "" {
		venta.EstadoEntrega = payload.EstadoEntrega
	}
	if payload.FechaEntrega != "" {
		fechaEntrega, err := time.Parse("2006-01-02", payload.FechaEntrega)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Fecha de entrega no valida", http.StatusBadRequest)
			return
		}
		venta.FechaEntrega = &fechaEntrega
	}
	if payload.Observacion != "" {
		venta.Observacion = payload.Observacion
	}

	if !ventaDescuentaStock(estadoAnterior) && ventaDescuentaStock(venta.EstadoVenta) {
		if err := descontarCantidadVehiculo(tx, venta.IDVehiculo, venta.Cantidad); err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	if ventaDescuentaStock(estadoAnterior) && !ventaDescuentaStock(venta.EstadoVenta) {
		if err := restaurarCantidadVehiculo(tx, venta.IDVehiculo, venta.Cantidad); err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Save(&venta).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al actualizar venta", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&venta)
}

func AnularVentaVehiculo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	tx := db.GDB.Begin()
	var venta models.VentaVehiculo
	if err := tx.Where("id = ?", id).First(&venta).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Venta no encontrada", http.StatusNotFound)
		return
	}
	if venta.EstadoVenta == estadoVentaAnulada {
		tx.Rollback()
		http.Error(w, "La venta ya esta anulada", http.StatusBadRequest)
		return
	}
	if ventaDescuentaStock(venta.EstadoVenta) {
		if err := restaurarCantidadVehiculo(tx, venta.IDVehiculo, venta.Cantidad); err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	venta.EstadoVenta = estadoVentaAnulada
	if err := tx.Save(&venta).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al anular venta", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&venta)
}

func CompletarReservaVehiculo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var payload CompletarReservaDAO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	var venta models.VentaVehiculo
	if err := tx.Where("id = ?", id).First(&venta).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Reserva no encontrada", http.StatusNotFound)
		return
	}
	if venta.TipoVenta != tipoVentaReserva {
		tx.Rollback()
		http.Error(w, "La venta seleccionada no es una reserva", http.StatusBadRequest)
		return
	}
	if venta.EstadoVenta == estadoVentaCompletada {
		tx.Rollback()
		http.Error(w, "La reserva ya fue completada", http.StatusBadRequest)
		return
	}
	if venta.EstadoVenta == estadoVentaAnulada {
		tx.Rollback()
		http.Error(w, "La reserva esta anulada", http.StatusBadRequest)
		return
	}

	montoPago := roundMoney(payload.MontoPago)
	saldoPendiente := roundMoney(venta.Saldo)
	if montoPago <= 0 {
		tx.Rollback()
		http.Error(w, "Monto de pago requerido", http.StatusBadRequest)
		return
	}
	if payload.IDUsuarioPago == 0 {
		tx.Rollback()
		http.Error(w, "Usuario que acepta el pago requerido", http.StatusBadRequest)
		return
	}
	if montoPago != saldoPendiente {
		tx.Rollback()
		http.Error(w, "El monto de pago debe completar el saldo pendiente", http.StatusBadRequest)
		return
	}
	metodoPago, err := normalizarMetodoPago(payload.MetodoPago)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := descontarCantidadVehiculo(tx, venta.IDVehiculo, venta.Cantidad); err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	venta.Saldo = 0
	venta.EstadoVenta = estadoVentaCompletada
	venta.EstadoPago = "Pagado completo"
	venta.MetodoPago = metodoPago
	venta.EstadoEntrega = strings.TrimSpace(payload.EstadoEntrega)
	if venta.EstadoEntrega == "" {
		venta.EstadoEntrega = "Pendiente"
	}
	if strings.TrimSpace(payload.FechaEntrega) != "" {
		fechaEntrega, err := time.Parse("2006-01-02", payload.FechaEntrega)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Fecha de entrega no valida", http.StatusBadRequest)
			return
		}
		venta.FechaEntrega = &fechaEntrega
	}
	if strings.TrimSpace(payload.Observacion) != "" {
		if strings.TrimSpace(venta.Observacion) != "" {
			venta.Observacion += " | "
		}
		venta.Observacion += strings.TrimSpace(payload.Observacion)
	}
	if payload.IDUsuarioPago != 0 {
		venta.IDUsuarioPagoReserva = &payload.IDUsuarioPago
	}

	if err := tx.Save(&venta).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al completar reserva", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&venta)
}

func construirVentaVehiculo(payload VentaVehiculoDAO) (models.VentaVehiculo, error) {
	if payload.IDCliente == 0 || payload.IDVehiculo == 0 {
		return models.VentaVehiculo{}, errors.New("Cliente y vehiculo son requeridos")
	}
	if payload.Cantidad == 0 {
		return models.VentaVehiculo{}, errors.New("Cantidad requerida")
	}

	fecha, err := time.Parse("2006-01-02", payload.Fecha)
	if err != nil {
		return models.VentaVehiculo{}, errors.New("Fecha no valida")
	}

	var vehiculo models.Vehiculo
	if err := db.GDB.Where("id = ?", payload.IDVehiculo).First(&vehiculo).Error; err != nil {
		return models.VentaVehiculo{}, errors.New("Vehiculo no encontrado")
	}
	if !vehiculo.Estado {
		return models.VentaVehiculo{}, errors.New("Vehiculo no disponible")
	}

	tipoVenta, err := normalizarTipoVenta(payload.TipoVenta)
	if err != nil {
		return models.VentaVehiculo{}, err
	}
	estadoVenta := estadoVentaRegistrada
	if strings.TrimSpace(payload.EstadoVenta) != "" {
		estadoVenta, err = normalizarEstadoVenta(payload.EstadoVenta)
		if err != nil {
			return models.VentaVehiculo{}, err
		}
	}

	validez := payload.ValidezProformaDias
	if validez == 0 {
		validez = 15
	}
	precioUnidad := roundMoney(vehiculo.Precio)
	precioTotal := roundMoney(precioUnidad * float64(payload.Cantidad))
	tipoCambio := roundMoney(payload.TipoCambio)
	if tipoCambio <= 0 {
		return models.VentaVehiculo{}, errors.New("Tipo de cambio requerido")
	}
	montoBOBCalculado := roundMoney(precioTotal * tipoCambio)
	pagoUSD := roundMoney(payload.PagoUSD)
	pagoBOB := roundMoney(payload.PagoBOB)
	cuotaInicial := 0.0
	saldo := 0.0
	if esVentaCredito(tipoVenta) {
		cuotaInicial = roundMoney(payload.MontoInicial)
		if cuotaInicial <= 0 {
			cuotaInicial = roundMoney(precioTotal * 0.30)
		}
		if cuotaInicial >= precioTotal {
			return models.VentaVehiculo{}, errors.New("Monto inicial debe ser menor al precio total")
		}
		saldo = roundMoney(precioTotal - cuotaInicial)
		if payload.NumeroCuotas == 0 {
			return models.VentaVehiculo{}, errors.New("Numero de cuotas requerido")
		}
		if strings.TrimSpace(payload.FechaInicioCredito) == "" {
			return models.VentaVehiculo{}, errors.New("Fecha de inicio de credito requerida")
		}
		if _, err := normalizarFrecuenciaPago(payload.FrecuenciaPago); err != nil {
			return models.VentaVehiculo{}, err
		}
		estadoVenta = estadoVentaEnCredito
		payload.EstadoPago = "Pendiente"
		if tipoVenta == tipoVentaCreditoDirecto {
			if !payload.TieneRespaldo || strings.TrimSpace(payload.TipoGarantia) == "" || strings.TrimSpace(payload.DatosGarante) == "" {
				return models.VentaVehiculo{}, errors.New("Credito directo requiere respaldo, tipo de garantia y datos del garante")
			}
		}
	}
	if tipoVenta == tipoVentaReserva {
		montoReserva := roundMoney(payload.MontoReserva)
		if montoReserva <= 0 {
			return models.VentaVehiculo{}, errors.New("Monto de reserva requerido")
		}
		if montoReserva > precioTotal {
			return models.VentaVehiculo{}, errors.New("Monto de reserva mayor al precio total")
		}
		cuotaInicial = montoReserva
		if len(payload.Pagos) > 0 {
			pagosNormalizados, totalUSD, totalBOB, pagadoEquivalenteUSD, err := normalizarPagosVenta(payload.Pagos, tipoCambio)
			if err != nil {
				return models.VentaVehiculo{}, err
			}
			if pagadoEquivalenteUSD > montoReserva {
				return models.VentaVehiculo{}, errors.New("Pago mayor al monto de reserva")
			}
			payload.Pagos = pagosNormalizados
			pagoUSD = totalUSD
			pagoBOB = totalBOB
			cuotaInicial = pagadoEquivalenteUSD
			payload.MetodoPago = resumenMetodoPago(pagosNormalizados)
			if pagadoEquivalenteUSD >= montoReserva {
				payload.EstadoPago = "Pagado completo"
			} else {
				payload.EstadoPago = "Parcial"
			}
		}
		saldo = roundMoney(precioTotal - cuotaInicial)
		estadoVenta = estadoVentaRegistrada
		if strings.TrimSpace(payload.EstadoPago) == "" {
			payload.EstadoPago = "Pagado completo"
		}
		payload.EstadoEntrega = "Pendiente"
		payload.FechaEntrega = ""
	}

	var idUsuario *uint
	if payload.IDUsuario != 0 {
		idUsuario = &payload.IDUsuario
	}
	var fechaEntrega *time.Time
	if strings.TrimSpace(payload.FechaEntrega) != "" {
		parsed, err := time.Parse("2006-01-02", payload.FechaEntrega)
		if err != nil {
			return models.VentaVehiculo{}, errors.New("Fecha de entrega no valida")
		}
		fechaEntrega = &parsed
	}

	estadoPago := strings.TrimSpace(payload.EstadoPago)
	if estadoPago == "" {
		estadoPago = "Pendiente"
	}
	metodoPago, err := normalizarMetodoPago(payload.MetodoPago)
	if err != nil {
		return models.VentaVehiculo{}, err
	}
	saldoBOB := roundMoney(saldo * tipoCambio)
	if metodoPago == "Mixto" && tipoVenta != tipoVentaReserva && !esVentaCredito(tipoVenta) {
		pagosNormalizados, totalUSD, totalBOB, pagadoEquivalenteUSD, err := normalizarPagosVenta(payload.Pagos, tipoCambio)
		if err != nil {
			return models.VentaVehiculo{}, err
		}
		if pagadoEquivalenteUSD > precioTotal {
			return models.VentaVehiculo{}, errors.New("Pago mixto mayor al precio total")
		}
		payload.Pagos = pagosNormalizados
		pagoUSD = totalUSD
		pagoBOB = totalBOB
		metodoPago = resumenMetodoPago(pagosNormalizados)
		cuotaInicial = pagadoEquivalenteUSD
		saldo = roundMoney(precioTotal - pagadoEquivalenteUSD)
		saldoBOB = roundMoney(saldo * tipoCambio)
		if saldo == 0 {
			estadoVenta = estadoVentaCompletada
			estadoPago = "Pagado completo"
		} else if pagadoEquivalenteUSD > 0 {
			estadoVenta = estadoVentaRegistrada
			estadoPago = "Parcial"
		} else {
			estadoVenta = estadoVentaRegistrada
			estadoPago = "Pendiente"
		}
	}
	fechaInicioCredito, frecuenciaPago, numeroCuotas, montoFinanciado, montoCuota := planCreditoValores(payload, tipoVenta, saldo)
	estadoEntrega := strings.TrimSpace(payload.EstadoEntrega)
	if estadoEntrega == "" {
		estadoEntrega = "Pendiente"
	}

	return models.VentaVehiculo{
		IDCliente:                payload.IDCliente,
		IDVehiculo:               payload.IDVehiculo,
		IDUsuario:                idUsuario,
		Fecha:                    fecha,
		FechaVenta:               time.Now(),
		TipoVenta:                tipoVenta,
		Cantidad:                 payload.Cantidad,
		PrecioUnidad:             precioUnidad,
		PrecioTotal:              precioTotal,
		PrecioUSD:                precioTotal,
		TipoCambioUsado:          tipoCambio,
		MontoBOBCalculado:        montoBOBCalculado,
		PagoUSD:                  pagoUSD,
		PagoBOB:                  pagoBOB,
		SaldoBOB:                 saldoBOB,
		CuotaInicial:             cuotaInicial,
		Saldo:                    saldo,
		MontoFinanciado:          montoFinanciado,
		NumeroCuotas:             numeroCuotas,
		MontoCuota:               montoCuota,
		FechaInicioCredito:       fechaInicioCredito,
		FrecuenciaPago:           frecuenciaPago,
		TieneRespaldo:            payload.TieneRespaldo,
		TipoGarantia:             strings.TrimSpace(payload.TipoGarantia),
		DocumentoGarantia:        strings.TrimSpace(payload.DocumentoGarantia),
		DatosGarante:             strings.TrimSpace(payload.DatosGarante),
		ValidezProformaDias:      validez,
		FechaVencimientoProforma: fecha.AddDate(0, 0, int(validez)),
		EstadoVenta:              estadoVenta,
		EstadoPago:               estadoPago,
		MetodoPago:               metodoPago,
		EstadoEntrega:            estadoEntrega,
		FechaEntrega:             fechaEntrega,
		ReferenciaBancaria:       strings.TrimSpace(payload.ReferenciaBancaria),
		EstadoDesembolso:         strings.TrimSpace(payload.EstadoDesembolso),
		Observacion:              strings.TrimSpace(payload.Observacion),
	}, nil
}

func esVentaCredito(tipoVenta string) bool {
	return tipoVenta == tipoVentaCredito || tipoVenta == tipoVentaCreditoDirecto || tipoVenta == tipoVentaCreditoBancario
}

func ventaDescuentaStock(estadoVenta string) bool {
	return estadoVenta == estadoVentaCompletada || estadoVenta == estadoVentaEnCredito || estadoVenta == estadoVentaPagado
}

func planCreditoValores(payload VentaVehiculoDAO, tipoVenta string, saldo float64) (*time.Time, string, uint, float64, float64) {
	if !esVentaCredito(tipoVenta) {
		return nil, "", 0, 0, 0
	}

	fechaInicio, _ := time.Parse("2006-01-02", payload.FechaInicioCredito)
	frecuencia, _ := normalizarFrecuenciaPago(payload.FrecuenciaPago)
	numeroCuotas := payload.NumeroCuotas
	montoFinanciado := roundMoney(saldo)
	montoCuota := roundMoney(montoFinanciado / float64(numeroCuotas))
	return &fechaInicio, frecuencia, numeroCuotas, montoFinanciado, montoCuota
}

func generarCuotasCredito(venta models.VentaVehiculo) []models.CuotaCredito {
	if venta.NumeroCuotas == 0 || venta.FechaInicioCredito == nil {
		return nil
	}

	cuotas := make([]models.CuotaCredito, 0, venta.NumeroCuotas)
	acumulado := 0.0
	for i := uint(1); i <= venta.NumeroCuotas; i++ {
		monto := venta.MontoCuota
		if i == venta.NumeroCuotas {
			monto = roundMoney(venta.MontoFinanciado - acumulado)
		}
		acumulado = roundMoney(acumulado + monto)
		cuotas = append(cuotas, models.CuotaCredito{
			IDVentaVehiculo:  venta.ID,
			Numero:           i,
			FechaVencimiento: calcularVencimientoCuota(*venta.FechaInicioCredito, venta.FrecuenciaPago, i),
			Monto:            monto,
			Estado:           "pendiente",
		})
	}
	return cuotas
}

func calcularVencimientoCuota(fechaInicio time.Time, frecuencia string, numero uint) time.Time {
	offset := int(numero - 1)
	switch frecuencia {
	case "semanal":
		return fechaInicio.AddDate(0, 0, offset*7)
	case "quincenal":
		return fechaInicio.AddDate(0, 0, offset*15)
	case "mensual":
		return fechaInicio.AddDate(0, offset, 0)
	default:
		return fechaInicio.AddDate(0, offset, 0)
	}
}

func actualizarEstadoVentaPorCuotas(tx *gorm.DB, idVenta uint) error {
	var pendientes int64
	if err := tx.Model(&models.CuotaCredito{}).
		Where("id_venta_vehiculo = ? and estado <> ?", idVenta, "pagada").
		Count(&pendientes).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{}
	if pendientes == 0 {
		updates["estado_venta"] = estadoVentaPagado
		updates["estado_pago"] = "Pagado completo"
		updates["saldo"] = 0
		updates["saldo_bob"] = 0
	} else {
		updates["estado_venta"] = estadoVentaEnCredito
		updates["estado_pago"] = "Pendiente"
	}
	return tx.Model(&models.VentaVehiculo{}).Where("id = ?", idVenta).Updates(updates).Error
}

func validarDisponibilidadVenta(tx *gorm.DB, idVehiculo uint, cantidad uint, excluirVentaID uint) error {
	var vehiculo models.Vehiculo
	if err := tx.Where("id = ?", idVehiculo).First(&vehiculo).Error; err != nil {
		return errors.New("Vehiculo no encontrado")
	}

	reservada, err := cantidadReservadaActiva(tx, idVehiculo, excluirVentaID)
	if err != nil {
		return errors.New("Error al validar disponibilidad")
	}
	disponible := int64(vehiculo.CantidadDisponible) - int64(reservada)
	if int64(cantidad) > disponible {
		return errors.New("Cantidad mayor a la disponibilidad del vehiculo")
	}
	return nil
}

func cantidadReservadaActiva(tx *gorm.DB, idVehiculo uint, excluirVentaID uint) (uint, error) {
	var total uint
	query := tx.Model(&models.VentaVehiculo{}).
		Select("coalesce(sum(cantidad), 0)").
		Where("id_vehiculo = ? and estado_venta = ? and fecha_vencimiento_proforma >= ?", idVehiculo, estadoVentaRegistrada, time.Now().Format("2006-01-02"))
	if excluirVentaID != 0 {
		query = query.Where("id <> ?", excluirVentaID)
	}
	if err := query.Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func descontarCantidadVehiculo(tx *gorm.DB, idVehiculo uint, cantidad uint) error {
	var vehiculo models.Vehiculo
	if err := tx.Where("id = ?", idVehiculo).First(&vehiculo).Error; err != nil {
		return errors.New("Vehiculo no encontrado")
	}
	if vehiculo.CantidadDisponible < cantidad {
		return errors.New("Cantidad mayor a la disponibilidad del vehiculo")
	}
	vehiculo.CantidadDisponible -= cantidad
	return tx.Save(&vehiculo).Error
}

func restaurarCantidadVehiculo(tx *gorm.DB, idVehiculo uint, cantidad uint) error {
	return tx.Model(&models.Vehiculo{}).Where("id = ?", idVehiculo).
		UpdateColumn("cantidad_disponible", gorm.Expr("cantidad_disponible + ?", cantidad)).Error
}

func normalizarTipoVenta(tipo string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(tipo)) {
	case "contado":
		return tipoVentaContado, nil
	case "credito", "crédito":
		return tipoVentaCredito, nil
	case "credito_directo", "credito directo", "crédito directo":
		return tipoVentaCreditoDirecto, nil
	case "credito_bancario", "credito bancario", "crédito bancario":
		return tipoVentaCreditoBancario, nil
	case "reserva":
		return tipoVentaReserva, nil
	default:
		return "", errors.New("Tipo de venta no valido")
	}
}

func normalizarEstadoVenta(estado string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(estado)) {
	case "registrada":
		return estadoVentaRegistrada, nil
	case "completada":
		return estadoVentaCompletada, nil
	case "anulada":
		return estadoVentaAnulada, nil
	case "en_credito", "en credito":
		return estadoVentaEnCredito, nil
	case "pagado_completo", "pagado completo":
		return estadoVentaPagado, nil
	default:
		return "", errors.New("Estado de venta no valido")
	}
}

func normalizarFrecuenciaPago(frecuencia string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(frecuencia)) {
	case "semanal":
		return "semanal", nil
	case "quincenal":
		return "quincenal", nil
	case "", "mensual":
		return "mensual", nil
	default:
		return "", errors.New("Frecuencia de pago no valida")
	}
}

func normalizarMetodoPago(metodo string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(metodo)) {
	case "", "efectivo":
		return "Efectivo", nil
	case "qr":
		return "QR", nil
	case "transferencia":
		return "Transferencia", nil
	case "mixto":
		return "Mixto", nil
	default:
		return "", errors.New("Metodo de pago no valido")
	}
}

func normalizarMetodoPagoSimple(metodo string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(metodo)) {
	case "", "efectivo":
		return "Efectivo", nil
	case "qr":
		return "QR", nil
	case "transferencia":
		return "Transferencia", nil
	default:
		return "", errors.New("Metodo de pago no valido")
	}
}

func normalizarMonedaPago(moneda string) (string, error) {
	switch strings.ToUpper(strings.TrimSpace(moneda)) {
	case "USD":
		return "USD", nil
	case "BOB":
		return "BOB", nil
	default:
		return "", errors.New("Moneda de pago no valida")
	}
}

func normalizarPagosVenta(pagos []PagoVentaDAO, tipoCambio float64) ([]PagoVentaDAO, float64, float64, float64, error) {
	if len(pagos) == 0 {
		return nil, 0, 0, 0, errors.New("Agregue al menos una linea en el detalle de pago")
	}

	normalizados := make([]PagoVentaDAO, 0, len(pagos))
	totalUSD := 0.0
	totalBOB := 0.0
	for _, pago := range pagos {
		moneda, err := normalizarMonedaPago(pago.Moneda)
		if err != nil {
			return nil, 0, 0, 0, err
		}
		metodo, err := normalizarMetodoPagoSimple(pago.Metodo)
		if err != nil {
			return nil, 0, 0, 0, err
		}
		monto := roundMoney(pago.Monto)
		if monto <= 0 {
			return nil, 0, 0, 0, errors.New("Cada linea de pago debe tener monto mayor a cero")
		}
		if moneda == "USD" {
			totalUSD = roundMoney(totalUSD + monto)
		} else {
			totalBOB = roundMoney(totalBOB + monto)
		}
		normalizados = append(normalizados, PagoVentaDAO{
			Moneda: moneda,
			Metodo: metodo,
			Monto:  monto,
		})
	}

	pagadoEquivalenteUSD := roundMoney(totalUSD + (totalBOB / tipoCambio))
	return normalizados, totalUSD, totalBOB, pagadoEquivalenteUSD, nil
}

func resumenMetodoPago(pagos []PagoVentaDAO) string {
	if len(pagos) == 0 {
		return "Efectivo"
	}
	resumen := pagos[0].Metodo
	for _, pago := range pagos {
		if pago.Metodo != resumen || pago.Moneda != pagos[0].Moneda {
			return "Mixto"
		}
	}
	return resumen
}

func roundMoney(value float64) float64 {
	return math.Round(value*100) / 100
}

func ventasVehiculosQuery() string {
	return `
		select
			vv.id,
			vv.id_cliente,
			vv.id_vehiculo,
			vv.id_usuario,
			to_char(vv.fecha, 'YYYY-MM-DD') as fecha,
			to_char(coalesce(vv.fecha_venta, vv.fecha), 'YYYY-MM-DD HH24:MI:SS') as fecha_venta,
			vv.tipo_venta,
			vv.cantidad,
			vv.precio_unidad,
			vv.precio_total,
			coalesce(vv.precio_usd, vv.precio_total) as precio_usd,
			coalesce(vv.tipo_cambio_usado, 0) as tipo_cambio_usado,
			coalesce(vv.monto_bob_calculado, 0) as monto_bob_calculado,
			coalesce(pv.total_usd, vv.pago_usd, 0) as pago_usd,
			coalesce(pv.total_bob, vv.pago_bob, 0) as pago_bob,
			coalesce(pv.pagos, '[]'::json) as pagos,
			coalesce(pv.detalle_pago, '') as detalle_pago,
			coalesce(vv.saldo_bob, 0) as saldo_bob,
			vv.cuota_inicial,
			vv.saldo,
			vv.validez_proforma_dias,
			to_char(vv.fecha_vencimiento_proforma, 'YYYY-MM-DD') as fecha_vencimiento_proforma,
			vv.estado_venta,
			vv.estado_pago,
			case
				when coalesce(pv.cantidad, 0) > 1 then 'Mixto'
				when coalesce(pv.cantidad, 0) = 1 then pv.unico_metodo
				else coalesce(vv.metodo_pago, 'Efectivo')
			end as metodo_pago,
			vv.estado_entrega,
			coalesce(to_char(vv.fecha_entrega, 'YYYY-MM-DD'), '') as fecha_entrega,
			coalesce(vv.referencia_bancaria, '') as referencia_bancaria,
			coalesce(vv.estado_desembolso, '') as estado_desembolso,
			coalesce(vv.monto_financiado, 0) as monto_financiado,
			coalesce(vv.numero_cuotas, 0) as numero_cuotas,
			coalesce(vv.monto_cuota, 0) as monto_cuota,
			coalesce(to_char(vv.fecha_inicio_credito, 'YYYY-MM-DD'), '') as fecha_inicio_credito,
			coalesce(vv.frecuencia_pago, '') as frecuencia_pago,
			coalesce(vv.tiene_respaldo, false) as tiene_respaldo,
			coalesce(vv.tipo_garantia, '') as tipo_garantia,
			coalesce(vv.documento_garantia, '') as documento_garantia,
			coalesce(vv.datos_garante, '') as datos_garante,
			coalesce(vv.observacion, '') as observacion,
			(vv.estado_venta = 'Registrada' and vv.fecha_vencimiento_proforma < current_date) as proforma_vencida,
			concat_ws(' ', c.nombre, c.apellido) as cliente,
			coalesce(c.ci, '') as ci_cliente,
			coalesce(nullif(concat_ws(' ', nullif(trim(v.marca), ''), nullif(trim(v.modelo), ''), nullif(v.anio::text, '0')), ''), v.nombre) as vehiculo,
			coalesce(cat.nombre, '') as categoria,
			coalesce(seg.nombre, '') as segmento,
			coalesce(concat_ws(' ', u.nombre, u.apellido), '') as vendedor,
			coalesce(concat_ws(' ', upr.nombre, upr.apellido), '') as usuario_pago_reserva
		from ventas_vehiculos vv
		inner join clientes c on c.id = vv.id_cliente
		inner join vehiculos v on v.id = vv.id_vehiculo
		left join categoria_vehiculo cat on cat.id = v.id_categoria
		left join segmento_vehiculo seg on seg.id = v.id_segmento
		left join usuarios u on u.id = vv.id_usuario
		left join usuarios upr on upr.id = vv.id_usuario_pago_reserva
		left join lateral (
			select
				json_agg(json_build_object('id', p.id, 'moneda', p.moneda, 'metodo', p.metodo, 'monto', p.monto) order by p.id) as pagos,
				string_agg(concat(p.moneda, ' ', p.metodo, ' ', p.monto::text), ' | ' order by p.id) as detalle_pago,
				count(*) as cantidad,
				max(p.metodo) as unico_metodo,
				coalesce(sum(case when p.moneda = 'USD' then p.monto else 0 end), 0) as total_usd,
				coalesce(sum(case when p.moneda = 'BOB' then p.monto else 0 end), 0) as total_bob
			from pagos_venta p
			where p.venta_id = vv.id
		) pv on true`
}
