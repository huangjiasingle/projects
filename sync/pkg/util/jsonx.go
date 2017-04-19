package util

import (
	"encoding/json"
)

//-----------------------
//ToJson struct to json
func ToJson(v interface{}) string {
	vbyte, _ := json.Marshal(v)
	return string(vbyte)
}
