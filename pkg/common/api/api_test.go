package api

import (
	"encoding/json"
	"fmt"
	"mindgpt/pkg/common/query"
	"testing"
)

func TestSuccess(t *testing.T) {
	s := SuccessPagination(12345, &query.Pagination{
		Page:      1,
		Limit:     10,
		Total:     11,
		TotalPage: 2,
		Sort:      "+id",
	})
	v, _ := json.MarshalIndent(s, "", "  ")
	fmt.Println(string(v))
}
