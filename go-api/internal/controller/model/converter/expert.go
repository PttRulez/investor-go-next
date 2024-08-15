package converter

import (
	"github.com/pttrulez/investor-go/internal/entity"
	"github.com/pttrulez/investor-go/pkg/api"
)

func FromCreateExpertRequestToExpert(req api.CreateExpertRequest) *entity.Expert {
	return &entity.Expert{
		AvatarURL: req.AvatarUrl,
		Name:      req.Name,
	}
}

func FromExpertToExpertResponse(e *entity.Expert) api.ExpertResponse {
	return api.ExpertResponse{
		AvatarUrl: e.AvatarURL,
		Id:        e.ID,
		Name:      e.Name,
	}
}
