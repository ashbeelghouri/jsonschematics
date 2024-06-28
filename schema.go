package jsonschematics

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ashbeelghouri/jsonschematics/operators"
	"github.com/ashbeelghouri/jsonschematics/utils"
	"github.com/ashbeelghouri/jsonschematics/validators"
)

var logs utils.Logger

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
	Version string  `json:"version"`
	Fields  []Field `json:"fields"`
}

type Field struct {
	DependsOn             []string               `json:"depends_on"`
	Name                  string                 `json:"name"`
	Type                  string                 `json:"type"`
	TargetKey             string                 `json:"target_key"`
	Description           string                 `json:"description"`
	Validators            map[string]Constant    `json:"validators"`
	Operators             map[string]Constant    `json:"operators"`
	L10n                  map[string]interface{} `json:"l10n"`
	AdditionalInformation map[string]interface{} `json:"additional_information"`
}

type Constant struct {
	Attributes map[string]interface{} `json:"attributes"`
	ErrMsg     string                 `json:"error"`
	L10n       map[string]interface{} `json:"l10n"`
}

func (s *Schematics) Configs() {
	logs = s.Logging
	if s.Logging.PrintDebugLogs {
		log.Println("debugger is on")
	}
	if s.Logging.PrintErrorLogs {
		log.Println("error logging is on")
	}
	s.Validators.Logger = logs
	s.Operators.Logger = logs
}

func (s *Schematics) LoadSchemaFromFile(path string) error {
	s.Configs()
	content, err := os.ReadFile(path)
	if err != nil {
		logs.ERROR("Failed to load schema file", err)
		return err
	}
	schema, err := HandleSchemaVersions(content)
	if err != nil {
		return err
	}
	s.Schema = *schema
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

func (s *Schematics) LoadSchemaFromMap(m *map[string]interface{}) error {
	s.Configs()
	jsonData, err := json.Marshal(m)
	if err != nil {
		logs.ERROR("Failed to load schema file", err)
		return nil
	}
	schema, err := HandleSchemaVersions(jsonData)
	s.Schema = *schema
	if err != nil {
		logs.ERROR("Failed to load schema file", err)
		return err
	}

	s.Validators.BasicValidators()
	logs.DEBUG("basic validator loaded")
	s.Operators.LoadBasicOperations()
	if s.Separator == "" {
		logs.DEBUG("separator set to '.'")
		s.Separator = "."
	}
	if s.Locale == "" {
		logs.DEBUG("locale set to 'en'")
		s.Locale = "en"
	}
	logs.DEBUG("loaded the file successfully")
	return nil
}

func (f *Field) Validate(value interface{}, allValidators map[string]validators.Validator, locale *string) (*string, error) {
	logs.DEBUG("validation is being performed on:", f.Name, fmt.Sprintf("[%s]", f.TargetKey))
	nameOfValidator := "unknown"
	for name, constants := range f.Validators {
		if name != "" {
			nameOfValidator = name
		} else {
			logs.ERROR("name of the validator is not defined!", name, "constants are:", constants)
		}

		logs.DEBUG("name of the validator is:", name)
		if stringExists(name, []string{"Exist", "Required", "IsRequired"}) {
			logs.DEBUG("skipping required validator as it has already been checked out")
			continue
		}
		if customValidator, exists := allValidators[name]; exists {
			logs.DEBUG("validating with", name)
			if err := customValidator(value, constants.Attributes); err != nil {
				logs.DEBUG("we have an error from our validator", err)
				if constants.ErrMsg != "" {
					logs.ERROR("Validation Error", err)
					var localeError = constants.ErrMsg
					if locale != nil && *locale != "" && *locale != "en" {
						logs.DEBUG("locale is loaded and is configured")
						_, ok := f.L10n[*locale].(string)
						if !ok {
							localeError = constants.ErrMsg
						} else {
							localeError = f.L10n[*locale].(string)
						}
					}
					logs.DEBUG("custom error from the schema is being sent")
					return &nameOfValidator, errors.New(localeError)
				}
				logs.DEBUG("sending the error from validation function")
				return &nameOfValidator, err
			}
		} else {
			logs.DEBUG("this validator is not registered", &nameOfValidator)
			return &nameOfValidator, errors.New("validator not registered")
		}
	}
	return &nameOfValidator, nil
}

func (f *Field) Operate(value interface{}, allOperations map[string]operators.Op) interface{} {
	logs.DEBUG("operation is being performed on:", f.Name, fmt.Sprintf("[%s]", f.TargetKey))
	for operationName, operationConstants := range f.Operators {
		logs.DEBUG("performing operation:", operationName)
		result := f.PerformOperation(value, operationName, allOperations, operationConstants)
		if result != nil {
			logs.ERROR("operation successful", result)
			value = result
		}
	}
	logs.DEBUG("all operations are performed on", f.TargetKey)
	return value
}

func (f *Field) PerformOperation(value interface{}, operation string, allOperations map[string]operators.Op, constants Constant) interface{} {
	customValidator, exists := allOperations[operation]
	if !exists {
		logs.ERROR("This operation does not exists in basic or custom operators", operation)
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
	var upperLevelErrors ErrorMessages

	bytes, err := json.Marshal(data)
	if err != nil {
		logs.ERROR("error converting the data into bytes", err)
		upperLevelErrors.AddError("BYTES", "MARSHAL DATA", err.Error(), "validate")
		return &upperLevelErrors
	}

	dataType, item := canConvert(bytes)
	if item == nil {
		logs.ERROR("error occurred when checking if this data is an array or object")
		errMsg := "unknown error"
		upperLevelErrors.AddError("BYTES", "DETERMINE_IS_JSON", errMsg, "validate")
		return &upperLevelErrors
	}
	logs.DEBUG("data type is:", dataType)
	if dataType == "object" {
		logs.DEBUG("data is an object")
		if obj, ok := item.(map[string]interface{}); ok {
			return s.validateSingle(obj)
		} else {
			logs.ERROR("unable to recognize the object for validations")
			upperLevelErrors.AddError("BYTES", "IS UNKNOWN TYPE", "unable to recognize the object for validation", "validate")
			return &upperLevelErrors
		}

	} else if dataType == "array" {
		logs.DEBUG("data is an array")
		if obj, ok := item.([]map[string]interface{}); ok {
			return s.validateArray(obj)
		} else {
			logs.ERROR("unable to recognize the array for validations")
			upperLevelErrors.AddError("BYTES", "IS UNKNOWN TYPE", "unable to recognize the array for validation", "validate")
			return &upperLevelErrors
		}
	} else {
		upperLevelErrors.AddError("BYTES", "IS UNKNOWN TYPE", "MUST PROVIDE VALID OBJ map[string]interface{} OR []map[string]interface{}", "validate")
		return &upperLevelErrors
	}
}

func (s *Schematics) validateSingle(d map[string]interface{}) *ErrorMessages {
	var errs ErrorMessages
	var missingFromDependants []string
	data := *s.makeFlat(d)
	for _, field := range s.Schema.Fields {
		allKeys := GetConstantMapKeys(field.Validators)
		matchingKeys := FindMatchingKeys(data, field.TargetKey)
		logs.DEBUG("matching keys to the target", matchingKeys)
		if len(matchingKeys) == 0 {
			logs.DEBUG("matching key not found")
			if stringsInSlice(allKeys, []string{"MustHave", "Exist", "Required", "IsRequired"}) {
				errs.AddError("IsRequired", field.TargetKey, "is required", "")
			}
			continue
		} else if len(field.DependsOn) > 0 {
			logs.DEBUG("checking for the pre-requisites")
			for _, d := range field.DependsOn {
				matchDependsOn := FindMatchingKeys(data, d)
				if len(matchDependsOn) < 1 || StringLikePatterns(d, missingFromDependants) {
					logs.ERROR("the field on which this field depends on not found", matchDependsOn)
					errs.AddError("Depends On", field.TargetKey, "this value depends on other values which do not exists", d)
					missingFromDependants = append(missingFromDependants, field.TargetKey)
					break
				}
			}

		} else {
			for key, value := range matchingKeys {
				validator, err := field.Validate(value, s.Validators.ValidationFns, &s.Locale)
				if err != nil {
					logs.ERROR("validator error occurred:", err)
					var fieldName = key
					if s.Locale != "" && s.Locale != "en" {
						logs.DEBUG("locale is different than en:", s.Locale)
						fieldNameLocales, ok := field.L10n["name"].(map[string]interface{})
						if ok {
							_, ok = fieldNameLocales[s.Locale].(string)
							if ok {
								fieldName = fieldNameLocales[s.Locale].(string)
							}
						}
					}
					if validator != nil {
						logs.DEBUG("this is a validation error", err, "adding to the errors")
						errs.AddError(*validator, fieldName, err.Error(), value)
					}
				}
			}
		}
	}
	if errs.HaveErrors() {
		return &errs
	}
	logs.DEBUG("found no errors")
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
			logs.DEBUG("array does not have ids defined, so defining them by row number")
			arrayId = fmt.Sprintf("row-%d", i)
			exists = true
		}
		logs.DEBUG("arrayID", arrayId)
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
	var upperLevelErrors ErrorMessages
	bytes, err := json.Marshal(data)
	if err != nil {
		logs.ERROR("[operate] error converting the data into bytes", err)
		upperLevelErrors.AddError("BYTES", "MARSHAL DATA", err.Error(), "operate")
		return &upperLevelErrors
	}

	dataType, item := canConvert(bytes)
	if item == nil {
		logs.ERROR("[operate] error occurred when checking if this data is an array or object")
		errMsg := "unknown error"
		upperLevelErrors.AddError("BYTES", "DETERMINE_IS_JSON", errMsg, "operate")
		return &upperLevelErrors
	}

	if dataType == "object" {
		logs.DEBUG("[operate] data is an object")
		if obj, ok := item.(map[string]interface{}); ok {
			return s.performOperationSingle(obj)
		} else {
			logs.ERROR("unable to recognize the object for operations")
			upperLevelErrors.AddError("BYTES", "IS UNKNOWN TYPE", "unable to recognize the object for operation", "operate")
			return &upperLevelErrors
		}

	} else if dataType == "array" {
		logs.DEBUG("[operate] data is an array")
		if obj, ok := item.([]map[string]interface{}); ok {
			return s.performOperationArray(obj)
		} else {
			logs.ERROR("unable to recognize the array for operations")
			upperLevelErrors.AddError("BYTES", "IS UNKNOWN TYPE", "unable to recognize the array for operation", "operate")
			return &upperLevelErrors
		}
	} else {
		upperLevelErrors.AddError("BYTES", "IS UNKNOWN TYPE", "MUST PROVIDE VALID OBJ map[string]interface{} OR []map[string]interface{}", "operate")
		return &upperLevelErrors
	}
}

func (s *Schematics) performOperationSingle(data map[string]interface{}) *map[string]interface{} {
	logs.DEBUG("performing all operations")
	data = *s.makeFlat(data)
	for _, field := range s.Schema.Fields {
		matchingKeys := FindMatchingKeys(data, field.TargetKey)
		logs.DEBUG("matching keys for operations", matchingKeys)
		for key, value := range matchingKeys {
			data[key] = field.Operate(value, s.Operators.OpFunctions)
			logs.DEBUG("data after operation", data[key])
		}
	}
	d := s.deflate(data)
	logs.DEBUG("deflated data", d)
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

func (s *Schematics) MergeFields(sc2 *Schematics) *Schematics {
	s.Schema.Fields = append(s.Schema.Fields, sc2.Schema.Fields...)
	return s
}
