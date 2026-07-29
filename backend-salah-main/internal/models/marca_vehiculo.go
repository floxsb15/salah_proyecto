package models

import "time"

type MarcaVehiculo struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Nombre    string    `gorm:"column:nombre;not null;uniqueIndex" json:"nombre"`
	Estado    bool      `gorm:"column:estado" json:"estado"`
	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (MarcaVehiculo) TableName() string {
	return "marcas_vehiculos"
}
