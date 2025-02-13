package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPagination_Init(t *testing.T) {
	p := Pagination{
		Page:  -1,
		Limit: -1,
	}
	p.Correct()
	assert.Equal(t, DefaultPage, p.Page)
	assert.Equal(t, DefaultLimit, p.Limit)
	p.Limit = MaxLimit + 1
	p.Correct()
	assert.Equal(t, MaxLimit, p.Limit)
}

func TestPagination_SetTotal(t *testing.T) {
	p := Pagination{
		Page:  1,
		Limit: 10,
	}
	p.SetTotal(100)
	assert.Equal(t, int64(100), p.Total)
	assert.Equal(t, int64(10), p.TotalPage)
	p.SetTotal(101)
	assert.Equal(t, int64(101), p.Total)
	assert.Equal(t, int64(11), p.TotalPage)
}

func TestPagination_GetSort(t *testing.T) {
	p := Pagination{
		Sort: "-name,+age,",
	}
	sort := p.GetSort()
	assert.Equal(t, 2, len(sort))
	assert.Equal(t, -1, sort["name"])
	assert.Equal(t, 1, sort["age"])
}
