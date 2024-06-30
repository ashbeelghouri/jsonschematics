package v2

import (
	"encoding/json"
	"github.com/ashbeelghouri/jsonschematics"
	basic "github.com/ashbeelghouri/jsonschematics/api/v0"
	"github.com/ashbeelghouri/jsonschematics/utils"
	"net/http"
	"os"
)

type Schema struct {
	Version   string              `json:"version"`
	Global    Global              `json:"global"`
	Endpoints map[string]Endpoint `json:"endpoints"`
	Locale    string
	Logger    utils.Logger
}

type Global struct {
	Headers []Field `json:"headers"`
}

type Endpoint struct {
	Path    string  `json:"path"`
	Type    string  `json:"type"`
	Body    []Field `json:"body"`
	Headers []Field `json:"headers"`
	Query   []Field `json:"query"`
}

type Field struct {
	DependsOn             []string
	Key                   string                 `json:"target_key"`
	Validators            []Component            `json:"validators"`
	Operators             []Component            `json:"operators"`
	L10n                  map[string]interface{} `json:"l10n"`
	AdditionalInformation map[string]interface{} `json:"additional_information"`
}

type Component struct {
	Name       string
	Attributes map[string]interface{}
	ErrMsg     string
	L10n       map[string]interface{}
}

func LoadJsonSchemaFile(path string) (*Schema, error) {
	var schema Schema
	schema.Logger = jsonschematics.Logs
	content, err := os.ReadFile(path)
	if err != nil {
		jsonschematics.Logs.ERROR("Failed to load schema file", err)
		return nil, err
	}

	err = json.Unmarshal(content, &schema)
	if err != nil {
		return nil, err
	}
	return &schema, nil
}

func transformComponents(components []Component) map[basic.TargetKey]basic.Constant {
	validators := map[basic.TargetKey]basic.Constant{}
	for _, validator := range components {
		validators[basic.TargetKey(validator.Name)] = basic.Constant{
			Attributes: validator.Attributes,
			ErrMsg:     validator.ErrMsg,
			L10n:       validator.L10n,
		}
	}
	return validators
}

func (s *Schema) transformTov0() *basic.Schema {
	var baseSchema basic.Schema
	baseSchema.Version = s.Version
	baseSchema.Locale = s.Locale
	baseSchema.Logger = s.Logger
	global := basic.Global{Headers: map[basic.TargetKey]basic.Field{}}
	endpoints := map[basic.EndpointKey]basic.Endpoint{}

	for _, field := range s.Global.Headers {
		global.Headers[basic.TargetKey(field.Key)] = basic.Field{
			DependsOn:  field.DependsOn,
			Validators: transformComponents(field.Validators),
			Operators:  transformComponents(field.Operators),
			L10n:       field.L10n,
		}
	}

	for path, endpoint := range s.Endpoints {
		headers := map[basic.TargetKey]basic.Field{}
		for _, field := range endpoint.Headers {
			headers[basic.TargetKey(field.Key)] = basic.Field{
				DependsOn:  field.DependsOn,
				Validators: transformComponents(field.Validators),
				Operators:  transformComponents(field.Operators),
				L10n:       field.L10n,
			}
		}
		body := map[basic.TargetKey]basic.Field{}
		for _, field := range endpoint.Body {
			body[basic.TargetKey(field.Key)] = basic.Field{
				DependsOn:  field.DependsOn,
				Validators: transformComponents(field.Validators),
				Operators:  transformComponents(field.Operators),
				L10n:       field.L10n,
			}
		}

		query := map[basic.TargetKey]basic.Field{}
		for _, field := range endpoint.Body {
			query[basic.TargetKey(field.Key)] = basic.Field{
				DependsOn:  field.DependsOn,
				Validators: transformComponents(field.Validators),
				Operators:  transformComponents(field.Operators),
				L10n:       field.L10n,
			}
		}
		endpoints[basic.EndpointKey(path)] = basic.Endpoint{
			Type:    endpoint.Type,
			Body:    body,
			Headers: headers,
			Query:   query,
		}
	}
	baseSchema.Global = global
	baseSchema.Endpoints = endpoints
	return &baseSchema
}

func (s *Schema) ValidateRequest(r *http.Request) *jsonschematics.ErrorMessages {
	baseSchema := s.transformTov0()
	return baseSchema.ValidateRequest(r)
}