package main

import (
	"fmt"
	"net/http"

	"github.com/keithfy96/go-project/auth"
	"github.com/keithfy96/go-project/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (c *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondwithError(w, 403, fmt.Sprintf("Error getting API Key: %v", err))
			return
		}

		user, err := c.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondwithError(w, 403, fmt.Sprintf("Error getting user by API Key: %v", err))
			return
		}

		handler(w, r, user)
	}
}
