package jsonschematics

import (
	"encoding/json"
	"fmt"
	"github.com/ashbeelghouri/jsonschematics/validators"
	"log"
	"os"
)

type Schematics struct {
	Schema     Schema
	Validators validators.Validators
	Prefix     string
	Separator  string
}

type Schema struct {
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

func Load(filePath string) (*Schematics, error) {
	var s Schematics
	err := s.LoadSchema(filePath)
	if err != nil {
		log.Fatalf("Can not load the schema: %v", err)
		return nil, err
	}
	s.Prefix = ""
	s.Separator = "."
	return &s, nil
}

func (s *Schematics) LoadSchema(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to load schema file: %v", err)
		return err
	}
	err = json.Unmarshal(content, &s.Schema)
	if err != nil {
		log.Fatalf("Failed to parse the data: %v", err)
		return err
	}
	return nil
}

func (f *Field) Validate(value interface{}, allValidators map[string]validators.Validator) error {
	for _, validator := range f.Validators {
		if customValidator, exists := allValidators[validator]; exists {
			if err := customValidator(value, f.Constants[validator].Attributes); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("validator not registered: %s", validator)
		}
	}
	return nil
}

func (s *Schematics) MakeFlat(data map[string]interface{}) *map[string]interface{} {
	var flatMap map[string]interface{}
	FlattenTheMap(data, s.Prefix, s.Separator, flatMap)
	return &flatMap
}

func (s *Schematics) Validate(d map[string]interface{}) []error {
	var errs ErrorMessages
	data := *s.MakeFlat(d)
	for _, field := range s.Schema.Fields {
		value, exists := data[field.TargetKey]
		if !exists {
			if stringsInSlice([]string{"MustHave", "Exist", "Required", "IsRequired"}, field.Validators) {
				errs.AddError(fmt.Sprintf("%s(%s) is a required field", field.TargetKey, field.Name))
			}
			continue
		}
		if err := field.Validate(value, s.Validators.ValidationFns); err != nil {
			errs.AddError(fmt.Sprintf("Validation Failed for: %s. Error: %v", field.TargetKey, err))
		}
	}
	if errs.HaveErrors() {
		return errs.Errors()
	}
	return nil
}
