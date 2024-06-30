package data

import (
	"encoding/json"
	v0 "github.com/ashbeelghouri/jsonschematics/data/v0"
)

func IsValidJson(content []byte) (string, interface{}) {
	var arr []map[string]interface{}
	var obj map[string]interface{}
	const IsArray = "array"
	const IsObject = "object"

	if err := json.Unmarshal(content, &arr); err == nil {
		return IsArray, arr
	}

	if err := json.Unmarshal(content, &obj); err == nil {
		return IsObject, obj
	}
	return "invalid format", nil
}

func GetConstantMapKeys(mapper map[string]v0.Constant) []string {
	keys := make([]string, 0, len(mapper))
	for k := range mapper {
		keys = append(keys, k)
	}
	return keys
}

func StringInStrings(str string, slice []string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}
