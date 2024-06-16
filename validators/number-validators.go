package validators

import (
	"errors"
	"fmt"
)

func IsNumber(i interface{}, _ map[string]interface{}) error {
	_, ok := i.(float64)
	if !ok {
		return errors.New(fmt.Sprintf("%s is not a number", i))
	}
	return nil
}

func MaxAllowed(i interface{}, attributes map[string]interface{}) error {
	number := i.(float64)
	_max := attributes["max"].(float64)
	if number > _max {
		return errors.New(fmt.Sprintf("(%v) is greater than maximum allowed: (%v)", number, _max))
	}
	return nil
}

func MinAllowed(i interface{}, attributes map[string]interface{}) error {
	number := i.(float64)
	_min := attributes["min"].(float64)
	if number < _min {
		return errors.New(fmt.Sprintf("(%v) is less than minimum allowed: (%v)", number, _min))
	}
	return nil
}

func InBetween(i interface{}, attributes map[string]interface{}) error {
	number := i.(float64)
	_min := attributes["min"].(float64)
	_max := attributes["max"].(float64)
	if number < _min || number > _max {
		return errors.New(fmt.Sprintf("(%v) should be in between (%v) and (%v)", number, _min, _max))
	}
	return nil
}
