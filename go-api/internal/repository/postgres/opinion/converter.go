package opinion

import "github.com/pttrulez/investor-go/internal/model"

func FromDBToModelOpinion(opinion *Opinion) *model.Opinion {
	return &model.Opinion{
		Date:         opinion.Date,
		Exchange:     opinion.Exchange,
		ExpertId:     opinion.ExpertId,
		Id:           opinion.Id,
		SecurityId:   opinion.SecurityId,
		SecurityType: opinion.SecurityType,
		SourceLink:   opinion.SourceLink,
		TargetPrice:  opinion.TargetPrice,
		Type:         opinion.Type,
		UserId:       opinion.UserId,
	}
}

func FromOpinionToDBOpinion(opinion *model.Opinion) *Opinion {
	return &Opinion{
		Date:         opinion.Date,
		Exchange:     opinion.Exchange,
		ExpertId:     opinion.ExpertId,
		Id:           opinion.Id,
		SecurityId:   opinion.SecurityId,
		SecurityType: opinion.SecurityType,
		SourceLink:   opinion.SourceLink,
		TargetPrice:  opinion.TargetPrice,
		Type:         opinion.Type,
		UserId:       opinion.UserId,
	}
}
