package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var connectionString = "postgres://postgres:root@localhost:5432/todo?sslmode=disable"

func NewDB() *sqlx.DB {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if err := db.Ping(); err != nil {
		fmt.Println(err)
		return nil
	}
	return db
}
