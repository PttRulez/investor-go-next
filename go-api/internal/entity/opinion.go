package entity

type Opinion struct {
	Date         Date
	Exchange     Exchange
	ExpertID     int
	ID           int
	SecurityID   int
	SecurityType SecurityType
	SourceLink   *string
	TargetPrice  *float64
	Text         string
	Type         OpinionType
	UserID       int
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
