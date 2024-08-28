package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-chi/jwtauth/v5"
)

func GetCurrentUserID(ctx context.Context) int {
	_, claims, _ := jwtauth.FromContext(ctx)
	return int(claims["id"].(float64))
}

func SignsAfterDot(f float64) int {
	s := fmt.Sprintf("%v", f)
	decimalPart := strings.Split(s, ".")[1]
	return len(decimalPart)
}
