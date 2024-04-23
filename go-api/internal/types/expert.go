package types

type Expert struct {
	AvatarUrl string `json:"avatarUrl" db:"avatar_url"`
	Id        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	UserId    int    `json:"-" db:"user_id"`
}
