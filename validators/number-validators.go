package validators

import (
	"fmt"
)

func IsNumber(i interface{}, _ map[string]interface{}) error {
	if _, ok := i.(float64); !ok {
		return errors.New(fmt.Sprintf("%s is not a number", i))
	}
	return nil
}

func MaxAllowed(i interface{}, attributes map[string]interface{}) error {
	number, ok := i.(float64)
	if !ok {
		return errors.New(fmt.Sprintf("%v is not a number", i))
	}
	if _max, ok := attributes["max"].(float64); !ok || number > _max {
		if !ok {
			return errors.New("max attribute is not a number")
		}
		return errors.New(fmt.Sprintf("(%v) is greater than maximum allowed: (%v)", number, _max))
	}
	return nil
}

func MinAllowed(i interface{}, attributes map[string]interface{}) error {
	number, ok := i.(float64)
	if !ok {
		return errors.New(fmt.Sprintf("%v is not a number", i))
	}
	if _min, ok := attributes["min"].(float64); !ok || number < _min {
		if !ok {
			return errors.New("min attribute is not a number")
		}
		return errors.New(fmt.Sprintf("(%v) is less than minimum allowed: (%v)", number, _min))
	}
	return nil
}

func InBetween(i interface{}, attributes map[string]interface{}) error {
	number, ok := i.(float64)
	if !ok {
		return errors.New(fmt.Sprintf("%v is not a number", i))
	}
	_min, minOk := attributes["min"].(float64)
	_max, maxOk := attributes["max"].(float64)
	if !minOk || !maxOk || number < _min || number > _max {
		if !minOk || !maxOk {
			return errors.New("min or max attribute is not a number")
		}
		return errors.New(fmt.Sprintf("(%v) should be in between (%v) and (%v)", number, _min, _max))
	}
	return nil
}
