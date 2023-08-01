package database

import (
	"log"
	"path/filepath"
	"github.com/joho/godotenv"
)

type ListPreferences interface { 

};

func init() { 
	envFilePath, err := filepath.Abs(".env")
    if err != nil {
		log.Fatalf("Error[InitFunc]: %v",err)
    }
  
	if err:=godotenv.Load(envFilePath); err!=nil { 
			log.Fatalf("Error[InitFunc]: %v",err)
	}
}