package entity

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
	Role           Role

	Password string
}
