package v0

import (
	"github.com/ashbeelghouri/jsonschematics/api/parsers"
	jsonschematics "github.com/ashbeelghouri/jsonschematics/data/v0"
	"github.com/ashbeelghouri/jsonschematics/errorHandler"
	"github.com/ashbeelghouri/jsonschematics/utils"
	"net/http"
	"regexp"
	"strings"
)

type TargetKey string
type EndpointKey string
type Name string

type Field struct {
	DependsOn  []string
	Name       string
	Type       string
	Required   bool
	Validators map[TargetKey]Constant
	Operators  map[TargetKey]Constant
	L10n       map[string]interface{}
}

type Constant struct {
	Attributes map[string]interface{} `json:"attributes"`
	ErrMsg     string                 `json:"error"`
	L10n       map[string]interface{} `json:"l10n"`
}

type Global struct {
	Headers map[TargetKey]Field
}

type Endpoint struct {
	Type    string
	Body    map[TargetKey]Field
	Headers map[TargetKey]Field
	Query   map[TargetKey]Field
}

type Schema struct {
	Version   string
	Global    Global
	Locale    string
	Logger    utils.Logger
	Endpoints map[EndpointKey]Endpoint
}

func (s *Schema) GetSchematics(fieldType string, fields *map[TargetKey]Field) (*jsonschematics.Schematics, error) {
	var schematics jsonschematics.Schematics
	FieldKeys := jsonschematics.Field{
		DependsOn:  []string{},
		Type:       fieldType,
		Validators: map[string]jsonschematics.Constant{},
		Operators:  map[string]jsonschematics.Constant{},
	}

	schema := jsonschematics.Schema{
		Version: s.Version,
		Fields:  make(map[jsonschematics.TargetKey]jsonschematics.Field),
	}

	for target, f := range *fields {
		var allValidators map[string]jsonschematics.Constant

		for key, validator := range f.Validators {
			allValidators[string(key)] = jsonschematics.Constant{
				Attributes: validator.Attributes,
				Error:      validator.ErrMsg,
				L10n:       validator.L10n,
			}
		}
		var allOperations map[string]jsonschematics.Constant
		for key, operator := range f.Operators {
			allOperations[string(key)] = jsonschematics.Constant{
				Attributes: operator.Attributes,
				Error:      operator.ErrMsg,
				L10n:       operator.L10n,
			}
		}
		FieldKeys.Type = f.Type
		FieldKeys.Validators = allValidators
		FieldKeys.Operators = allOperations
		FieldKeys.L10n = s.Global.Headers[target].L10n
		schema.Fields[jsonschematics.TargetKey(target)] = FieldKeys
	}

	schematics.Schema = schema
	return &schematics, nil
}

func (s *Schema) ValidateRequest(r *http.Request) *errorHandler.Errors {
	internalErrors := "internal-errors"

	var errorMessages errorHandler.Errors
	var errMsg errorHandler.Error
	errMsg.Validator = "request"
	errMsg.Value = "all"
	transformedRequest, err := parsers.ParseRequest(r)
	if err != nil {
		s.Logger.ERROR(err.Error())
		errMsg.AddMessage("en", "unable to transform request")
		errorMessages.AddError(internalErrors, errMsg)
		return &errorMessages
	}

	globalHeadersSchematics, err := s.GetSchematics("Global Headers", &s.Global.Headers)
	if err != nil {
		s.Logger.ERROR(err.Error())
		errMsg.AddMessage("en", "schema conversion error")
		errorMessages.AddError(internalErrors, errMsg)
		return &errorMessages
	}
	errs := globalHeadersSchematics.Validate(transformedRequest["headers"])
	if errs.HasErrors() {
		s.Logger.ERROR("all errors", err.Error())
		return errs
	}

	for path, endpoint := range s.Endpoints {
		regex := utils.GetPathRegex(string(path))
		matched, err := regexp.MatchString(regex, transformedRequest["path"].(string))
		if err != nil {
			errMsg.AddMessage("en", "path not matched - regex not matched")
			errorMessages.AddError(internalErrors, errMsg)
			return &errorMessages
		}
		if matched {
			s.Logger.DEBUG("url not matched")
			return nil
		}

		if strings.ToLower(endpoint.Type) == strings.ToLower(transformedRequest["method"].(string)) {
			headerSchematics, err := s.GetSchematics("Headers", &endpoint.Headers)
			if err != nil {
				s.Logger.ERROR(err.Error())
				errMsg.AddMessage("en", err.Error())
				errorMessages.AddError(internalErrors, errMsg)
				return &errorMessages
			}
			errs := headerSchematics.Validate(transformedRequest["headers"])
			if errs.HasErrors() {
				s.Logger.ERROR("validation errors on headers:", errs.GetStrings("en", "%validator: %message"))
				return errs
			}
			bodySchematics, err := s.GetSchematics("Body", &endpoint.Body)
			if err != nil {
				s.Logger.ERROR(err.Error())
				errMsg.AddMessage("en", err.Error())
				errorMessages.AddError(internalErrors, errMsg)
				return &errorMessages
			}
			errs = bodySchematics.Validate(transformedRequest["body"])
			if errs.HasErrors() {
				s.Logger.ERROR("validation errors on body:", errs.GetStrings("en", "%validator: %message"))
				return errs
			}
			querySchematics, err := s.GetSchematics("Query", &endpoint.Query)
			if err != nil {
				s.Logger.ERROR(err.Error())
				errMsg.AddMessage("en", err.Error())
				errorMessages.AddError(internalErrors, errMsg)
				return &errorMessages
			}
			errs = querySchematics.Validate(transformedRequest["query"])
			if errs.HasErrors() {
				s.Logger.ERROR("validation errors on query:", errs.GetStrings("en", "%validator: %message"))
				return errs
			}
		}
	}
	return nil
}
