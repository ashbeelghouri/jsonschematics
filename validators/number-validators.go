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
