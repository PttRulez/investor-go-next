package domain

type Expert struct {
	AvatarURL *string `json:"avatarUrl"`
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	UserID    int     `json:"userId"`
}
