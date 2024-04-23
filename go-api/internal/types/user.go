package types

type User struct {
	Email          string `json:"email" db:"email"`
	HashedPassword string `json:"-" db:"hashed_password"`
	Id             int    `json:"id" db:"id"`
	Name           string `json:"name" db:"name"`
	Role           Role   `json:"role" db:"role"`
}

type RegisterUser struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     Role   `json:"role" `
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
