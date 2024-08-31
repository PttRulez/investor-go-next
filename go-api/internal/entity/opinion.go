package entity

type Opinion struct {
	Date         Date         `json:"date"`
	Exchange     Exchange     `json:"exchange"`
	ExpertID     int          `json:"expert_id"`
	Expert       Expert       `json:"expert"`
	ID           int          `json:"id"`
	SecurityID   int          `json:"security_id"`
	SecurityType SecurityType `json:"security_type"`
	SourceLink   *string      `json:"sourceLink"`
	TargetPrice  *float64     `json:"targetPrice"`
	Text         string       `json:"text"`
	Ticker       string       `json:"ticker"`
	Type         OpinionType  `json:"type"`
	UserID       int          `json:"user_id"`
}

type OpinionType string

const (
	Flat      OpinionType = "FLAT"
	General   OpinionType = "GENERAL"
	Growth    OpinionType = "GROWTH"
	Reduction OpinionType = "REDUCTION"
)

func (e OpinionType) Validate() bool {
	switch e {
	case Flat:
	case General:
	case Growth:
	case Reduction:
	default:
		return false
	}
	return true
}

type OpinionFilters struct {
	ExpertID     *int
	SecurityID   *int
	Exchange     *Exchange
	SecurityType *SecurityType
}
