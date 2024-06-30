package v0

import (
	"encoding/json"
	"fmt"
	"github.com/ashbeelghouri/jsonschematics"
	"github.com/ashbeelghouri/jsonschematics/data"
	"github.com/ashbeelghouri/jsonschematics/errorHandler"
	"github.com/ashbeelghouri/jsonschematics/operators"
	"github.com/ashbeelghouri/jsonschematics/utils"
	"github.com/ashbeelghouri/jsonschematics/validators"
	"log"
	"os"
)

var Logs utils.Logger

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
	Version string           `json:"version"`
	Fields  map[string]Field `json:"fields"`
}

type Field struct {
	DependsOn             []string               `json:"depends_on"`
	Name                  string                 `json:"name"`
	Type                  string                 `json:"type"`
	IsRequired            bool                   `json:"is_required"`
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
			err.Validator = name
		}
		if f.IsRequired && value == nil {
			err.Validator = "Required"
			err.AddMessage("en", "this is a required field")
			return &err
		}

		if _, exists := allValidators[name]; exists {
			err.AddMessage("en", "validator not registered")
			return &err
		}
		if err1 := allValidators[name](value, constants.Attributes); err1 != nil {
			if !(constants.ErrMsg != "" && f.L10n != nil) {
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
	var dMap jsonschematics.DataMap
	dMap.FlattenTheMap(data, "", s.Separator)
	return &dMap.Data
}

func (s *Schematics) deflate(data map[string]interface{}) map[string]interface{} {
	return jsonschematics.DeflateMap(data, s.Separator)
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
	dataType, item := data.IsValidJson(dataBytes)
	if item == nil {
		baseError.AddMessage("en", "invalid format provided for the data, can only be map[string]interface or []map[string]interface")
		errs.AddError("whole-data", baseError)
		return &errs
	}
	if dataType == "object" {
		obj := item.(map[string]interface{})
		return s.ValidateObject(obj, nil)
	} else {
		arr := item.([]map[string]interface{})
		return s.ValidateArray(arr)
	}
}

func (s *Schematics) ValidateObject(jsonData map[string]interface{}, id *string) *errorHandler.Errors {
	var errorMessages errorHandler.Errors
	var baseError errorHandler.Error
	d := *s.makeFlat(jsonData)
	var missingFromDependants []string
	for _, field := range s.Schema.Fields {
		baseError.Validator = "is-required"
		matchingKeys := jsonschematics.FindMatchingKeys(d, field.TargetKey)
		if len(matchingKeys) == 0 {
			if field.IsRequired {
				baseError.AddMessage("en", "this field is required")
				errorMessages.AddError(field.TargetKey, baseError)
			}
			continue
		}
		//	check for dependencies
		if len(field.DependsOn) > 0 {
			missing := false
			for _, d := range field.DependsOn {
				matchDependsOn := jsonschematics.FindMatchingKeys(jsonData, d)
				if !(data.StringInStrings(field.TargetKey, missingFromDependants) == false && len(matchDependsOn) > 0) {
					baseError.AddMessage("en", "this field depends on other values which do not exists")
					errorMessages.AddError(field.TargetKey, baseError)
					missingFromDependants = append(missingFromDependants, field.TargetKey)
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
	var errs errorHandler.Errors
	i := 0
	for _, d := range jsonData {
		var errorMessages *errorHandler.Errors
		var dMap jsonschematics.DataMap
		dMap.FlattenTheMap(d, "", s.Separator)
		arrayId, exists := dMap.Data[s.ArrayIdKey]
		if !exists {
			arrayId = fmt.Sprintf("row-%d", i)
			exists = true
		}

		id := arrayId.(string)
		errorMessages = s.ValidateObject(d, &id)
		if errorMessages.HasErrors() {
			errs.MergeErrors(errorMessages)
		}
	}

	if errs.HasErrors() {
		return &errs
	}
	return nil
}
