package v1

import (
	"encoding/json"
	basic "github.com/ashbeelghouri/jsonschematics/api/v0"
	"github.com/ashbeelghouri/jsonschematics/utils"
	"log"
	"os"
)

var Logs utils.Logger

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

func (s *Schema) Configs() {
	Logs = s.Logger
	if s.Logger.PrintDebugLogs {
		log.Println("debugger is on")
	}
	if s.Logger.PrintErrorLogs {
		log.Println("error logging is on")
	}
}

func LoadJsonSchemaFile(path string) (*basic.Schema, error) {
	var schema Schema
	schema.Configs()
	content, err := os.ReadFile(path)
	if err != nil {
		Logs.ERROR("Failed to load schema file", err)
		return nil, err
	}
	err = json.Unmarshal(content, &schema)
	if err != nil {
		return nil, err
	}
	return schema.transformTov0(), nil
}

func LoadMap(schemaMap interface{}) (*basic.Schema, error) {
	var s *Schema
	s.Configs()
	jsonBytes, err := json.Marshal(schemaMap)
	if err != nil {
		Logs.ERROR("Schema should be valid json map[string]interface", err)
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &s)
	if err != nil {
		return nil, err
	}
	return s.transformTov0(), nil
}

func transformComponents(components map[string]Constant) map[basic.TargetKey]basic.Constant {
	results := map[basic.TargetKey]basic.Constant{}
	for key, constant := range components {
		results[basic.TargetKey(key)] = basic.Constant{
			Attributes: constant.Attributes,
			ErrMsg:     constant.ErrMsg,
			L10n:       constant.L10n,
		}
	}
	return results
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
