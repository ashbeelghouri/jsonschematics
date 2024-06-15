package jsonschematics

import (
	"encoding/json"
	"log"
	"testing"
)

func TestAll(t *testing.T) {
	schema, err := LoadFromJsonFile("json/schema.json")
	if err != nil {
		t.Error(err)
	}
	log.Println(schema)
	data, err := GetJsonFileAsMap("json/data.json")
	if err != nil {
		t.Error(err)
	}
	errs := schema.Validate(*data)
	b, err := json.Marshal(errs)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	log.Println(string(b))
}
