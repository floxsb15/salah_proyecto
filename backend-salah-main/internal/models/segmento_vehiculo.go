package models

import "time"

type SegmentoVehiculo struct {
	ID          uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Nombre      string `gorm:"column:nombre;not null" json:"nombre"`
	Descripcion string `gorm:"column:descripcion" json:"descripcion"`
	Estado      bool   `gorm:"column:estado" json:"estado"`
	IDCategoria uint   `gorm:"column:id_categoria;not null" json:"id_categoria"`

	Categoria CategoriaVehiculo `gorm:"foreignKey:IDCategoria;reference:ID" json:"-"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (SegmentoVehiculo) TableName() string {
	return "segmento_vehiculo"
}
