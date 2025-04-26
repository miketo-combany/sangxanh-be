package dto

import (
	"encoding/json"
)

type Metadata map[string]string

func (m Metadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string(m))
}
