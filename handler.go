package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

type CustomValidator func(interface{}, map[string]interface{}) error

var allValidators = map[string]CustomValidator{}

type SchemaDef struct {
	Version string  `json:"version"`
	Fields  []Field `json:"fields"`
}

type Field struct {
	Name        string              `json:"name"`
	Type        string              `json:"type"`
	DependsOn   string              `json:"depends_on"`
	TargetKey   string              `json:"target_key"`
	Description string              `json:"description"`
	Validators  []string            `json:"validators"`
	Constants   map[string]Constant `json:"constants"`
}

type Constant struct {
	Attributes map[string]interface{} `json:"attributes"`
}

func NewSchema() SchemaDef {
	return SchemaDef{}
}

func (schema *SchemaDef) LoadSchema(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to load schema file: %v", err)
		return err
	}
	err = json.Unmarshal(content, &schema)
	if err != nil {
		log.Fatalf("Failed to parse the data: %v", err)
		return err
	}
	return nil
}

func RegisterValidator(name string, fn CustomValidator) {
	allValidators[name] = fn
}

func (f *Field) Validate(value interface{}) error {
	for _, validator := range f.Validators {
		if customValidator, exists := allValidators[validator]; exists {
			if err := customValidator(value, f.Constants[validator].Attributes); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unknown validator: %s", validator)
		}
	}
	return nil
}

func (schema *SchemaDef) Validate(data map[string]interface{}) []error {
	var errorMessages []error
	var flatMap map[string]interface{}
	FlattenTheMap(data, "", ".", flatMap)
	for _, field := range schema.Fields {
		value, exists := flatMap[field.TargetKey]
		if !exists {
			if stringInSlice("IsRequired", field.Validators) {
				errorMessages = append(errorMessages, errors.New(fmt.Sprintf("%s(%s) is a required field", field.TargetKey, field.Name)))
			}
			continue // Or handle required fields
		}
		if err := field.Validate(value); err != nil {
			errorMessages = append(errorMessages, errors.New(fmt.Sprintf("Validation Failed for: %s. Error: %v", field.TargetKey, err)))
		}
	}
	if len(errorMessages) > 0 {
		return errorMessages
	}
	return nil
}
