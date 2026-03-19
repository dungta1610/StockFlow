package model

type Paging struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

func NewPaging() *Paging {
	return &Paging{
		Page:  1,
		Limit: 10,
	}
}

func (p *Paging) Normalize() {
	if p == nil {
		return
	}

	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 {
		p.Limit = 10
	}

	if p.Limit > 100 {
		p.Limit = 100
	}
}

func (p *Paging) Offset() int {
	if p == nil {
		return 0
	}

	return (p.Page - 1) * p.Limit
}
