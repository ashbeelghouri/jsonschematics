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
	schema.Logging.PrintErrorLogs = true
	schema.Logging.PrintDebugLogs = true
	if err != nil {
		t.Error(err)
	}
	data, err := GetJson("json/data.json")
	if err != nil {
		t.Error(err)
	}
	start := time.Now()
	errs := schema.Validate(data)
	log.Printf("[SINGLE OBJ] Validation Time: %v", time.Since(start))
	log.Print("[SINGLE OBJ] have single errors: ", errs.HaveSingleError("", ""))
	errorsFromValidate, err := json.Marshal(errs)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	log.Println("[SINGLE OBJ] errorsFromValidate: ", string(errorsFromValidate))
	start = time.Now()
	newData := schema.Operate(data)
	log.Printf("[SINGLE OBJ] Operaions Time: %v", time.Since(start))
	log.Printf("[SINGLE OBJ] Updated DATA: %v", newData)

	log.Printf("[SINGLE OBJ] total time taken: %v", time.Since(fnTimeStart))
	log.Println("-------------------------------------------")
}

func TestForArrayData(t *testing.T) {
	fnTimeStart := time.Now()
	var schema1 Schematics
	schema1.Logging.PrintErrorLogs = true
	schema1.Logging.PrintDebugLogs = true
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
	log.Printf("[ARRAY OF OBJ] Validation Time: %v", time.Since(start))
	if errs != nil {
		obj, err := json.Marshal(errs)
		if err != nil {
			log.Fatalf("err: %v", err)
		}
		log.Printf("array validations >>>> %v", string(obj))
	} else {
		start = time.Now()
		newData := schema1.Operate(data)
		log.Printf("[ARRAY OF OBJ] Operation Time: %v", time.Since(start))
		log.Printf("[ARRAY OF OBJ] Updated Data: %v", newData)
	}
	log.Printf("[ARRAY OF OBJ] total time taken: %v", time.Since(fnTimeStart))
	log.Println("-------------------------------------------")
}

func TestNestedArrays(t *testing.T) {
	fnTimeStart := time.Now()
	var schema Schematics
	schema.Logging.PrintErrorLogs = true
	schema.Logging.PrintDebugLogs = true
	err := schema.LoadSchemaFromFile("json/arr-inside-obj-schema.json")
	if err != nil {
		t.Error(err)
	}
	data, err := GetJson("json/arr-inside-obj-data.json")
	if err != nil {
		t.Error(err)
	}
	start := time.Now()
	errs := schema.Validate(data)
	log.Printf("[TestNestedArrays] Validation Time: %v", time.Since(start))
	if errs != nil {
		jsonErrors, err := json.Marshal(errs)
		if err != nil {
			log.Fatalf("[TestNestedArrays] err: %v", err)
		}
		log.Println("[TestNestedArrays] json errors:", string(jsonErrors))
	}

	start = time.Now()
	newData := schema.Operate(data)
	log.Printf("[TestNestedArrays] Operation Time: %v", time.Since(start))
	log.Println("[TestNestedArrays] after operations:", newData)
	log.Println("[TestNestedArrays] total time taken:", time.Since(fnTimeStart))
}

func TestDeepValidationInArray(t *testing.T) {
	fnTimeStart := time.Now()
	var schema Schematics
	schema.Logging.PrintErrorLogs = true
	schema.Logging.PrintDebugLogs = true
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
	start := time.Now()
	errs := schema.Validate(data)
	log.Printf("[TestDeepValidationInArray] Validation Time: %v", time.Since(start))

	if errs != nil {
		jsonErrors, err := json.Marshal(errs)
		if err != nil {
			log.Fatalf("[TestDeepValidationInArray] err: %v", err)
		}
		log.Println("[TestDeepValidationInArray] json errors:", string(jsonErrors))
	}
	start = time.Now()
	newData := schema.Operate(data)
	log.Printf("[TestDeepValidationInArray] Operation Time: %v", time.Since(start))
	log.Println("[TestDeepValidationInArray] after operations:", newData)
	log.Println("[TestDeepValidationInArray] total time taken:", time.Since(fnTimeStart))
}

func TestSchemaVersioning(t *testing.T) {
	fnTimeStart := time.Now()
	var schema Schematics
	err := schema.LoadSchemaFromFile("json/schema-v1.1.json")
	if err != nil {
		log.Println("[TestSchemaVersioning] unable to load the schema from json file: ", err)
		t.Error(err)
	}
	log.Println("Schema Version 1.1", schema)

	log.Println("[TestSchemaVersioning] total time taken:", time.Since(fnTimeStart))
}
