package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)
var table = `CREATE TABLE IF NOT EXISTS password_storage(
userid integer NOT NULL,
url varchar(255) NOT NULL,
login varchar(255) NOT NULL,
password varchar(255) NOT NULL
);`
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s  port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	_, err = db.Exec(table)
	if err != nil{
		return nil,  err
	}
	_, err = db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto")
	if err != nil {
		return nil, err
	}

	return db, nil
}
