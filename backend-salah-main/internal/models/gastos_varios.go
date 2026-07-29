package models

import "time"

type GastoVario struct {
	ID           uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Nombre       string `gorm:"column:nombre;not null" json:"nombre"`
	UnidadMedida string `gorm:"column:unidad_medida" json:"unidad_medida"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (GastoVario) TableName() string {
	return "gastos_varios"
}
