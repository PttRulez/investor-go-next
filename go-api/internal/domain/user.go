package domain

type Role string

const (
	Admin    Role = "ADMIN"
	Investor Role = "INVESTOR"
)

type User struct {
	Email          string
	HashedPassword string
	ID             int
	Name           string
	Password       string
	Role           Role
	TgChatID       *string
}
