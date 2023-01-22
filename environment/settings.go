package environment

import "encoding/json"

type Database struct {
	Dialect         string `envconfig:"default=postgres"`
	UserName        string `envconfig:"default=postgres"`
	Password        string `envconfig:"default=pass"`
	Host            string `envconfig:"default=localhost"`
	Port            string `envconfig:"default=5432"`
	DatabaseName    string `envconfig:"default=pass"`
	SetMaxOpenConns int    `envconfig:"default=100"`
	LogMode         bool   `envconfig:"default=false"`
}

type settings struct {
	Database
}

var Settings settings

func (s settings) String() string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}
