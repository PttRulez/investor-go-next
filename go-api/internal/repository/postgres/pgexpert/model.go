package pgexpert

type Expert struct {
	AvatarUrl string `db:"avatar_url"`
	Id        int    `db:"id"`
	Name      string `db:"name"`
	UserId    int    `db:"user_id"`
}
