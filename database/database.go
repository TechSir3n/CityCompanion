package database

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"path/filepath"
)

func init() {
	envFilePath, err := filepath.Abs(".env")
	if err != nil {
		log.Fatalf("Error[InitFunc]: %v", err)
	}

	if err := godotenv.Load(envFilePath); err != nil {
		log.Fatalf("Error[InitFunc]: %v", err)
	}
}

var DB *sql.DB

func ConnectDB() {
	db_host := os.Getenv("DB_HOST")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	db_user := os.Getenv("DB_USER")
	db_port := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db_host, db_port, db_user, db_password, db_name)
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("[ConnectDB]: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Bad connection to db: %v", err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS SaveLocation(id  SERIAL PRIMARY KEY,latitude FLOAT,longitude FLOAT)`)
	if err != nil {
		log.Fatalf("[ConnectDB]: %v ", err)
	}
}
