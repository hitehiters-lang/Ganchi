package connection

import (
	"database/sql"
	"fmt"
	"ganchi_app/additional"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var db *bun.DB
var (
	Host       string
	Port       string
	User       string
	Password   string
	Database   string
	ServerPort string
)

func InitDatabase() *bun.DB {
	if err := LoadENV(); err != nil {
		additional.PrintServerError(err.Error())
		return nil
	}
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(fmt.Sprintf(
			"postgres://%v:%v@%v:%v/%v?sslmode=disable",
			User, Password, Host, Port, Database))))

	db = bun.NewDB(sqldb, pgdialect.New())

	if err := db.Ping(); err != nil {
		additional.PrintServerError(err.Error())
		return nil
	}
	additional.PrintMessage("Успешно подключено к базе данных " + Database)
	return db
}

func GetDatabase() *bun.DB {
	return db
}
