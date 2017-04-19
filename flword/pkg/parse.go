package pkg

import (
	"encoding/json"
)

func ToJson(v interface{}) string {
	by, _ := json.MarshalIndent(&v, "", "  ")
	return string(by)
}
