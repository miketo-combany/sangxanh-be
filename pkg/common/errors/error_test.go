package errors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBadRequest(t *testing.T) {
	e := BadRequest("error %v", 1).WithDebug(map[string]any{"key": "value"})
	assert.Equal(t, "error 1", e.Error())
}
