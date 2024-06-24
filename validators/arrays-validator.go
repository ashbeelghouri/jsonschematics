package validators

import (
	"errors"
	"fmt"
	"reflect"
)

func isArray(i interface{}) bool {
	val := reflect.ValueOf(i)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		return true
	default:
		return false
	}
}

func ArrayLengthMax(i interface{}, attr map[string]interface{}) error {
	if !isArray(i) {
		return errors.New("only arrays are allowed")
	}
	arrLen := reflect.ValueOf(i).Len()
	maxLen := attr["max"].(float64)
	if arrLen > int(maxLen) {
		return errors.New(fmt.Sprintf("Array length can not be greater than %d", int(maxLen)))
	}

	return nil
}
func ArrayLengthMin(i interface{}, attr map[string]interface{}) error {
	if !isArray(i) {
		return errors.New("only arrays are allowed")
	}
	arrLen := reflect.ValueOf(i).Len()
	minLen := attr["min"].(float64)
	if arrLen < int(minLen) {
		return errors.New(fmt.Sprintf("Array length can not be lesser than %d", int(minLen)))
	}

	return nil
}
