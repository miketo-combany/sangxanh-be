package query

import (
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

const (
	DefaultPage  = int64(1)
	DefaultLimit = int64(20)
	MaxLimit     = int64(100)
)

type Pagination struct {
	Page      int64  `json:"page" query:"page"`
	Limit     int64  `json:"limit" query:"limit"`
	Total     int64  `json:"total"`
	TotalPage int64  `json:"total_page"`
	Sort      string `json:"sort" query:"sort"`
}

func (p *Pagination) Correct() {
	p.Page = max(p.Page, DefaultPage)
	p.Limit = max(p.Limit, DefaultLimit)
	p.Limit = min(p.Limit, MaxLimit)
}

func (p *Pagination) SetTotal(total int64) {
	p.Total = total
	p.TotalPage = total / p.Limit
	if total%p.Limit > 0 {
		p.TotalPage++
	}
}

func (p *Pagination) GetSort() bson.M {
	if len(p.Sort) == 0 {
		return nil
	}
	sort := make(bson.M)
	parts := strings.Split(p.Sort, ",")
	for _, part := range parts {
		if len(part) == 1 {
			continue
		}
		if strings.HasPrefix(part, "-") {
			sort[strings.TrimPrefix(part, "-")] = -1
		}
		if strings.HasPrefix(part, "+") {
			sort[strings.TrimPrefix(part, "+")] = 1
		}
	}
	return sort
}
