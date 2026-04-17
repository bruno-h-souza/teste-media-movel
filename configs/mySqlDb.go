package configs

import (
	"database/sql"
	"fmt"
	"teste-media-movel/internal/utils"

	_ "github.com/go-sql-driver/mysql" // Driver do MySQL
)

// var DbConnMySql *sql.DB

func NewMySqlConnection() (*sql.DB, error) {
	dbUser := utils.GetEnv("DB_USER", "root")
	dbPass := utils.GetEnv("DB_PASSWORD", "secret")
	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "3306")
	dbName := utils.GetEnv("DB_NAME", "usersdb")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
