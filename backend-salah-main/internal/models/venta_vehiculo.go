package models

import "time"

type VentaVehiculo struct {
	ID                       uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDCliente                uint       `gorm:"column:id_cliente;not null" json:"id_cliente"`
	IDVehiculo               *uint      `gorm:"column:id_vehiculo" json:"id_vehiculo"`
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
	Observacion              string     `gorm:"column:observacion" json:"observacion"`
	IDUsuarioPagoReserva     *uint      `gorm:"column:id_usuario_pago_reserva" json:"id_usuario_pago_reserva"`
	TipoReserva              string     `gorm:"column:tipo_reserva;size:20;not null;default:stock" json:"tipo_reserva"`
	PedidoMarca              string     `gorm:"column:pedido_marca;size:120" json:"pedido_marca"`
	PedidoModelo             string     `gorm:"column:pedido_modelo;size:120" json:"pedido_modelo"`
	PedidoAnio               uint       `gorm:"column:pedido_anio;not null;default:0" json:"pedido_anio"`
	PedidoColor              string     `gorm:"column:pedido_color;size:80" json:"pedido_color"`
	PedidoVersion            string     `gorm:"column:pedido_version" json:"pedido_version"`
	PedidoPaisOrigen         string     `gorm:"column:pedido_pais_origen;size:120" json:"pedido_pais_origen"`
	PedidoProveedor          string     `gorm:"column:pedido_proveedor;size:180" json:"pedido_proveedor"`
	PedidoLlegadaEstimada    string     `gorm:"column:pedido_llegada_estimada;size:180" json:"pedido_llegada_estimada"`

	Cliente            Cliente  `gorm:"foreignKey:IDCliente;reference:ID" json:"-"`
	Vehiculo           Vehiculo `gorm:"foreignKey:IDVehiculo;reference:ID" json:"-"`
	Usuario            Usuario  `gorm:"foreignKey:IDUsuario;reference:ID" json:"-"`
	UsuarioPagoReserva Usuario  `gorm:"foreignKey:IDUsuarioPagoReserva;reference:ID" json:"-"`

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

type CuotaCredito struct {
	ID               uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDVentaVehiculo  uint       `gorm:"column:id_venta_vehiculo;not null;index" json:"id_venta_vehiculo"`
	Numero           uint       `gorm:"column:numero;not null" json:"numero"`
	FechaVencimiento time.Time  `gorm:"column:fecha_vencimiento;not null" json:"fecha_vencimiento"`
	Monto            float64    `gorm:"column:monto;not null" json:"monto"`
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
