package validators

import (
	"errors"
	"fmt"
)

func IsInt(value interface{}, _ map[string]interface{}) error {
	_, ok := value.(int)
	if !ok {
		return errors.New(fmt.Sprintf("Is not an Int"))
	}
	return nil
}

func MaxAllowed(value interface{}, attributes map[string]interface{}) error {
	number := value.(int)
	_max := attributes["max"].(int)
	if number > _max {
		return errors.New(fmt.Sprintf("number is greater than maximum allowed: (%d)", _max))
	}
	return nil
}

func MinAllowed(value interface{}, attributes map[string]interface{}) error {
	number := value.(int)
	_min := attributes["min"].(int)
	if number < _min {
		return errors.New(fmt.Sprintf("number is less than minimum allowed: (%d)", _min))
	}
	return nil
}

func InBetween(value interface{}, attributes map[string]interface{}) error {
	number := value.(int)
	_min := attributes["min"].(int)
	_max := attributes["max"].(int)
	if number < _min || number > _max {
		return errors.New(fmt.Sprintf("number should be in between (%d) and (%d)", _min, _max))
	}
	return nil
}
