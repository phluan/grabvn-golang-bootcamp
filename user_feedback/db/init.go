package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var DBCon *sql.DB

func InitDB() {
	host := "localhost"
	port := 5432
	user := "luanpham"
	password := "123"
	dbname := "grab_week_3"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", // TODO: take care of SSL
		host, port, user, password, dbname)

	var err error
	DBCon, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = DBCon.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to Postgres")
}
