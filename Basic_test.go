package jsonschematics

import (
	"encoding/json"
	"log"
	"testing"
)

func TestLoader(t *testing.T) {
	schema, err := Load("json/schema.json")
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
