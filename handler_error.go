package main
import "net/http"

func handlerError (w http.ResponseWriter, r *http.Request) {
	// Simulate an error for demonstration purposes
	errMessage := "An unexpected error occurred"
	respondWithError(w, 400, errMessage)
}