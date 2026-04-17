package configs

import (
	"database/sql"
	"fmt"
	"teste-media-movel/internal/utils"
)

var DbConnMySql *sql.DB

func NewMySqlConnection() error {
	dbUser := utils.GetEnv("DB_USER", "root")
	dbPass := utils.GetEnv("DB_PASSWORD", "secret")
	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "3306")
	dbName := utils.GetEnv("DB_NAME", "usersdb")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}

	DbConnMySql = db

	return nil
}
