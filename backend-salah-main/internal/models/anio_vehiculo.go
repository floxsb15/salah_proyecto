package models

import "time"

type AnioVehiculo struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Valor     uint      `gorm:"column:valor;not null;uniqueIndex" json:"valor"`
	Estado    bool      `gorm:"column:estado" json:"estado"`
	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (AnioVehiculo) TableName() string {
	return "anios_vehiculos"
}
