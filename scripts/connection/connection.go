package connection

import _ "github.com/go-sql-driver/mysql"

import (
	"database/sql"
	"log"
	"os"
)

func Connect() *sql.DB {

	MYSQL_USERNAME := os.Getenv("MYSQL_USERNAME")
	MYSQL_PASSWORD := os.Getenv("MYSQL_PASSWORD")
	MYSQL_HOST := os.Getenv("MYSQL_HOST")
	MYSQL_PORT := os.Getenv("MYSQL_PORT")
	MYSQL_DB := os.Getenv("MYSQL_DB")

	db, err := sql.Open("mysql", MYSQL_USERNAME+":"+MYSQL_PASSWORD+":"+"@tcp("+MYSQL_HOST+":"+MYSQL_PORT+")/"+MYSQL_DB)

	if err != nil {
		log.Fatal(err)
	}

	return db

}
