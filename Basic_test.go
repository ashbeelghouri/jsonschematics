package jsonschematics

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	schema, err := LoadFromJsonFile("json/schema.json")
	if err != nil {
		t.Error(err)
	}
	//log.Println(schema)
	data, err := GetJsonFileAsMap("json/data.json")
	if err != nil {
		t.Error(err)
	}
	start := time.Now()
	errs := schema.Validate(*data)
	end := time.Now()

	log.Printf("Time taken to execute the function: %v", end.Sub(start))

	b, err := json.Marshal(errs)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	log.Println(string(b))
}

func TestDeFlatMap(t *testing.T) {
	flattened := map[string]interface{}{
		"person.name.first": "John",
		"person.name.last":  "Doe",
		"person.age":        30,
		"address.city":      "New York",
		"address.zip":       "10001",
	}

	d := &DataMap{}
	deflate := d.DeflateMap(flattened, ".")

	log.Printf("Unflattened Map: %+v\n", deflate)
}
