package models

import "time"

type ProformaVehicular struct {
	ID                uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDCliente         *uint     `gorm:"column:id_cliente" json:"id_cliente"`
	IDVehiculo        uint      `gorm:"column:id_vehiculo;not null" json:"id_vehiculo"`
	IDUsuario         *uint     `gorm:"column:id_usuario" json:"id_usuario"`
	Fecha             time.Time `gorm:"column:fecha;not null" json:"fecha"`
	ClienteNombre     string    `gorm:"column:cliente_nombre;size:255;not null" json:"cliente_nombre"`
	ClienteDireccion  string    `gorm:"column:cliente_direccion;size:255" json:"cliente_direccion"`
	ClienteTelefono   string    `gorm:"column:cliente_telefono;size:30" json:"cliente_telefono"`
	Modalidad         string    `gorm:"column:modalidad;size:50;not null" json:"modalidad"`
	PrecioUnidad      float64   `gorm:"column:precio_unidad;not null" json:"precio_unidad"`
	Cantidad          uint      `gorm:"column:cantidad;not null;default:1" json:"cantidad"`
	PrecioTotal       float64   `gorm:"column:precio_total;not null" json:"precio_total"`
	CuotaInicial      float64   `gorm:"column:cuota_inicial;not null;default:0" json:"cuota_inicial"`
	Saldo             float64   `gorm:"column:saldo;not null;default:0" json:"saldo"`
	ValidezDias       uint      `gorm:"column:validez_dias;not null;default:10" json:"validez_dias"`
	FechaVencimiento  time.Time `gorm:"column:fecha_vencimiento;not null" json:"fecha_vencimiento"`
	PrecioCatalogoRef float64   `gorm:"column:precio_catalogo_ref;not null;default:0" json:"precio_catalogo_ref"`
	CreatedAt         time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt         time.Time `json:"-"`

	Cliente  Cliente  `gorm:"foreignKey:IDCliente;references:ID" json:"-"`
	Vehiculo Vehiculo `gorm:"foreignKey:IDVehiculo;references:ID" json:"-"`
	Usuario  Usuario  `gorm:"foreignKey:IDUsuario;references:ID" json:"-"`
}

func (ProformaVehicular) TableName() string {
	return "proformas_vehiculares"
}
