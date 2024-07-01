package utils

import (
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var ExcludedValidators = []string{
	"REQUIRED",
	"IS_REQUIRED",
	"IS-REQUIRED",
	"ISREQUIRED",
}

type DataMap struct {
	Data map[string]interface{}
}

func (d *DataMap) FlattenTheMap(data map[string]interface{}, prefix string, separator string) {
	if d.Data == nil {
		d.Data = make(map[string]interface{})
	}
	if separator == "" {
		separator = "."
	}
	for key, value := range data {
		newKey := key
		if prefix != "" {
			newKey = prefix + separator + key
		}
		switch reflect.TypeOf(value).Kind() {
		case reflect.Map:
			if nestedMap, ok := value.(map[string]interface{}); ok {
				d.FlattenTheMap(nestedMap, newKey, separator)
			}
		case reflect.Slice:
			s := reflect.ValueOf(value)
			for i := 0; i < s.Len(); i++ {
				arrayKey := newKey + separator + strconv.Itoa(i)
				if nestedMap, ok := s.Index(i).Interface().(map[string]interface{}); ok {
					d.FlattenTheMap(nestedMap, arrayKey, separator)
				} else {
					d.Data[arrayKey] = s.Index(i).Interface()
				}
			}
		default:
			d.Data[newKey] = value
		}
	}
}

func DeflateMap(data map[string]interface{}, separator string) map[string]interface{} {
	result := make(map[string]interface{})

	for flatKey, value := range data {
		keys := strings.Split(flatKey, separator)
		subMap := result

		for i := 0; i < len(keys)-1; i++ {
			key := keys[i]

			if nextKeyIsIndex := i < len(keys)-1 && IsNumeric(keys[i+1]); nextKeyIsIndex {
				if _, exists := subMap[key]; !exists {
					subMap[key] = []interface{}{}
				}

				if reflect.TypeOf(subMap[key]).Kind() != reflect.Slice {
					subMap[key] = []interface{}{}
				}

				slice := subMap[key].([]interface{})
				index, _ := strconv.Atoi(keys[i+1])
				for len(slice) <= index {
					slice = append(slice, map[string]interface{}{})
				}
				subMap[key] = slice

				subMap = slice[index].(map[string]interface{})
				i++
			} else {
				if _, exists := subMap[key]; !exists {
					subMap[key] = map[string]interface{}{}
				}

				subMap = subMap[key].(map[string]interface{})
			}
		}

		subMap[keys[len(keys)-1]] = value
	}

	return result
}

func IsNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func StringInStrings(str string, slice []string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

func StringsInSlice(s []string, slice []string) bool {
	for _, str := range s {
		if StringInStrings(str, slice) {
			return true
		}
	}
	return false
}

func ConvertKeyToRegex(key string) string {
	// Escape special regex characters in the key except for *
	escapedKey := regexp.QuoteMeta(key)
	// Replace * with \d+ to match array indices
	regexPattern := strings.ReplaceAll(escapedKey, `\*`, `\d+`)
	// Add start and end of line anchors
	regexPattern = "^" + regexPattern + "$"
	return regexPattern
}

func FindMatchingKeys(data map[string]interface{}, keyPattern string) map[string]interface{} {
	matchingKeys := make(map[string]interface{})
	re := regexp.MustCompile(ConvertKeyToRegex(keyPattern))
	for key, value := range data {
		if re.MatchString(key) {
			matchingKeys[key] = value
		}
	}
	return matchingKeys
}

func IsValidJson(content []byte) (string, interface{}) {
	var arr []map[string]interface{}
	var obj map[string]interface{}
	const IsArray = "array"
	const IsObject = "object"

	if err := json.Unmarshal(content, &arr); err == nil {
		return IsArray, arr
	}

	if err := json.Unmarshal(content, &obj); err == nil {
		return IsObject, obj
	}
	return "invalid format", nil
}

func GetPathRegex(path string) string {
	path = strings.ReplaceAll(path, "*", ".*")
	path = strings.ReplaceAll(path, ":", "[^/]+")
	return "^" + path + "$"
}

func FormatError(id *string, message string, target string, validator string, value string, format string, data *map[string]interface{}) string {
	errorMessage := strings.Replace(format, "%message", message, -1)
	errorMessage = strings.Replace(errorMessage, "%target", target, -1)
	errorMessage = strings.Replace(errorMessage, "%validator", validator, -1)
	if id != nil {
		errorMessage = strings.Replace(errorMessage, "%id", validator, -1)
	}
	if data != nil {
		marshalled, err := json.Marshal(data)
		if err == nil {
			d := string(marshalled)
			errorMessage = strings.Replace(errorMessage, "%data", d, -1)
		}

	}
	errorMessage = strings.Replace(errorMessage, "%value", value, -1)
	return errorMessage
}

func BytesToMap(content []byte) (interface{}, error) {
	jsonType, obj := IsValidJson(content)
	if obj == nil {
		return nil, errors.New("invalid json file content found")
	}
	switch jsonType {
	case "object":
		mapper, err := getJsonFileAsMap(content)
		if err != nil {
			return nil, err
		}
		return mapper, err
	case "array":
		mapper, err := getJsonFileAsMapArray(content)
		if err != nil {
			return nil, err
		}
		return mapper, err
	default:
		return nil, errors.New("unknown json file content found")
	}
}

func getJsonFileAsMap(content []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getJsonFileAsMapArray(content []byte) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	err := json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
