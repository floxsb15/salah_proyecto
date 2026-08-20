package models

import (
	"time"
)

type Vehiculo struct {
	ID                 uint    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Nombre             string  `gorm:"column:nombre;not null" json:"nombre"`
	Descripcion        string  `gorm:"column:descripcion" json:"descripcion"`
	Precio             float64 `gorm:"not null" json:"precio" format:"%.2f"`
	PrecioCompra       float64 `gorm:"column:precio_compra;not null;default:0" json:"-" format:"%.2f"`
	CantidadDisponible uint    `gorm:"column:cantidad_disponible;not null;default:1" json:"cantidad_disponible"`
	Imagen             string  `gorm:"column:imagen" json:"imagen"`
	Estado             bool    `gorm:"column:estado" json:"estado"`
	IDCategoria        uint    `gorm:"column:id_categoria" json:"id_categoria"`
	IDSegmento         *uint   `gorm:"column:id_segmento" json:"id_segmento,omitempty"`
	Marca              string  `gorm:"column:marca" json:"marca"`
	Modelo             string  `gorm:"column:modelo" json:"modelo"`
	Anio               uint    `gorm:"column:anio" json:"anio"`
	Version            string  `gorm:"column:version" json:"version"`
	TipoTecho          string  `gorm:"column:tipo_techo" json:"tipo_techo"`
	Combustible        string  `gorm:"column:combustible" json:"combustible"`
	Traccion           string  `gorm:"column:traccion" json:"traccion"`
	Transmision        string  `gorm:"column:transmision" json:"transmision"`
	Asientos           *uint   `gorm:"column:asientos" json:"asientos,omitempty"`
	Garantia           string  `gorm:"column:garantia" json:"garantia"`
	Equipamiento       string  `gorm:"column:equipamiento" json:"equipamiento"`
	PedidoOrigenID     *uint   `gorm:"column:pedido_origen_id;index" json:"pedido_origen_id,omitempty"`

	Categoria CategoriaVehiculo `gorm:"foreignKey:IDCategoria;reference:ID" json:"-"`
	Segmento  SegmentoVehiculo  `gorm:"foreignKey:IDSegmento;reference:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Vehiculo) TableName() string {
	return "vehiculos"
}
