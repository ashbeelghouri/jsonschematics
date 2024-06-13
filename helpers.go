package main

import "reflect"

func FlattenTheMap(data map[string]interface{}, prefix string, separator string, result map[string]interface{}) {
	for key, value := range data {
		newKey := key
		if prefix != "" {
			newKey = prefix + separator + key
		}
		if reflect.TypeOf(value).Kind() == reflect.Map {
			if nestedMap, ok := value.(map[string]interface{}); ok {
				FlattenTheMap(nestedMap, newKey, separator, result)
			}
		} else {
			result[newKey] = value
		}
	}
}

func stringInSlice(s string, slice []string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
