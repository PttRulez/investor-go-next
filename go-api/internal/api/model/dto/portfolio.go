package dto

type CreatePortfolio struct {
	Compound bool   `json:"compound,omitempty"`
	Name     string `json:"name"`
}

type UpdatePortfolio struct {
	Compound bool   `json:"compound,omitempty"`
	Id       int    `json:"id"`
	Name     string `json:"name,omitempty"`
}
