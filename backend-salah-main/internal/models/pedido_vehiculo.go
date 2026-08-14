package models

import "time"

type PedidoVehiculo struct {
	ID                   uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDCliente            uint       `gorm:"column:id_cliente;not null;index" json:"id_cliente"`
	IDProveedor          *uint      `gorm:"column:id_proveedor;index" json:"id_proveedor"`
	IDVehiculo           *uint      `gorm:"column:id_vehiculo;index" json:"id_vehiculo"`
	IDUsuario            *uint      `gorm:"column:id_usuario;index" json:"id_usuario"`
	Fecha                time.Time  `gorm:"column:fecha;not null" json:"fecha"`
	Marca                string     `gorm:"column:marca;size:120;not null" json:"marca"`
	Modelo               string     `gorm:"column:modelo;size:120;not null" json:"modelo"`
	Anio                 uint       `gorm:"column:anio;not null" json:"anio"`
	Color                string     `gorm:"column:color;size:80" json:"color"`
	Version              string     `gorm:"column:version" json:"version"`
	PaisOrigen           string     `gorm:"column:pais_origen;size:120;not null" json:"pais_origen"`
	PrecioEstimadoUSD    float64    `gorm:"column:precio_estimado_usd;not null;default:0" json:"precio_estimado_usd"`
	TipoCambioUsado      float64    `gorm:"column:tipo_cambio_usado;not null;default:0" json:"tipo_cambio_usado"`
	FechaLlegadaEstimada time.Time  `gorm:"column:fecha_llegada_estimada;not null" json:"fecha_llegada_estimada"`
	AdelantoRequeridoUSD float64    `gorm:"column:adelanto_requerido_usd;not null;default:0" json:"adelanto_requerido_usd"`
	AdelantoPorcentaje   float64    `gorm:"column:adelanto_porcentaje;not null;default:0" json:"adelanto_porcentaje"`
	AdelantoPagadoUSD    float64    `gorm:"column:adelanto_pagado_usd;not null;default:0" json:"adelanto_pagado_usd"`
	AdelantoPagadoBOB    float64    `gorm:"column:adelanto_pagado_bob;not null;default:0" json:"adelanto_pagado_bob"`
	SaldoPendienteUSD    float64    `gorm:"column:saldo_pendiente_usd;not null;default:0" json:"saldo_pendiente_usd"`
	Estado               string     `gorm:"column:estado;size:40;not null;default:Pedido registrado" json:"estado"`
	FechaVencimiento     time.Time  `gorm:"column:fecha_vencimiento;not null" json:"fecha_vencimiento"`
	FechaRecepcion       *time.Time `gorm:"column:fecha_recepcion" json:"fecha_recepcion"`
	FechaCompletado      *time.Time `gorm:"column:fecha_completado" json:"fecha_completado"`
	Observacion          string     `gorm:"column:observacion" json:"observacion"`
	IDUsuarioRecibe      *uint      `gorm:"column:id_usuario_recibe" json:"id_usuario_recibe"`
	IDUsuarioCompleta    *uint      `gorm:"column:id_usuario_completa" json:"id_usuario_completa"`

	Cliente         Cliente       `gorm:"foreignKey:IDCliente;reference:ID" json:"-"`
	Proveedor       ProveedorAuto `gorm:"foreignKey:IDProveedor;reference:ID" json:"-"`
	Vehiculo        Vehiculo      `gorm:"foreignKey:IDVehiculo;reference:ID" json:"-"`
	Usuario         Usuario       `gorm:"foreignKey:IDUsuario;reference:ID" json:"-"`
	UsuarioRecibe   Usuario       `gorm:"foreignKey:IDUsuarioRecibe;reference:ID" json:"-"`
	UsuarioCompleta Usuario       `gorm:"foreignKey:IDUsuarioCompleta;reference:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (PedidoVehiculo) TableName() string {
	return "pedidos_vehiculos"
}

type PagoPedidoVehiculo struct {
	ID        uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PedidoID  uint           `gorm:"column:pedido_id;not null;index" json:"pedido_id"`
	Etapa     string         `gorm:"column:etapa;size:20;not null;default:adelanto" json:"etapa"`
	Moneda    string         `gorm:"column:moneda;size:3;not null" json:"moneda"`
	Metodo    string         `gorm:"column:metodo;size:40;not null" json:"metodo"`
	Monto     float64        `gorm:"column:monto;not null" json:"monto"`
	Pedido    PedidoVehiculo `gorm:"foreignKey:PedidoID;reference:ID" json:"-"`
	CreatedAt time.Time      `gorm:"default:now()" json:"created_at"`
}

func (PagoPedidoVehiculo) TableName() string {
	return "pagos_pedido_vehiculo"
}
