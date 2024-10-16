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

func GetCurrentUserTgChatID(ctx context.Context) string {
	_, claims, _ := jwtauth.FromContext(ctx)

	if chatId, ok := claims["tg_chat_id"].(string); ok {
		return chatId
	}

	return ""
}

func SignsAfterDot(f float64) int {
	s := fmt.Sprintf("%v", f)
	parts := strings.Split(s, ".")
	if len(parts) == 1 {
		return 0
	}
	decimalPart := parts[1]
	return len(decimalPart)
}
