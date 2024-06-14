package jsonschematics

import (
	"encoding/json"
	"errors"
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
	DependsOn   []string            `json:"depends_on"`
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
	s.Validators.BasicValidators()
	return nil
}

func (f *Field) Validate(value interface{}, allValidators map[string]validators.Validator) (*string, error) {
	for _, validator := range f.Validators {
		if stringExists(validator, []string{"Exist", "Required", "IsRequired"}) {
			continue
		}
		if customValidator, exists := allValidators[validator]; exists {
			if err := customValidator(value, f.Constants[validator].Attributes); err != nil {
				return &validator, err
			}
		} else {
			return &validator, errors.New("validator not registered")
		}
	}
	return nil, nil
}

func (s *Schematics) MakeFlat(data map[string]interface{}) *map[string]interface{} {
	var dMap DataMap
	dMap.FlattenTheMap(data, s.Prefix, s.Separator)
	return &dMap.Data
}

func (s *Schematics) Validate(d map[string]interface{}) *ErrorMessages {
	var errs ErrorMessages
	data := *s.MakeFlat(d)
	for _, field := range s.Schema.Fields {
		value, exists := data[field.TargetKey]
		if !exists {
			if stringsInSlice([]string{"MustHave", "Exist", "Required", "IsRequired"}, field.Validators) {
				errs.AddError("IsRequired", "field.TargetKey", "is required")
			}
			continue
		} else {
			validator, err := field.Validate(value, s.Validators.ValidationFns)
			if err != nil {
				errs.AddError(*validator, field.TargetKey, err.Error())
			}
		}
	}
	if errs.HaveErrors() {
		return &errs
	}
	return nil
}
