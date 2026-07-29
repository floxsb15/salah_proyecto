package models

import "time"

type Cliente struct {
	ID        uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Nombre    string `gorm:"column:nombre;size:255;not null" json:"nombre"`
	Apellido  string `gorm:"column:apellido;size:255" json:"apellido"`
	CI        string `gorm:"column:ci;size:30" json:"ci"`
	Celular   string `gorm:"column:celular;size:30" json:"celular"`
	Direccion string `gorm:"column:direccion;size:255" json:"direccion"`
	Estado    bool   `gorm:"column:estado" json:"estado"`

	CreatedAt time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (Cliente) TableName() string {
	return "clientes"
}
