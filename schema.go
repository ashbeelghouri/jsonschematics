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
	ArrayIdKey string
}

type Schema struct {
	Version string  `json:"version"`
	Fields  []Field `json:"fields"`
}

type Field struct {
	DependsOn   []string            `json:"depends_on"`
	TargetKey   string              `json:"target_key"`
	Description string              `json:"description"`
	Validators  []string            `json:"validators"`
	Constants   map[string]Constant `json:"constants"`
}

type Constant struct {
	Attributes map[string]interface{} `json:"attributes"`
	ErrMsg     string                 `json:"err"`
}

func LoadFromMap(s *map[string]interface{}) (*Schematics, error) {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	var sch Schematics
	err = json.Unmarshal(jsonData, &sch.Schema)
	if err != nil {
		sch.Validators.BasicValidators()
	}
	return &sch, nil
}

func LoadFromJsonFile(filePath string) (*Schematics, error) {
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
				if f.Constants[validator].ErrMsg != "" {
					log.Printf("Validation Error: %v", err)
					return &validator, errors.New(f.Constants[validator].ErrMsg)
				} else {
					return &validator, err
				}
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

func (s *Schematics) ValidateArray(data []map[string]interface{}) *[]ArrayOfErrors {
	var msg []ArrayOfErrors

	for _, d := range data {
		arrayId, exists := d[s.ArrayIdKey]
		if exists {
			errMessages := s.Validate(d)
			if errMessages != nil {
				msg = append(msg, ArrayOfErrors{
					Errors: *errMessages,
					ID:     arrayId,
				})
			}
		} else {
			msg = append(msg, ArrayOfErrors{
				Errors: ErrorMessages{
					Messages: []ErrorMessage{
						{
							Message:   "unable to validate the array, ArrayIdKey not defined in schematics",
							Validator: "ALL",
							Target:    "ALL",
						},
					},
				},
				ID: nil,
			})
		}
	}

	if len(msg) > 0 {
		return &msg
	}
	return nil
}
