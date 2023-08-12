package database

import (
	"database/sql"
	"fmt"
	utils "github.com/TechSir3n/CityCompanion/assistance"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
	"path/filepath"
)

func init() {
	envFilePath, err := filepath.Abs(".env")
	if err != nil {
		utils.Fatal("Error[InitFunc]: %v", err)
	}

	if err := godotenv.Load(envFilePath); err != nil {
		utils.Fatal("Error[InitFunc]: %v", err)
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
		utils.Fatal("Failed open database %v", err)
	}

	if err = DB.Ping(); err != nil {
		utils.Error("Bad database connection: ", err.Error)
	}

	if _, err = DB.Exec(`CREATE TABLE IF NOT EXISTS UserLocation(id  SERIAL PRIMARY KEY,
		userID INTEGER NOT NULL ,latitude FLOAT,longitude FLOAT)`); err != nil {
		utils.Error("Failed to create table SaveLocation: %v ", err)
		return
	}

	if _, err = DB.Exec(`CREATE TABLE IF NOT EXISTS Radius(id SERIAL PRIMARY KEY,
		userID INTEGER NOT NULL,radius FLOAT)`); err != nil {
		utils.Error("Failed to create table SaveRadius: %v ", err)
		return
	}

	if _, err = DB.Exec(`CREATE TABLE IF NOT EXISTS StoragePlace(id SERIAL PRIMARY KEY,
		userID INTEGER NOT NULL,name VARCHAR(233),address VARCHAR(233))`); err != nil {
		utils.Error("Failed to create table SavePlace: %v ", err)
		return
	}

	if _, err := DB.Exec(`CREATE TABLE IF NOT EXISTS SavedFavoritePlace(id SERIAL PRIMARY KEY,
		userID INTEGER NOT NULL,name VARCHAR(233),address VARCHAR(233))`); err != nil {
		utils.Error("Failed to create table SaveFavoritePlace: %v ", err)
		return
	}

	if _, err := DB.Exec(`CREATE TABLE IF NOT EXISTS Reviews(id SERIAL PRIMARY KEY,name VARCHAR(233),address VARCHAR(233),userName VARCHAR(50),
	rating INTEGER,comment TEXT,created_at TIMESTAMP)`); err != nil {
		utils.Error("Failed to create table Reviews: %v", err)
		return
	}
}
