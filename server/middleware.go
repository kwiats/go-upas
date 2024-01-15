package server

import (
	"fmt"
	"net/http"

	"github.com/kwiats/go-upas/models"
	"github.com/kwiats/go-upas/utils"
)

func withJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling JWT auth middleware")
		tokenString := r.Header.Get("x-jwt-token")
		_, err := utils.ValidateJWT(tokenString)
		if err != nil {
			if err := models.WriteJSONResponse(w, http.StatusForbidden, models.APIError{Error: "invalid token"}); err != nil {
				return
			}
			return
		}
		handlerFunc(w, r)
	}
}
