package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Environment variable PORT is not set")
	}

	router := chi.NewRouter()
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		ExposedHeaders:   []string{"Link"},
	}))
	log.Printf("Starting server on port %s", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

	fmt.Println("Environment variable PORT is set to:", portString)
}
