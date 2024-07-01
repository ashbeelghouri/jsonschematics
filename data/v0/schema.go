package v0

import (
	"encoding/json"
	"fmt"
	"github.com/ashbeelghouri/jsonschematics/errorHandler"
	"github.com/ashbeelghouri/jsonschematics/operators"
	"github.com/ashbeelghouri/jsonschematics/utils"
	"github.com/ashbeelghouri/jsonschematics/validators"
	"log"
	"os"
	"strings"
)

var Logs utils.Logger

type TargetKey string

type Schematics struct {
	Schema     Schema
	Validators validators.Validators
	Operators  operators.Operators
	Separator  string
	ArrayIdKey string
	Locale     string
	Logging    utils.Logger
}

type Schema struct {
	Version string              `json:"version"`
	Fields  map[TargetKey]Field `json:"fields"`
}

type Field struct {
	DependsOn             []string               `json:"depends_on"`
	DisplayName           string                 `json:"display_name"`
	Name                  string                 `json:"name"`
	Type                  string                 `json:"type"`
	IsRequired            bool                   `json:"required"`
	Description           string                 `json:"description"`
	Validators            map[string]Constant    `json:"validators"`
	Operators             map[string]Constant    `json:"operators"`
	L10n                  map[string]interface{} `json:"l10n"`
	AdditionalInformation map[string]interface{} `json:"additional_information"`
}

type Constant struct {
	Attributes map[string]interface{} `json:"attributes"`
	Error      string                 `json:"error"`
	L10n       map[string]interface{} `json:"l10n"`
}

func (s *Schematics) Configs() {
	Logs = s.Logging
	if s.Logging.PrintDebugLogs {
		log.Println("debugger is on")
	}
	if s.Logging.PrintErrorLogs {
		log.Println("error logging is on")
	}
	s.Validators.Logger = Logs
	s.Operators.Logger = Logs
}

func (s *Schematics) LoadJsonSchemaFile(path string) error {
	s.Configs()
	content, err := os.ReadFile(path)
	if err != nil {
		Logs.ERROR("Failed to load schema file", err)
		return err
	}
	var schema Schema
	err = json.Unmarshal(content, &schema)
	if err != nil {
		Logs.ERROR("Failed to unmarshall schema file", err)
		return err
	}
	Logs.DEBUG("Schema Loaded From File: ", schema)
	s.Schema = schema
	s.Validators.BasicValidators()
	s.Operators.LoadBasicOperations()
	if s.Separator == "" {
		s.Separator = "."
	}
	if s.Locale == "" {
		s.Locale = "en"
	}
	return nil
}

func (s *Schematics) LoadMap(schemaMap interface{}) error {
	JSON, err := json.Marshal(schemaMap)
	if err != nil {
		Logs.ERROR("Schema should be valid json map[string]interface", err)
		return err
	}
	var schema Schema
	err = json.Unmarshal(JSON, &schema)
	if err != nil {
		Logs.ERROR("Invalid Schema", err)
		return err
	}
	Logs.DEBUG("Schema Loaded From MAP: ", schema)
	s.Schema = schema
	s.Validators.BasicValidators()
	s.Operators.LoadBasicOperations()
	if s.Separator == "" {
		s.Separator = "."
	}
	if s.Locale == "" {
		s.Locale = "en"
	}
	return nil
}

func (f *Field) Validate(value interface{}, allValidators map[string]validators.Validator, id *string) *errorHandler.Error {
	var err errorHandler.Error
	err.Value = value
	err.ID = id
	err.Validator = "unknown"
	for name, constants := range f.Validators {
		if name != "" {
			Logs.DEBUG("Name of the validator is not given: ", name)
			err.Validator = name
		}
		if f.IsRequired && value == nil {
			err.Validator = "Required"
			err.AddMessage("en", "this is a required field")
			Logs.DEBUG("ERR: ", err)
			return &err
		}

		if utils.StringInStrings(strings.ToUpper(name), utils.ExcludedValidators) {
			continue
		}

		var fn validators.Validator
		fn, exists := allValidators[name]
		if !exists {
			log.Println("does not exists here!!", name)
			err.AddMessage("en", "validator not registered")
			return &err
		}
		if err1 := fn(value, constants.Attributes); err1 != nil {
			if !(constants.Error != "" && f.L10n != nil) {
				err.AddMessage("en", err1.Error())
				return &err
			}
			for locale, msg := range f.L10n {
				if msg == nil {
					err.AddMessage(locale, msg.(string))
				}
			}
			return &err
		}
	}
	return nil
}

func (s *Schematics) makeFlat(data map[string]interface{}) *map[string]interface{} {
	var dMap utils.DataMap
	dMap.FlattenTheMap(data, "", s.Separator)
	return &dMap.Data
}

func (s *Schematics) deflate(data map[string]interface{}) map[string]interface{} {
	return utils.DeflateMap(data, s.Separator)
}

func (s *Schematics) Validate(jsonData interface{}) *errorHandler.Errors {
	var baseError errorHandler.Error
	var errs errorHandler.Errors
	baseError.Validator = "validate-object"
	dataBytes, err := json.Marshal(jsonData)
	if err != nil {
		baseError.AddMessage("en", "data is not valid json")
		errs.AddError("whole-data", baseError)
		return &errs
	}
	dataType, item := utils.IsValidJson(dataBytes)
	if item == nil {
		baseError.AddMessage("en", "invalid format provided for the data, can only be map[string]interface or []map[string]interface")
		errs.AddError("whole-data", baseError)
		return &errs
	}
	if dataType == "object" {
		obj, ok := item.(map[string]interface{})
		if !ok {
			baseError.AddMessage("en", "invalid format provided for the data, can only be map[string]interface or []map[string]interface")
			errs.AddError("whole-data-obj", baseError)
			return &errs
		}
		Logs.DEBUG("validating the object", obj)
		return s.ValidateObject(obj, nil)
	} else {
		arr, ok := item.([]map[string]interface{})
		if !ok {
			baseError.AddMessage("en", "invalid format provided for the data, can only be map[string]interface or []map[string]interface")
			errs.AddError("whole-data-arr", baseError)
			return &errs
		}
		Logs.DEBUG("validating the array", arr)
		return s.ValidateArray(arr)
	}
}

func (s *Schematics) ValidateObject(jsonData map[string]interface{}, id *string) *errorHandler.Errors {
	log.Println("validating the object")
	var errorMessages errorHandler.Errors
	var baseError errorHandler.Error
	flatData := *s.makeFlat(jsonData)
	var missingFromDependants []string
	for target, field := range s.Schema.Fields {
		baseError.Validator = "is-required"
		matchingKeys := utils.FindMatchingKeys(flatData, string(target))
		if len(matchingKeys) == 0 {
			if field.IsRequired {
				baseError.AddMessage("en", "this field is required")
				errorMessages.AddError(string(target), baseError)
			}
			continue
		}
		//	check for dependencies
		if len(field.DependsOn) > 0 {
			missing := false
			for _, d := range field.DependsOn {
				matchDependsOn := utils.FindMatchingKeys(flatData, d)
				if !(utils.StringInStrings(string(target), missingFromDependants) == false && len(matchDependsOn) > 0) {
					log.Println(matchDependsOn)
					baseError.Validator = "depends-on"
					baseError.AddMessage("en", "this field depends on other values which do not exists")
					errorMessages.AddError(string(target), baseError)
					missingFromDependants = append(missingFromDependants, string(target))
					missing = true
					break
				}
			}
			if missing {
				continue
			}
		}

		for key, value := range matchingKeys {
			validationError := field.Validate(value, s.Validators.ValidationFns, id)
			Logs.DEBUG(validationError)
			Logs.DEBUG("validation error with pointer", *validationError)
			if validationError != nil {
				errorMessages.AddError(key, *validationError)
			}
		}

	}

	if errorMessages.HasErrors() {
		return &errorMessages
	}
	return nil
}

func (s *Schematics) ValidateArray(jsonData []map[string]interface{}) *errorHandler.Errors {
	Logs.DEBUG("validating the array")
	var errs errorHandler.Errors
	i := 0
	for _, d := range jsonData {
		var errorMessages *errorHandler.Errors
		var dMap utils.DataMap
		dMap.FlattenTheMap(d, "", s.Separator)
		arrayId, exists := dMap.Data[s.ArrayIdKey]
		if !exists {
			arrayId = fmt.Sprintf("row-%d", i)
			exists = true
		}

		id := arrayId.(string)
		errorMessages = s.ValidateObject(d, &id)
		if errorMessages.HasErrors() {
			log.Println("has errors", errorMessages.GetStrings("en", "%data\n"))
			errs.MergeErrors(errorMessages)
		}
		i = i + 1
	}

	if errs.HasErrors() {
		return &errs
	}
	return nil
}

// operators

func (f *Field) Operate(value interface{}, allOperations map[string]operators.Op) interface{} {
	for operationName, operationConstants := range f.Operators {
		customValidator, exists := allOperations[operationName]
		if !exists {
			Logs.ERROR("This operation does not exists in basic or custom operators", operationName)
			return nil
		}
		result := customValidator(value, operationConstants.Attributes)
		if result != nil {
			value = result
		}
	}
	return value
}

func (s *Schematics) Operate(data interface{}) (interface{}, *errorHandler.Errors) {
	var errorMessages errorHandler.Errors
	var baseError errorHandler.Error
	baseError.Validator = "operate-on-schema"
	bytes, err := json.Marshal(data)
	if err != nil {
		Logs.ERROR("[operate] error converting the data into bytes", err)
		baseError.AddMessage("en", "data is not valid json")
		errorMessages.AddError("whole-data", baseError)
		return nil, &errorMessages
	}

	dataType, item := utils.IsValidJson(bytes)
	if item == nil {
		Logs.ERROR("[operate] error occurred when checking if this data is an array or object")
		baseError.AddMessage("en", "can not convert the data into json")
		errorMessages.AddError("whole-data", baseError)
		return nil, &errorMessages
	}

	if dataType == "object" {
		obj := item.(map[string]interface{})
		results := s.OperateOnObject(obj)
		if results != nil {
			return results, nil
		} else {
			baseError.AddMessage("en", "operation on object unsuccessful")
			errorMessages.AddError("whole-data", baseError)
			return nil, &errorMessages
		}
	} else if dataType == "array" {
		arr := item.([]map[string]interface{})
		results := s.OperateOnArray(arr)
		if results != nil && len(*results) > 0 {
			return results, nil
		} else {
			baseError.AddMessage("en", "operation on array unsuccessful")
			errorMessages.AddError("whole-data", baseError)
			return nil, &errorMessages
		}
	}

	return data, nil
}

func (s *Schematics) OperateOnObject(data map[string]interface{}) *map[string]interface{} {
	data = *s.makeFlat(data)
	for target, field := range s.Schema.Fields {
		matchingKeys := utils.FindMatchingKeys(data, string(target))
		for key, value := range matchingKeys {
			data[key] = field.Operate(value, s.Operators.OpFunctions)
		}
	}
	d := s.deflate(data)
	return &d
}

func (s *Schematics) OperateOnArray(data []map[string]interface{}) *[]map[string]interface{} {
	var obj []map[string]interface{}
	for _, d := range data {
		results := s.OperateOnObject(d)
		obj = append(obj, *results)
	}
	if len(obj) > 0 {
		return &obj
	}
	return nil
}

// General

func (s *Schematics) MergeFields(sc2 *Schematics) *Schematics {
	for target, field := range sc2.Schema.Fields {
		if s.Schema.Fields[target].Type == "" {
			s.Schema.Fields[target] = field
		}
	}
	return s
}
