package operators

import (
	"strings"
)

func Capitalize(i interface{}, _ map[string]interface{}) *interface{} {
	str := i.(string)
	var opResult interface{} = strings.ToUpper(string(str[0])) + strings.ToLower(str[1:])
	return &opResult
}

func UpperCase(i interface{}, _ map[string]interface{}) *interface{} {
	str := i.(string)
	var opResult interface{} = strings.ToUpper(str)
	return &opResult
}
func LowerCase(i interface{}, _ map[string]interface{}) *interface{} {
	str := i.(string)
	var opResult interface{} = strings.ToLower(str)
	return &opResult
}
