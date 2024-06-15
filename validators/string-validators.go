package validators

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func IsString(i interface{}, _ map[string]interface{}) error {
	_, ok := i.(string)
	if !ok {
		return errors.New(fmt.Sprintf("is not a string"))
	}
	return nil
}

func NotEmpty(i interface{}, _ map[string]interface{}) error {
	str := i.(string)
	if strings.TrimSpace(str) == "" {
		return errors.New(fmt.Sprintf("this string can not be empty"))
	}
	return nil
}

func StringInArr(i interface{}, attr map[string]interface{}) error {
	str := i.(string)
	strArr := attr["arr"].([]string)
	found := false
	if len(strArr) > 0 {
		for _, item := range strArr {
			if item == str {
				found = true
			}
		}
	}
	if !found {
		return errors.New(fmt.Sprintf("string not found in array"))
	}
	return nil
}

func LIKE(i interface{}, attr map[string]interface{}) error {
	str := i.(string)

	pattern, ok := attr["pattern"].(string)
	if ok {
		replacer := strings.NewReplacer(
			".", "\\.",
			"+", "\\+",
			"?", "\\?",
			"(", "\\(",
			")", "\\)",
			"[", "\\[",
			"]", "\\]",
			"{", "\\{",
			"}", "\\}",
			"^", "\\^",
			"$", "\\$",
		)
		regexPattern := replacer.Replace(pattern)
		regexPattern = strings.ReplaceAll(regexPattern, "%", ".*")
		regexPattern = strings.ReplaceAll(regexPattern, "_", ".")
		regexPattern = "^" + regexPattern + "$"

		matched, _ := regexp.MatchString(regexPattern, str)

		if !matched {
			return errors.New(fmt.Sprintf("%s is not a LIKE %s", str))
		}
	} else {
		return errors.New(fmt.Sprintf("like pattern is invalid or is not provided"))
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
	if !re.MatchString(str) {
		return errors.New("special Characters are not allowed")
	}
	return nil
}

func HaveSpecialCharacters(i interface{}, _ map[string]interface{}) error {
	err := NoSpecialCharacters(i, nil)
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

func IsURL(i interface{}, _ map[string]interface{}) error {
	str := i.(string)
	urlRegex := regexp.MustCompile(`^(http|https)://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`)
	if !urlRegex.MatchString(str) {
		return errors.New(fmt.Sprintf("%s is not a valid url", str))
	}
	return nil
}

func IsNotURL(i interface{}, _ map[string]interface{}) error {
	str := i.(string)
	urlRegex := regexp.MustCompile(`^(http|https)://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`)
	if urlRegex.MatchString(str) {
		return errors.New(fmt.Sprintf("%s is url", str))
	}
	return nil
}

func HaveURLHostName(i interface{}, attr map[string]interface{}) error {
	str := i.(string)
	shouldHaveHost := attr["HostName"].(string)
	parsedURL, err := url.Parse(str)
	if err != nil {
		return err
	}
	hostname := parsedURL.Hostname()

	if !strings.HasSuffix(hostname, shouldHaveHost) {
		return errors.New(fmt.Sprintf("%s have different hostname than %s", str, hostname))
	}
	return nil
}

func HaveQueryParameter(i interface{}, attr map[string]interface{}) error {
	str := i.(string)
	queryParams := attr["params"].(string)
	params := strings.Split(queryParams, ",")
	parsedURL, err := url.Parse(str)
	if err != nil {
		return err
	}
	queryParameters := parsedURL.Query()
	for _, p := range params {
		_, exists := queryParameters[strings.TrimSpace(p)]
		if !exists {
			return errors.New(fmt.Sprintf("url %s have missing parameter: %s", str, p))
		}
	}
	return nil
}

func IsHttps(i interface{}, _ map[string]interface{}) error {
	str := i.(string)
	parsedURL, err := url.Parse(str)
	if err != nil {
		return err
	}

	if parsedURL.Scheme != "https" {
		return errors.New(fmt.Sprintf("url %s is not https", str))
	}

	return nil
}

func IsValidUuid(i interface{}, _ map[string]interface{}) error {
	str := i.(string)
	var uuidRegex = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
	if !uuidRegex.MatchString(str) {
		return errors.New(fmt.Sprintf("url %s is not a valid uuid", str))
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
