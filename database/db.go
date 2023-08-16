package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func ConnectDBRailway() {
	connStr := fmt.Sprintf("postgres://postgres:BVEZR16AFLJMjQrMpVVp@containers-us-west-98.railway.app:5754/railway")

	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed open database ", err)
	}

	if err = DB.Ping(); err != nil {
		log.Printf("Bad database connection: %v", err)
	}
}
