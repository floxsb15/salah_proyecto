package models

import "time"

type CategoriaVehiculo struct {
	ID          uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Nombre      string `gorm:"column:nombre;not null" json:"nombre"`
	Descripcion string `gorm:"column:descripcion" json:"descripcion"`
	Estado      bool   `gorm:"column:estado" json:"estado"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (CategoriaVehiculo) TableName() string {
	return "categoria_vehiculo"
}
