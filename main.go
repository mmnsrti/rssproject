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
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		ExposedHeaders:   []string{"Link"},
	}))
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	router.Mount("/v1", v1Router)
	v1Router.Get("/error", handlerError)
	log.Printf("Starting server on port %s", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

	fmt.Println("Environment variable PORT is set to:", portString)
}
