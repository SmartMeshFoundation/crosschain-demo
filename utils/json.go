package utils

import "encoding/json"

// ToJSONString :
func ToJSONString(v interface{}) string {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

// ToFormatJSONString :
func ToFormatJSONString(v interface{}) string {
	buf, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(buf)
}
