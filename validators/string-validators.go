package validators

import (
	"errors"
	"fmt"
	"regexp"
)

func IsString(value interface{}, _ map[string]interface{}) error {
	_, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprintf("Is not a string"))
	}
	return nil
}

func IsEmail(i interface{}, _ map[string]interface{}) error {
	str := i.(string)
	const pattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(str) {
		return errors.New(fmt.Sprintf("%s is not a valid email address", str))
	}
	return nil
}

func MaxLengthAllowed(i interface{}, attr map[string]interface{}) error {
	str := i.(string)
	length, ok := attr["max"].(float64)
	intLength := int(length)
	if !ok {
		return errors.New("max is not provided as an int in attributes of schema")
	}
	if len(str) > intLength {
		return errors.New(fmt.Sprintf("length of the string is greater than %d", intLength))
	}
	return nil
}

func MinLengthAllowed(i interface{}, attr map[string]interface{}) error {
	str := i.(string)
	length, ok := attr["min"].(float64)
	intLength := int(length)
	if !ok {
		return errors.New("min is not provided as an int in attributes of schema")
	}
	if len(str) < intLength {
		return errors.New(fmt.Sprintf("length of the string is less than %d", intLength))
	}
	return nil
}

func InBetweenLengthAllowed(i interface{}, attr map[string]interface{}) error {
	str := i.(string)
	minlength, ok := attr["min"].(float64)
	if !ok {
		return errors.New("min is not provided as an int in attributes of schema")
	}
	maxlength, ok := attr["max"].(float64)
	if !ok {
		return errors.New("max is not provided as an int in attributes of schema")
	}
	intMinLength := int(minlength)
	intMaxLength := int(maxlength)
	if len(str) < intMinLength || len(str) > intMaxLength {
		return errors.New(fmt.Sprintf("length of the string should be greater than %d and less than %s", intMinLength, intMaxLength))
	}
	return nil
}

func NoSpecialCharacters(i interface{}, _ map[string]interface{}) error {
	str := i.(string)
	pattern := `[^a-zA-Z0-9]`
	re := regexp.MustCompile(pattern)
	if re.MatchString(str) {
		return errors.New("special Characters are not allowed")
	}
	return nil
}

func HaveSpecialCharacters(i interface{}, attr map[string]interface{}) error {
	err := NoSpecialCharacters(i, attr)
	if err == nil {
		return errors.New("special characters are required")
	}
	return nil
}

func LeastOneUpperCase(i interface{}, _ map[string]interface{}) error {
	str := i.(string)
	pattern := `.*[A-Z].*`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(str) {
		return errors.New("at least one uppercase letter is required")
	}
	return nil
}

func LeastOneLowerCase(i interface{}, _ map[string]interface{}) error {
	str := i.(string)
	pattern := `.*[a-z].*`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(str) {
		return errors.New("at least one lowercase letter is required")
	}
	return nil
}

func LeastOneDigit(i interface{}, _ map[string]interface{}) error {
	str := i.(string)
	pattern := `.*\d.*`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(str) {
		return errors.New("at least one numeric digit is required")
	}
	return nil
}

func Regex(i interface{}, attr map[string]interface{}) error {
	str := i.(string)
	pattern := attr["regex"].(string)
	re := regexp.MustCompile(pattern)
	if !re.MatchString(str) {
		return errors.New("regex failed")
	}
	return nil
}
