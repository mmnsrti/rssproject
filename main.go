package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mmnsrti/rssproject/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Environment variable PORT is not set")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Environment variable DB_URL is not set")
	}
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	queries := database.New(conn)

	apicfg := apiConfig{
		DB: queries,
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
	v1Router.Get("/error", handlerError)
	v1Router.Post("/users", apicfg.handlerCreateUser)
	v1Router.Get("/users", apicfg.authMiddleware(apicfg.handlerGetUserByAPIKey))
	v1Router.Post("/feeds", apicfg.authMiddleware(apicfg.handlerCreateFeed))
	v1Router.Get("/feeds", apicfg.handlerGetFeeds)
	v1Router.Post("/feed_follows", apicfg.authMiddleware(apicfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apicfg.authMiddleware(apicfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apicfg.authMiddleware(apicfg.handlerFeedFollowDelete))
	router.Mount("/v1", v1Router)
	log.Printf("Starting server on port %s", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

	fmt.Println("Environment variable PORT is set to:", portString)
}
