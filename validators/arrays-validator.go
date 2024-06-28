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
	if maxLen, ok := attr["max"].(float64); !ok || maxLen < 0 {
		return errors.New("attribute 'max' must be a non-negative float64")
	} else if arrLen := reflect.ValueOf(i).Len(); arrLen > int(maxLen) {
		return fmt.Errorf("array length can not be greater than %d", int(maxLen))
	}
	return nil
}

func ArrayLengthMin(i interface{}, attr map[string]interface{}) error {
	if !isArray(i) {
		return errors.New("only arrays are allowed")
	}
	if minLen, ok := attr["min"].(float64); !ok || minLen < 0 {
		return errors.New("attribute 'min' must be a non-negative float64")
	} else if arrLen := reflect.ValueOf(i).Len(); arrLen < int(minLen) {
		return fmt.Errorf("array length can not be lesser than %d", int(minLen))
	}
	return nil
}

func StringsTakenFromOptions(i interface{}, attr map[string]interface{}) error {
	if !isArray(i) {
		return errors.New("only arrays are allowed")
	}
	STRINGS := i.([]interface{})
	for _, str := range STRINGS {
		stringDoesNotExists := StringTakenFromOptions(str, attr)
		if stringDoesNotExists != nil {
			return stringDoesNotExists
		}
	}
	return nil
}

func SpecificStringIsProvidedInArray(i interface{}, attr map[string]interface{}) error {
	if !isArray(i) {
		return errors.New("only arrays are allowed")
	}
	if _, ok := attr["shouldExists"]; !ok {
		return errors.New("attribute 'shouldExists' is not provided")
	}
	STRINGS := i.([]interface{})
	shouldExist := attr["shouldExists"].(string)
	for _, str := range STRINGS {
		if str == shouldExist {
			return nil
		}
	}
	return fmt.Errorf("the string %s is not provided in the array", shouldExist)
}
