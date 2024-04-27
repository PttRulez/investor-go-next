package pgportfolio

type Portfolio struct {
	Compound bool   `db:"compound"`
	Id       int    `db:"id"`
	Name     string `db:"name"`
	UserId   int    `db:"user_id"`
}
