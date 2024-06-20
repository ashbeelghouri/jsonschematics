package jsonschematics

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

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

			if nextKeyIsIndex := i < len(keys)-1 && isNumeric(keys[i+1]); nextKeyIsIndex {
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

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func stringExists(s string, slice []string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func stringsInSlice(s []string, slice []string) bool {
	for _, str := range s {
		if stringExists(str, slice) {
			return true
		}
	}
	return false
}

func isJSON(content []byte) (string, error) {
	var result interface{}
	if err := json.Unmarshal(content, &result); err != nil {
		return "", err
	}

	switch result.(type) {
	case map[string]interface{}:
		return "object", nil
	case []interface{}:
		return "array", nil
	default:
		return "unknown", fmt.Errorf("content is neither a JSON object nor array")
	}
}

func GetJson(path string) (interface{}, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to load schema file: %v", err)
		return nil, err
	}
	jsonType, err := isJSON(content)
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
		log.Fatalf("[GetJsonFileAsMap] Failed to parse the data: %v", err)
		return nil, err
	}
	return data, nil
}

func getJsonFileAsMapArray(content []byte) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	err := json.Unmarshal(content, &data)
	if err != nil {
		log.Fatalf("[GetJsonFileAsMapArray] Failed to parse the data: %v", err)
		return nil, err
	}
	return data, nil
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

func StringLikePatterns(str string, keyPatterns []string) bool {
	for _, pattern := range keyPatterns {
		re := regexp.MustCompile(ConvertKeyToRegex(pattern))
		if re.MatchString(str) {
			return true
		}
	}
	return false
}

func GetConstantMapKeys(mapper map[string]Constant) []string {
	keys := make([]string, 0, len(mapper))
	for k := range mapper {
		keys = append(keys, k)
	}
	return keys
}
