package v0

import (
	"github.com/ashbeelghouri/jsonschematics"
	"github.com/ashbeelghouri/jsonschematics/api"
	"github.com/ashbeelghouri/jsonschematics/api/parsers"
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
		Fields:  []jsonschematics.Field{},
	}

	for target, f := range *fields {
		FieldKeys.TargetKey = string(target)
		var allValidators map[string]jsonschematics.Constant

		for key, validator := range f.Validators {
			allValidators[string(key)] = jsonschematics.Constant{
				Attributes: validator.Attributes,
				ErrMsg:     validator.ErrMsg,
				L10n:       validator.L10n,
			}
		}
		var allOperations map[string]jsonschematics.Constant
		for key, operator := range f.Operators {
			allOperations[string(key)] = jsonschematics.Constant{
				Attributes: operator.Attributes,
				ErrMsg:     operator.ErrMsg,
				L10n:       operator.L10n,
			}
		}
		FieldKeys.Validators = allValidators
		FieldKeys.Operators = allOperations
		FieldKeys.L10n = s.Global.Headers[target].L10n
		schema.Fields = append(schema.Fields, FieldKeys)
	}

	schematics.Schema = schema
	return &schematics, nil
}

func (s *Schema) ValidateRequest(r *http.Request) *jsonschematics.ErrorMessages {
	var errorMessages jsonschematics.ErrorMessages

	transformedRequest, err := parsers.ParseRequest(r)
	if err != nil {
		s.Logger.ERROR(err.Error())
		errorMessages.AddError("Request Transformation", "request", err.Error(), "")
		return &errorMessages
	}

	globalHeadersSchematics, err := s.GetSchematics("Global Headers", &s.Global.Headers)
	if err != nil {
		s.Logger.ERROR(err.Error())
		errorMessages.AddError("Global Headers", "global.headers", err.Error(), "")
		return &errorMessages
	}
	errs := globalHeadersSchematics.Validate(transformedRequest["headers"])
	if errs.HaveErrors() {
		s.Logger.ERROR("all errors", err.Error())
		return errs
	}

	for path, endpoint := range s.Endpoints {
		regex := api.GetPathRegex(string(path))
		matched, err := regexp.MatchString(regex, transformedRequest["path"].(string))
		if err != nil {
			errorMessages.AddError("REGEX-MATCHING-FOR-ENDPOINT", "global.headers", err.Error(), "")
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
				errorMessages.AddError("Global Headers", "headers", err.Error(), "")
				return &errorMessages
			}
			errs := headerSchematics.Validate(transformedRequest["headers"])
			if errs.HaveErrors() {
				s.Logger.ERROR("validation errors on headers:", errs.ExtractAsStrings(""))
				return errs
			}
			bodySchematics, err := s.GetSchematics("Body", &endpoint.Body)
			if err != nil {
				s.Logger.ERROR(err.Error())
				errorMessages.AddError("BODY-Schema", "body", err.Error(), "")
				return &errorMessages
			}
			errs = bodySchematics.Validate(transformedRequest["body"])
			if errs.HaveErrors() {
				s.Logger.ERROR("validation errors on body:", errs.ExtractAsStrings(""))
				return errs
			}
			querySchematics, err := s.GetSchematics("Query", &endpoint.Query)
			if err != nil {
				s.Logger.ERROR(err.Error())
				errorMessages.AddError("Query-Schema", "query", err.Error(), "")
				return &errorMessages
			}
			errs = querySchematics.Validate(transformedRequest["query"])
			if errs.HaveErrors() {
				s.Logger.ERROR("validation errors on query:", errs.ExtractAsStrings(""))
				return errs
			}
		}

	}
	return nil
}
