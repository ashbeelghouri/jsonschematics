package jsonschematics

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
)

type DataMap struct {
	Data map[string]interface{}
}

func (d *DataMap) FlattenTheMap(data map[string]interface{}, prefix string, separator string) {
	// Ensure result is not nil
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
		// Check if the value is a map and recursively flatten it
		if reflect.TypeOf(value).Kind() == reflect.Map {
			if nestedMap, ok := value.(map[string]interface{}); ok {
				d.FlattenTheMap(nestedMap, newKey, separator)
			}
		} else {
			// Assign the value to the result map
			d.Data[newKey] = value
		}
	}
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

func GetJsonFileAsMap(path string) (*map[string]interface{}, error) {
	var data map[string]interface{}
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to load schema file: %v", err)
		return nil, err
	}
	err = json.Unmarshal(content, &data)
	if err != nil {
		log.Fatalf("Failed to parse the data: %v", err)
		return nil, err
	}
	return &data, nil
}
