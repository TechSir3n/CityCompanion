package database

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"path/filepath"
)

func init() {
	envFilePath, err := filepath.Abs(".env")
	if err != nil {
		log.Fatal("Error[InitFunc]:", err)
	}

	if err := godotenv.Load(envFilePath); err != nil {
		log.Fatal("Error[InitFunc]: ", err)
	}
}

var DB *sql.DB

func ConnectDB() {
	connStr := "postgres://postgres:BVEZR16AFLJMjQrMpVVp@containers-us-west-98.railway.app:5754/railway"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed open database ", err)
	}

	if err = DB.Ping(); err != nil {
		log.Printf("Bad database connection: %v", err)
	}

	if _, err = DB.Exec(`CREATE TABLE IF NOT EXISTS UserLocation(id  SERIAL PRIMARY KEY,
		userID INTEGER NOT NULL ,latitude FLOAT,longitude FLOAT)`); err != nil {
		log.Printf("Failed to create table SaveLocation: %v ", err)
		return
	}

	if _, err = DB.Exec(`CREATE TABLE IF NOT EXISTS Radius(id SERIAL PRIMARY KEY,
		userID INTEGER NOT NULL,radius FLOAT)`); err != nil {
		log.Printf("Failed to create table SaveRadius: %v ", err)
		return
	}

	if _, err = DB.Exec(`CREATE TABLE IF NOT EXISTS StoragePlace(id SERIAL PRIMARY KEY,
		userID INTEGER NOT NULL,name VARCHAR(233),address VARCHAR(233))`); err != nil {
		log.Printf("Failed to create table SavePlace: %v ", err)
		return
	}

	if _, err := DB.Exec(`CREATE TABLE IF NOT EXISTS SavedFavoritePlace(id SERIAL PRIMARY KEY,
		userID INTEGER NOT NULL,name VARCHAR(233),address VARCHAR(233))`); err != nil {
		log.Printf("Failed to create table SaveFavoritePlace: %v ", err)
		return
	}

	if _, err := DB.Exec(`CREATE TABLE IF NOT EXISTS Reviews(id SERIAL PRIMARY KEY,name VARCHAR(233),address VARCHAR(233),userName VARCHAR(50),
	rating INTEGER,comment TEXT,created_at TIMESTAMP)`); err != nil {
		log.Printf("Failed to create table Reviews: %v", err)
		return
	}

	fmt.Println("Success connected db")
}
