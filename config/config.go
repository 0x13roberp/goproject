package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type configDataBase struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

var DataBase configDataBase

var SecretJWTKey string
var DateFormat string

func Init() {
	DataBase.Host = os.Getenv("DB_HOST")
	DataBase.Port = os.Getenv("DB_PORT")
	DataBase.User = os.Getenv("DB_USER")
	DataBase.Pass = os.Getenv("DB_PASS")
	DataBase.Name = os.Getenv("DB_NAME")
	SecretJWTKey = os.Getenv("JWT_SECRET")
	DateFormat = os.Getenv("DATE_FORMAT")
}
