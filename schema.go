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
	Validators  map[string]Constant `json:"validators"`
	Operators   map[string]Constant `json:"operators"`
}

type Constant struct {
	Attributes map[string]interface{} `json:"attributes"`
	ErrMsg     string                 `json:"error"`
}

func (s *Schematics) LoadSchemaFromFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to load schema file: %v", err)
		return err
	}
	schema, err := HandleSchemaVersions(content)
	s.Schema = *schema
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
		return nil
	}
	schema, err := HandleSchemaVersions(jsonData)
	s.Schema = *schema
	if err != nil {
		return err
	}
	s.Validators.BasicValidators()
	s.Operators.LoadBasicOperations()
	return nil
}

func (f *Field) Validate(value interface{}, allValidators map[string]validators.Validator) (*string, error) {
	nameOfValidator := "unknown"

	for name, constants := range f.Validators {
		if name != "" {
			nameOfValidator = name
		}

		log.Println("name of the validator is:", name)
		if stringExists(name, []string{"Exist", "Required", "IsRequired"}) {
			continue
		}
		if customValidator, exists := allValidators[name]; exists {
			if err := customValidator(value, constants.Attributes); err != nil {
				if constants.ErrMsg != "" {
					log.Printf("Validation Error: %v", err)
					return &nameOfValidator, errors.New(constants.ErrMsg)
				}
				return &nameOfValidator, err
			}
		} else {
			return &nameOfValidator, errors.New("validator not registered")
		}
	}
	return &nameOfValidator, nil
}

func (f *Field) Operate(value interface{}, allOperations map[string]operators.Op) interface{} {
	for operationName, operationConstants := range f.Operators {
		result := f.PerformOperation(value, operationName, allOperations, operationConstants)
		if result != nil {
			value = result
		}
	}
	return value
}

func (f *Field) PerformOperation(value interface{}, operation string, allOperations map[string]operators.Op, constants Constant) interface{} {
	customValidator, exists := allOperations[operation]
	if !exists {
		return nil
	}
	result := customValidator(value, constants.Attributes)
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

func (s *Schematics) Validate(data interface{}) *ErrorMessages {
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
	var missingFromDependants []string
	data := *s.makeFlat(d)
	for _, field := range s.Schema.Fields {
		allKeys := GetConstantMapKeys(field.Validators)
		matchingKeys := FindMatchingKeys(data, field.TargetKey)
		if len(matchingKeys) == 0 {
			if stringsInSlice(allKeys, []string{"MustHave", "Exist", "Required", "IsRequired"}) {
				errs.AddError("IsRequired", field.TargetKey, "is required", "")
			}
			continue
		} else if len(field.DependsOn) > 0 {
			for _, d := range field.DependsOn {
				matchDependsOn := FindMatchingKeys(data, d)
				if len(matchDependsOn) < 1 || StringLikePatterns(d, missingFromDependants) {
					errs.AddError("Depends On", field.TargetKey, "this value depends on other values which do not exists", d)
					missingFromDependants = append(missingFromDependants, field.TargetKey)
					break
				}
			}

		} else {
			for key, value := range matchingKeys {
				validator, err := field.Validate(value, s.Validators.ValidationFns)
				if err != nil {
					log.Println("validator error:", err)
					if validator != nil {
						errs.AddError(*validator, key, err.Error(), value)
					}

				}
			}
		}
	}
	if errs.HaveErrors() {
		return &errs
	}
	return nil
}

func (s *Schematics) validateArray(data []map[string]interface{}) *ErrorMessages {
	var errs ErrorMessages
	i := 0
	for _, d := range data {
		var errorMessages *ErrorMessages
		i = i + 1
		var dMap DataMap
		dMap.FlattenTheMap(d, "", s.Separator)
		arrayId, exists := dMap.Data[s.ArrayIdKey]
		if !exists {
			arrayId = fmt.Sprintf("row-%d", i)
			exists = true
		}
		log.Println("arrayID", arrayId)
		errorMessages = s.validateSingle(d)
		if errorMessages != nil {
			for _, msg := range errorMessages.Messages {
				errs.AddErrorsForArray(msg.Validator, msg.Target, msg.Message, msg.Value, arrayId)
			}
		}
	}

	if len(errs.Messages) > 0 {
		return &errs
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
