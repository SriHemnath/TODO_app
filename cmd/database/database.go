package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
)

func NewDB() *sqlx.DB {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	usr, ok := viper.Get("userID").(string)
	if !ok {
		usr = os.Getenv("userID")
	}
	pwd, ok := viper.Get("pwd").(string)
	if !ok {
		pwd = os.Getenv("pwd")
	}
	connectionString := fmt.Sprintf("postgres://%s:%s@localhost:5432/todo?sslmode=disable", usr, pwd)
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
