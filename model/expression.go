package model

type Expression struct {
	ID         int    `gorm:"column:id" json:"id"`
	Definition string `gorm:"column:definition" json:"definition"`
}
