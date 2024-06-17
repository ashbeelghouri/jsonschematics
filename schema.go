package jsonschematics

import (
	"encoding/json"
	"errors"
	"github.com/ashbeelghouri/jsonschematics/operators"
	"github.com/ashbeelghouri/jsonschematics/validators"
	"log"
	"os"
)

type Schematics struct {
	Schema     Schema
	Validators validators.Validators
	Operators  operators.Operators
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
	Operators   []string            `json:"operators"`
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
		sch.Operators.LoadBasicOperations()
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
	s.Operators.LoadBasicOperations()
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

func (f *Field) Operate(value interface{}, allOperations map[string]operators.Op) interface{} {
	for _, operator := range f.Operators {
		result := f.PerformOperation(value, operator, allOperations)
		if result != nil {
			value = result
		}
	}
	return value
}

func (f *Field) PerformOperation(value interface{}, operation string, allOperations map[string]operators.Op) interface{} {
	customValidator, exists := allOperations[operation]
	if !exists {
		return nil
	}
	constants := f.Constants["_"+operation].Attributes
	result := customValidator(value, constants)
	return *result
}

func (s *Schematics) MakeFlat(data map[string]interface{}) *map[string]interface{} {
	var dMap DataMap
	dMap.FlattenTheMap(data, "", s.Separator)
	return &dMap.Data
}

func (s *Schematics) Deflate(data map[string]interface{}) map[string]interface{} {
	return DeflateMap(data, s.Separator)
}

func (s *Schematics) Validate(d map[string]interface{}) *ErrorMessages {
	var errs ErrorMessages
	data := *s.MakeFlat(d)
	for _, field := range s.Schema.Fields {
		matchingKeys := FindMatchingKeys(data, field.TargetKey)
		if len(matchingKeys) == 0 {
			if stringsInSlice(field.Validators, []string{"MustHave", "Exist", "Required", "IsRequired"}) {
				errs.AddError("IsRequired", field.TargetKey, "is required", "")
			}
			continue
		} else {
			for key, value := range matchingKeys {
				validator, err := field.Validate(value, s.Validators.ValidationFns)
				if err != nil {
					errs.AddError(*validator, key, err.Error(), value)
				}
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
		var dMap DataMap
		dMap.FlattenTheMap(d, "", s.Separator)
		arrayId, exists := dMap.Data[s.ArrayIdKey]
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

func (s *Schematics) PerformOperations(data map[string]interface{}) *map[string]interface{} {
	data = *s.MakeFlat(data)
	for _, field := range s.Schema.Fields {
		matchingKeys := FindMatchingKeys(data, field.TargetKey)
		for key, value := range matchingKeys {
			data[key] = field.Operate(value, s.Operators.OpFunctions)
		}
	}
	d := s.Deflate(data)
	return &d
}

func (s *Schematics) PerformArrOperations(data []map[string]interface{}) *[]map[string]interface{} {
	var obj []map[string]interface{}
	for _, d := range data {
		results := s.PerformOperations(d)
		obj = append(obj, *results)
	}
	if len(obj) > 0 {
		return &obj
	}
	return nil
}
