package jsonschematics

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

func TestForObjData(t *testing.T) {
	fnTimeStart := time.Now()
	schema, err := LoadFromJsonFile("json/schema.json")
	if err != nil {
		t.Error(err)
	}
	data, err := GetJsonFileAsMap("json/data.json")
	if err != nil {
		t.Error(err)
	}
	start := time.Now()
	errs := schema.Validate(*data)
	end := time.Now()

	log.Printf("[SINGLE OBJ] Validation Time: %v", end.Sub(start))

	_, err = json.Marshal(errs)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	start = time.Now()
	newData := schema.PerformOperations(*data)
	end = time.Now()
	log.Printf("[SINGLE OBJ] Operaions Time: %v", end.Sub(start))
	log.Printf("[SINGLE OBJ] Updated DATA: %v", newData)

	log.Printf("[SINGLE OBJ] total time taken: %v", time.Now().Sub(fnTimeStart))
	log.Println("-------------------------------------------")
}

func TestForArrayData(t *testing.T) {
	fnTimeStart := time.Now()
	schema1, err := LoadFromJsonFile("json/schema.json")
	schema1.ArrayIdKey = "user.id"
	if err != nil {
		t.Error(err)
	}
	data, err := GetJsonFileAsMapArray("json/arr-data.json")
	if err != nil {
		t.Error(err)
	}
	start := time.Now()
	errs := schema1.ValidateArray(*data)
	end := time.Now()

	log.Printf("[ARRAY OF OBJ] Validation Time: %v", end.Sub(start))
	if errs != nil && len(*errs) > 0 {
		for _, j := range *errs {
			obj, err := json.Marshal(j)
			if err != nil {
				log.Fatalf("err: %v", err)
			}
			log.Printf("array validations >>>> %v", string(obj))
		}
	}
	if errs == nil || !(len(*errs) > 0) {
		start = time.Now()
		newData := schema1.PerformArrOperations(*data)
		end = time.Now()
		log.Printf("[ARRAY OF OBJ] Operation Time: %v", end.Sub(start))
		log.Printf("[ARRAY OF OBJ] Updated Data: %v", newData)
	}
	log.Printf("[ARRAY OF OBJ] total time taken: %v", time.Now().Sub(fnTimeStart))
	log.Println("-------------------------------------------")
}

func TestNestedArrays(t *testing.T) {
	//fnTimeStart := time.Now()
	schema, err := LoadFromJsonFile("json/arr-inside-obj-schema.json")
	if err != nil {
		t.Error(err)
	}
	data, err := GetJsonFileAsMap("json/arr-inside-obj-data.json")
	if err != nil {
		t.Error(err)
	}
	flatData := schema.MakeFlat(*data)
	log.Println("flat data:", flatData)

	deflated := schema.Deflate(*flatData)
	log.Println("________________________")
	log.Println("deflated data: ", deflated)
	log.Println("________________________")

	errs := schema.Validate(*data)

	if errs != nil {
		jsonErrors, err := json.Marshal(errs)
		if err != nil {
			log.Fatalf("err: %v", err)
		}
		log.Println("json errors:", string(jsonErrors))
	}

	newData := schema.PerformOperations(*data)
	log.Println("after operations:", newData)

}
