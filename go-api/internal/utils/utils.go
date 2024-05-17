package utils

import (
	"context"
	"github.com/go-chi/jwtauth/v5"
)

func GetCurrentUserId(ctx context.Context) int {
	_, claims, _ := jwtauth.FromContext(ctx)
	return int(claims["id"].(float64))
}
