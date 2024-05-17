package entity

type Role string

const (
	Admin    Role = "ADMIN"
	Investor Role = "INVESTOR"
)

type User struct {
	Email          string
	HashedPassword string
	Id             int
	Name           string
	Role           Role

	Password string
}
