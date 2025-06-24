package main

import (
	"fmt"
	_ "net"
	_"time"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=127.0.0.1 port=5432 user=user dbname=mydb password=pswd sslmode=disable",
		))
	if err != nil {
		fmt.Print(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Print(err)
	}

}