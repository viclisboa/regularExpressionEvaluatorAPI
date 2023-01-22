package model

type Expression struct {
	ID         int    `gorm:"column:id" json:"id"`
	Definition string `gorm:"column:definition" json:"definition"`
}

type Response struct {
	Definition string `json:"definition"`
	Values     string `json:"values"`
	Result     bool   `json:"result"`
}
