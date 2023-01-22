package model

import "encoding/json"

type Expression struct {
	ID         int    `gorm:"column:id" json:"id"`
	Definition string `gorm:"column:definition" json:"definition"`
}

func (f Expression) String() string {
	bytes, _ := json.Marshal(f)
	return string(bytes)
}

func (Expression) TableName() string {
	return "expression"
}

type Response struct {
	Definition string `json:"definition"`
	Values     string `json:"values"`
	Result     bool   `json:"result"`
}
