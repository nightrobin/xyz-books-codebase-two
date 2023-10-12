package main

import (
	// "fmt"
	"os"
	"log"
	"path/filepath"

	// "xyz-books/router"
	
	"github.com/joho/godotenv"

)


func main() {
	// Load Environment Variables
	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)

	err = godotenv.Load(exPath + "/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}