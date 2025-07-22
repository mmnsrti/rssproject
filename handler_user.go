package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mmnsrti/rssproject/internal/database"
)

func (apicfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `name`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	apicfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: uuid.New(),
	})

	respondWithJSON(w, 200, struct{}{})
}
