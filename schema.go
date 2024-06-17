package jsonschematics

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (s *Schematics) LoadSchemaFromFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to load schema file: %v", err)
		return err
	}
	err = json.Unmarshal(content, &s.Schema)
	if err != nil {
		log.Fatalf("[LoadSchema] Failed to parse the data: %v", err)
		return err
	}
	s.Validators.BasicValidators()
	s.Operators.LoadBasicOperations()
	if s.Separator == "" {
		s.Separator = "."
	}
	return nil
}

func (s *Schematics) LoadSchemaFromMap(m *map[string]interface{}) error {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, &s.Schema)
	if err != nil {
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

func (s *Schematics) makeFlat(data map[string]interface{}) *map[string]interface{} {
	var dMap DataMap
	dMap.FlattenTheMap(data, "", s.Separator)
	return &dMap.Data
}

func (s *Schematics) deflate(data map[string]interface{}) map[string]interface{} {
	return DeflateMap(data, s.Separator)
}

func (s *Schematics) Validate(data interface{}) interface{} {
	switch d := data.(type) {
	case map[string]interface{}:
		return s.validateSingle(d)
	case []map[string]interface{}:
		return s.validateArray(d)
	default:
		return nil
	}
}

func (s *Schematics) validateSingle(d map[string]interface{}) *ErrorMessages {
	var errs ErrorMessages
	data := *s.makeFlat(d)
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

func (s *Schematics) validateArray(data []map[string]interface{}) *[]ArrayOfErrors {
	var msg []ArrayOfErrors
	i := 0
	for _, d := range data {
		i = i + 1
		var dMap DataMap
		dMap.FlattenTheMap(d, "", s.Separator)
		arrayId, exists := dMap.Data[s.ArrayIdKey]
		if !exists {
			arrayId = fmt.Sprintf("row-%d", i)
			exists = true
		}
		log.Println("arrayID", arrayId)
		errMessages := s.validateSingle(d)
		if errMessages != nil {
			msg = append(msg, ArrayOfErrors{
				Errors: *errMessages,
				ID:     arrayId,
			})
		}
	}

	if len(msg) > 0 {
		return &msg
	}
	return nil
}

func (s *Schematics) Operate(data interface{}) interface{} {
	switch d := data.(type) {
	case map[string]interface{}:
		return s.performOperationSingle(d)
	case []map[string]interface{}:
		return s.performOperationArray(d)
	default:
		return nil
	}
}

func (s *Schematics) performOperationSingle(data map[string]interface{}) *map[string]interface{} {
	data = *s.makeFlat(data)
	for _, field := range s.Schema.Fields {
		matchingKeys := FindMatchingKeys(data, field.TargetKey)
		for key, value := range matchingKeys {
			data[key] = field.Operate(value, s.Operators.OpFunctions)
		}
	}
	d := s.deflate(data)
	return &d
}

func (s *Schematics) performOperationArray(data []map[string]interface{}) *[]map[string]interface{} {
	var obj []map[string]interface{}
	for _, d := range data {
		results := s.performOperationSingle(d)
		obj = append(obj, *results)
	}
	if len(obj) > 0 {
		return &obj
	}
	return nil
}
