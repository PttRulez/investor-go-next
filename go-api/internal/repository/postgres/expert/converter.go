package expert

import "github.com/pttrulez/investor-go/internal/model"

func FromModelToDBExpert(e *model.Expert) *Expert {
	return &Expert{
		AvatarUrl: e.AvatarUrl,
		Id:        e.Id,
		Name:      e.Name,
		UserId:    e.UserId,
	}
}

func FromDBToModelExpert(e *Expert) *model.Expert {
	return &model.Expert{
		AvatarUrl: e.AvatarUrl,
		Id:        e.Id,
		Name:      e.Name,
		UserId:    e.UserId,
	}
}
