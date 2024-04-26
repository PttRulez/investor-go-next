package user

type User struct {
	Id             int    `db:"id"`
	Email          string `db:"email"`
	HashedPassword string `db:"hashed_password"`
	Name           string `db:"name"`
	Role           string `db:"role"`
}
