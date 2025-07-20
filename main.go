package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)
	


func main() {
	fmt.Print("Hello, World!\n")
	godotenv.Load(".env")
	portString := os.Getenv("PORT") 
	if portString == "" {
		log.Fatal("Environment variable PORT is not set")
	}

	fmt.Println("Environment variable PORT is set to:", portString)
}