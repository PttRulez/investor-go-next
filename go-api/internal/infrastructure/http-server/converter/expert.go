package converter

import (
	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server/contracts"
)

func FromCreateExpertRequestToExpert(req contracts.CreateExpertRequest) domain.Expert {
	return domain.Expert{
		AvatarURL: req.AvatarUrl,
		Name:      req.Name,
	}
}

func FromExpertToExpertResponse(e domain.Expert) contracts.ExpertResponse {
	return contracts.ExpertResponse{
		AvatarUrl: e.AvatarURL,
		Id:        e.ID,
		Name:      e.Name,
	}
}
