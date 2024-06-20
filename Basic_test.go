package jsonschematics

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

func TestForObjData(t *testing.T) {
	fnTimeStart := time.Now()
	var schema Schematics
	err := schema.LoadSchemaFromFile("json/schema.json")
	if err != nil {
		t.Error(err)
	}
	data, err := GetJson("json/data.json")
	if err != nil {
		t.Error(err)
	}
	start := time.Now()
	errs := schema.Validate(data)
	end := time.Now()

	log.Printf("[SINGLE OBJ] Validation Time: %v", end.Sub(start))

	errorsFromValidate, err := json.Marshal(errs)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	log.Println("[SINGLE OBJ] errorsFromValidate: ", string(errorsFromValidate))
	start = time.Now()
	newData := schema.Operate(data)
	end = time.Now()
	log.Printf("[SINGLE OBJ] Operaions Time: %v", end.Sub(start))
	log.Printf("[SINGLE OBJ] Updated DATA: %v", newData)

	log.Printf("[SINGLE OBJ] total time taken: %v", time.Now().Sub(fnTimeStart))
	log.Println("-------------------------------------------")
}

func TestForArrayData(t *testing.T) {
	fnTimeStart := time.Now()
	var schema1 Schematics
	err := schema1.LoadSchemaFromFile("json/schema.json")
	schema1.ArrayIdKey = "user.id"
	if err != nil {
		t.Error(err)
	}
	data, err := GetJson("json/arr-data.json")
	if err != nil {
		t.Error(err)
	}
	start := time.Now()
	errs := schema1.Validate(data)
	end := time.Now()

	log.Printf("[ARRAY OF OBJ] Validation Time: %v", end.Sub(start))
	if errs != nil {
		obj, err := json.Marshal(errs)
		if err != nil {
			log.Fatalf("err: %v", err)
		}
		log.Printf("array validations >>>> %v", string(obj))
	} else {
		start = time.Now()
		newData := schema1.Operate(data)
		end = time.Now()
		log.Printf("[ARRAY OF OBJ] Operation Time: %v", end.Sub(start))
		log.Printf("[ARRAY OF OBJ] Updated Data: %v", newData)
	}
	log.Printf("[ARRAY OF OBJ] total time taken: %v", time.Now().Sub(fnTimeStart))
	log.Println("-------------------------------------------")
}

func TestNestedArrays(t *testing.T) {
	fnTimeStart := time.Now()
	var schema Schematics
	err := schema.LoadSchemaFromFile("json/arr-inside-obj-schema.json")
	if err != nil {
		t.Error(err)
	}
	data, err := GetJson("json/arr-inside-obj-data.json")
	if err != nil {
		t.Error(err)
	}
	errs := schema.Validate(data)
	if errs != nil {
		jsonErrors, err := json.Marshal(errs)
		if err != nil {
			log.Fatalf("[TestNestedArrays] err: %v", err)
		}
		log.Println("[TestNestedArrays] json errors:", string(jsonErrors))
	}

	newData := schema.Operate(data)
	log.Println("[TestNestedArrays] after operations:", newData)
	log.Println("[TestNestedArrays] total time taken:", time.Now().Sub(fnTimeStart))
}

func TestDeepValidationInArray(t *testing.T) {
	fnTimeStart := time.Now()
	var schema Schematics
	err := schema.LoadSchemaFromFile("json/arr-inside-obj-schema.json")
	if err != nil {
		log.Println("[TestDeepValidationInArray] unable to load the schema from json file: ", err)
		t.Error(err)
	}
	data, err := GetJson("json/arr-inside-arr-obj-data.json")
	if err != nil {
		log.Println("[TestDeepValidationInArray] unable to load the data from json file: ", err)
		t.Error(err)
	}
	errs := schema.Validate(data)

	if errs != nil {
		jsonErrors, err := json.Marshal(errs)
		if err != nil {
			log.Fatalf("[TestDeepValidationInArray] err: %v", err)
		}
		log.Println("[TestDeepValidationInArray] json errors:", string(jsonErrors))
	}
	newData := schema.Operate(data)
	log.Println("[TestDeepValidationInArray] after operations:", newData)
	log.Println("[TestDeepValidationInArray] total time taken:", time.Now().Sub(fnTimeStart))
}
