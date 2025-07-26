package main

import (
	"fmt"
	"net/http"

	"github.com/mmnsrti/rssproject/internal/auth"
	"github.com/mmnsrti/rssproject/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apicfg *apiConfig) authMiddleware(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Error getting API key: %v", err))
			return
		}
		user, err := apicfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error getting user by API key: %v", err))
			return
		}
		handler(w, r, user)
	}
}
