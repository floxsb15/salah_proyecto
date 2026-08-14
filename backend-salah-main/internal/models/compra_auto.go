package models

import "time"

type ProveedorAuto struct {
	ID            uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Nombre        string    `gorm:"column:nombre;size:255;not null" json:"nombre"`
	CINIT         string    `gorm:"column:ci_nit;size:30" json:"ci_nit"`
	Telefono      string    `gorm:"column:telefono;size:30" json:"telefono"`
	Email         string    `gorm:"column:email;size:120" json:"email"`
	Direccion     string    `gorm:"column:direccion;size:255" json:"direccion"`
	Tipo          string    `gorm:"column:tipo;size:40" json:"tipo"`
	Observaciones string    `gorm:"column:observaciones" json:"observaciones"`
	Estado        bool      `gorm:"column:estado;not null;default:true" json:"estado"`
	CreatedAt     time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

func (ProveedorAuto) TableName() string {
	return "proveedores_autos"
}

type CompraAuto struct {
	ID                uint          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDVehiculo        uint          `gorm:"column:id_vehiculo;not null;index" json:"id_vehiculo"`
	IDProveedor       uint          `gorm:"column:id_proveedor;not null;index" json:"id_proveedor"`
	IDUsuario         *uint         `gorm:"column:id_usuario" json:"id_usuario"`
	FechaCompra       time.Time     `gorm:"column:fecha_compra;not null" json:"fecha_compra"`
	MonedaPrecio      string        `gorm:"column:moneda_precio;size:3;not null;default:USD" json:"moneda_precio"`
	PrecioCompraUSD   float64       `gorm:"column:precio_compra_usd;not null;default:0" json:"precio_compra_usd"`
	PrecioCompraBOB   float64       `gorm:"column:precio_compra_bob;not null;default:0" json:"precio_compra_bob"`
	TipoCambioUsado   float64       `gorm:"column:tipo_cambio_usado;not null;default:0" json:"tipo_cambio_usado"`
	GastoImportacion  float64       `gorm:"column:gasto_importacion;not null;default:0" json:"gasto_importacion"`
	GastoReparacion   float64       `gorm:"column:gasto_reparacion;not null;default:0" json:"gasto_reparacion"`
	GastoTransporte   float64       `gorm:"column:gasto_transporte;not null;default:0" json:"gasto_transporte"`
	GastoPapeleo      float64       `gorm:"column:gasto_papeleo;not null;default:0" json:"gasto_papeleo"`
	GastosAdicionales float64       `gorm:"column:gastos_adicionales;not null;default:0" json:"gastos_adicionales"`
	CostoTotalUSD     float64       `gorm:"column:costo_total_usd;not null;default:0" json:"costo_total_usd"`
	MetodoPago        string        `gorm:"column:metodo_pago;size:40;not null;default:Contado" json:"metodo_pago"`
	EstadoPago        string        `gorm:"column:estado_pago;size:40;not null;default:Pendiente" json:"estado_pago"`
	Observacion       string        `gorm:"column:observacion" json:"observacion"`
	Vehiculo          Vehiculo      `gorm:"foreignKey:IDVehiculo;reference:ID" json:"-"`
	Proveedor         ProveedorAuto `gorm:"foreignKey:IDProveedor;reference:ID" json:"-"`
	Usuario           Usuario       `gorm:"foreignKey:IDUsuario;reference:ID" json:"-"`
	CreatedAt         time.Time     `gorm:"default:now()" json:"-"`
	UpdatedAt         time.Time     `json:"-"`
}

func (CompraAuto) TableName() string {
	return "compras_autos"
}

type PagoCompraAuto struct {
	ID        uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CompraID  uint       `gorm:"column:compra_id;not null;index" json:"compra_id"`
	Moneda    string     `gorm:"column:moneda;size:3;not null" json:"moneda"`
	Metodo    string     `gorm:"column:metodo;size:40;not null" json:"metodo"`
	Monto     float64    `gorm:"column:monto;not null" json:"monto"`
	Compra    CompraAuto `gorm:"foreignKey:CompraID;reference:ID" json:"-"`
	CreatedAt time.Time  `gorm:"default:now()" json:"created_at"`
}

func (PagoCompraAuto) TableName() string {
	return "pagos_compra_auto"
}
