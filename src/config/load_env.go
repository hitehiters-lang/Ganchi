package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	Host     string
	Port     string
	User     string
	Password string
	Database string
)

const (
	Rub float64 = 1
	Dol float64 = 0.013
	Ten float64 = 6.2
	Yua float64 = 0.09
)

func LoadENV() (error, string, string, string, string, string) {
	err := godotenv.Load(".env")
	if err != nil {
		return err, "", "", "", "", ""
	}
	Host = os.Getenv("HOST")
	Port = os.Getenv("DB_PORT")
	User = os.Getenv("USER")
	Password = os.Getenv("PASSWORD")
	Database = os.Getenv("DB_NAME")
	return err, Host, Port, User, Password, Database
}
func GetServerPort() string {
	return ":" + os.Getenv("SERVER_PORT")
}
func GetEmailData() (email_from, email_to, pass string) {
	email_from = os.Getenv("CONTACT_MAIL")
	email_to = os.Getenv("MAIL_TO")
	pass = os.Getenv("MAIL_PASSWORD")
	return email_from, email_to, pass
}
