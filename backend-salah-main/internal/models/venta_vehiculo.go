package models

import "time"

type VentaVehiculo struct {
	ID                       uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDCliente                uint       `gorm:"column:id_cliente;not null" json:"id_cliente"`
	IDVehiculo               uint       `gorm:"column:id_vehiculo;not null" json:"id_vehiculo"`
	IDUsuario                *uint      `gorm:"column:id_usuario" json:"id_usuario"`
	Fecha                    time.Time  `gorm:"column:fecha;not null" json:"fecha"`
	FechaVenta               time.Time  `gorm:"column:fecha_venta;not null;default:now()" json:"fecha_venta"`
	TipoVenta                string     `gorm:"column:tipo_venta;size:40;not null" json:"tipo_venta"`
	Cantidad                 uint       `gorm:"column:cantidad;not null" json:"cantidad"`
	PrecioUnidad             float64    `gorm:"column:precio_unidad;not null" json:"precio_unidad"`
	PrecioTotal              float64    `gorm:"column:precio_total;not null" json:"precio_total"`
	PrecioUSD                float64    `gorm:"column:precio_usd;not null;default:0" json:"precio_usd"`
	TipoCambioUsado          float64    `gorm:"column:tipo_cambio_usado;not null;default:0" json:"tipo_cambio_usado"`
	MontoBOBCalculado        float64    `gorm:"column:monto_bob_calculado;not null;default:0" json:"monto_bob_calculado"`
	PagoUSD                  float64    `gorm:"column:pago_usd;not null;default:0" json:"pago_usd"`
	PagoBOB                  float64    `gorm:"column:pago_bob;not null;default:0" json:"pago_bob"`
	SaldoBOB                 float64    `gorm:"column:saldo_bob;not null;default:0" json:"saldo_bob"`
	CuotaInicial             float64    `gorm:"column:cuota_inicial;not null;default:0" json:"cuota_inicial"`
	Saldo                    float64    `gorm:"column:saldo;not null;default:0" json:"saldo"`
	MontoFinanciado          float64    `gorm:"column:monto_financiado;not null;default:0" json:"monto_financiado"`
	NumeroCuotas             uint       `gorm:"column:numero_cuotas;not null;default:0" json:"numero_cuotas"`
	MontoCuota               float64    `gorm:"column:monto_cuota;not null;default:0" json:"monto_cuota"`
	FechaInicioCredito       *time.Time `gorm:"column:fecha_inicio_credito" json:"fecha_inicio_credito"`
	FrecuenciaPago           string     `gorm:"column:frecuencia_pago;size:30" json:"frecuencia_pago"`
	TieneRespaldo            bool       `gorm:"column:tiene_respaldo;not null;default:false" json:"tiene_respaldo"`
	TipoGarantia             string     `gorm:"column:tipo_garantia" json:"tipo_garantia"`
	DocumentoGarantia        string     `gorm:"column:documento_garantia" json:"documento_garantia"`
	DatosGarante             string     `gorm:"column:datos_garante" json:"datos_garante"`
	ValidezProformaDias      uint       `gorm:"column:validez_proforma_dias;not null;default:15" json:"validez_proforma_dias"`
	FechaVencimientoProforma time.Time  `gorm:"column:fecha_vencimiento_proforma;not null" json:"fecha_vencimiento_proforma"`
	EstadoVenta              string     `gorm:"column:estado_venta;size:40;not null" json:"estado_venta"`
	EstadoPago               string     `gorm:"column:estado_pago;size:60;not null" json:"estado_pago"`
	MetodoPago               string     `gorm:"column:metodo_pago;size:40;not null;default:Efectivo" json:"metodo_pago"`
	EstadoEntrega            string     `gorm:"column:estado_entrega;size:40;not null" json:"estado_entrega"`
	FechaEntrega             *time.Time `gorm:"column:fecha_entrega" json:"fecha_entrega"`
	ReferenciaBancaria       string     `gorm:"column:referencia_bancaria" json:"referencia_bancaria"`
	EstadoDesembolso         string     `gorm:"column:estado_desembolso;size:60" json:"estado_desembolso"`
	BancoEntidadFinanciera   string     `gorm:"column:banco_entidad_financiera;size:120" json:"banco_entidad_financiera"`
	EstadoTramiteBancario    string     `gorm:"column:estado_tramite_bancario;size:40" json:"estado_tramite_bancario"`
	MontoFinanciarBanco      float64    `gorm:"column:monto_financiar_banco;not null;default:0" json:"monto_financiar_banco"`
	FechaEstimadaDesembolso  *time.Time `gorm:"column:fecha_estimada_desembolso" json:"fecha_estimada_desembolso"`
	NumeroOperacionBanco     string     `gorm:"column:numero_operacion_banco;size:80" json:"numero_operacion_banco"`
	FechaDesembolsoBanco     *time.Time `gorm:"column:fecha_desembolso_banco" json:"fecha_desembolso_banco"`
	MontoDesembolsadoBanco   float64    `gorm:"column:monto_desembolsado_banco;not null;default:0" json:"monto_desembolsado_banco"`
	IDUsuarioDesembolsoBanco *uint      `gorm:"column:id_usuario_desembolso_banco" json:"id_usuario_desembolso_banco"`
	Observacion              string     `gorm:"column:observacion" json:"observacion"`
	IDUsuarioPagoReserva     *uint      `gorm:"column:id_usuario_pago_reserva" json:"id_usuario_pago_reserva"`

	Cliente            Cliente  `gorm:"foreignKey:IDCliente;reference:ID" json:"-"`
	Vehiculo           Vehiculo `gorm:"foreignKey:IDVehiculo;reference:ID" json:"-"`
	Usuario            Usuario  `gorm:"foreignKey:IDUsuario;reference:ID" json:"-"`
	UsuarioPagoReserva Usuario  `gorm:"foreignKey:IDUsuarioPagoReserva;reference:ID" json:"-"`
	UsuarioDesembolso  Usuario  `gorm:"foreignKey:IDUsuarioDesembolsoBanco;reference:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (VentaVehiculo) TableName() string {
	return "ventas_vehiculos"
}

type PagoVenta struct {
	ID        uint          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	VentaID   uint          `gorm:"column:venta_id;not null;index" json:"venta_id"`
	Moneda    string        `gorm:"column:moneda;size:3;not null" json:"moneda"`
	Metodo    string        `gorm:"column:metodo;size:40;not null" json:"metodo"`
	Monto     float64       `gorm:"column:monto;not null" json:"monto"`
	Venta     VentaVehiculo `gorm:"foreignKey:VentaID;reference:ID" json:"-"`
	CreatedAt time.Time     `gorm:"default:now()" json:"created_at"`
}

func (PagoVenta) TableName() string {
	return "pagos_venta"
}

type PagoEngancheBancario struct {
	ID              uint          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDVentaVehiculo uint          `gorm:"column:id_venta_vehiculo;not null;index" json:"id_venta_vehiculo"`
	Moneda          string        `gorm:"column:moneda;size:3;not null" json:"moneda"`
	Metodo          string        `gorm:"column:metodo;size:40;not null" json:"metodo"`
	Monto           float64       `gorm:"column:monto;not null" json:"monto"`
	TipoCambio      float64       `gorm:"column:tipo_cambio;not null" json:"tipo_cambio"`
	EquivalenteUSD  float64       `gorm:"column:equivalente_usd;not null" json:"equivalente_usd"`
	VentaVehiculo   VentaVehiculo `gorm:"foreignKey:IDVentaVehiculo;reference:ID" json:"-"`
	CreatedAt       time.Time     `gorm:"default:now()" json:"created_at"`
}

func (PagoEngancheBancario) TableName() string {
	return "pagos_enganche_bancario"
}

type CuotaCredito struct {
	ID               uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDVentaVehiculo  uint       `gorm:"column:id_venta_vehiculo;not null;index" json:"id_venta_vehiculo"`
	Numero           uint       `gorm:"column:numero;not null" json:"numero"`
	FechaVencimiento time.Time  `gorm:"column:fecha_vencimiento;not null" json:"fecha_vencimiento"`
	Monto            float64    `gorm:"column:monto;not null" json:"monto"`
	MontoPagado      float64    `gorm:"column:monto_pagado;not null;default:0" json:"monto_pagado"`
	SaldoPendiente   float64    `gorm:"column:saldo_pendiente;not null;default:0" json:"saldo_pendiente"`
	TipoCambioPago   float64    `gorm:"column:tipo_cambio_pago;not null;default:0" json:"tipo_cambio_pago"`
	MontoBOBPagado   float64    `gorm:"column:monto_bob_pagado;not null;default:0" json:"monto_bob_pagado"`
	Estado           string     `gorm:"column:estado;size:20;not null;default:pendiente" json:"estado"`
	FechaPago        *time.Time `gorm:"column:fecha_pago" json:"fecha_pago"`
	IDUsuarioPago    *uint      `gorm:"column:id_usuario_pago" json:"id_usuario_pago"`

	VentaVehiculo VentaVehiculo `gorm:"foreignKey:IDVentaVehiculo;reference:ID" json:"-"`
	UsuarioPago   Usuario       `gorm:"foreignKey:IDUsuarioPago;reference:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (CuotaCredito) TableName() string {
	return "cuotas_credito"
}

type PagoCuotaCredito struct {
	ID              uint          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDCuotaCredito  uint          `gorm:"column:id_cuota_credito;not null;index" json:"id_cuota_credito"`
	IDVentaVehiculo uint          `gorm:"column:id_venta_vehiculo;not null;index" json:"id_venta_vehiculo"`
	MontoUSD        float64       `gorm:"column:monto_usd;not null" json:"monto_usd"`
	TipoCambio      float64       `gorm:"column:tipo_cambio;not null" json:"tipo_cambio"`
	MontoBOB        float64       `gorm:"column:monto_bob;not null" json:"monto_bob"`
	FechaPago       time.Time     `gorm:"column:fecha_pago;not null;default:now()" json:"fecha_pago"`
	IDUsuarioPago   uint          `gorm:"column:id_usuario_pago;not null" json:"id_usuario_pago"`
	Observacion     string        `gorm:"column:observacion" json:"observacion"`
	CuotaCredito    CuotaCredito  `gorm:"foreignKey:IDCuotaCredito;reference:ID" json:"-"`
	VentaVehiculo   VentaVehiculo `gorm:"foreignKey:IDVentaVehiculo;reference:ID" json:"-"`
	UsuarioPago     Usuario       `gorm:"foreignKey:IDUsuarioPago;reference:ID" json:"-"`
	CreatedAt       time.Time     `gorm:"default:now()" json:"created_at"`
}

func (PagoCuotaCredito) TableName() string {
	return "pagos_cuotas_credito"
}
