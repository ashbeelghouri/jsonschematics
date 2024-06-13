package main

import (
	"errors"
	"fmt"
)

func IsString(value interface{}, _ map[string]interface{}) error {
	_, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprintf("Is not a string"))
	}
	return nil
}

func IsInt(value interface{}, _ map[string]interface{}) error {
	_, ok := value.(int)
	if !ok {
		return errors.New(fmt.Sprintf("Is not an Int"))
	}
	return nil
}
