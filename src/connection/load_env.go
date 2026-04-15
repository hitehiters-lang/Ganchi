package connection

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadENV() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	Host = os.Getenv("HOST")
	Port = os.Getenv("DB_PORT")
	User = os.Getenv("USER")
	Password = os.Getenv("PASSWORD")
	Database = os.Getenv("DB_NAME")
	ServerPort = os.Getenv("SERVER_PORT")
	return err
}
