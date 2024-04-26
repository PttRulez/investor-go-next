package controller

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func getUserIdFromJwt(r *http.Request) int {
	_, claims, _ := jwtauth.FromContext(r.Context())
	return int(claims["id"].(float64))
}
