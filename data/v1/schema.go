package v1

import (
	"encoding/json"
	v0 "github.com/ashbeelghouri/jsonschematics/data/v0"
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
	Version string  `json:"version"`
	Fields  []Field `json:"fields"`
}

type Field struct {
	DependsOn             []string               `json:"depends_on"`
	DisplayName           string                 `json:"display_name"`
	Name                  string                 `json:"name"`
	TargetKey             string                 `json:"target_key"`
	Type                  string                 `json:"type"`
	IsRequired            bool                   `json:"required"`
	Description           string                 `json:"description"`
	Validators            map[string]Component   `json:"validators"`
	Operators             map[string]Component   `json:"operators"`
	L10n                  map[string]interface{} `json:"l10n"`
	AdditionalInformation map[string]interface{} `json:"additional_information"`
}

type Component struct {
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
	s.Validators.BasicValidators()
	s.Operators.LoadBasicOperations()
}

func LoadJsonSchemaFile(path string) (*v0.Schematics, error) {
	var s *Schematics
	s.Configs()
	content, err := os.ReadFile(path)
	if err != nil {
		Logs.ERROR("Failed to load schema file", err)
		return nil, err
	}
	var schema Schema
	err = json.Unmarshal(content, &schema)
	if err != nil {
		Logs.ERROR("Failed to unmarshall schema file", err)
		return nil, err
	}
	s.Schema = schema

	return transformSchematics(*s), nil
}

func LoadMap(schemaMap interface{}) (*v0.Schematics, error) {
	var s *Schematics
	s.Configs()
	jsonBytes, err := json.Marshal(schemaMap)
	if err != nil {
		Logs.ERROR("Schema should be valid json map[string]interface", err)
		return nil, err
	}
	var schema Schema
	err = json.Unmarshal(jsonBytes, &schema)
	if err != nil {
		Logs.ERROR("Failed to unmarshall schema file", err)
		return nil, err
	}
	s.Schema = schema
	return transformSchematics(*s), nil
}

func transformSchematics(s Schematics) *v0.Schematics {
	var baseSchematics v0.Schematics
	baseSchematics.Locale = s.Locale
	baseSchematics.Logging = s.Logging
	baseSchematics.ArrayIdKey = s.ArrayIdKey
	baseSchematics.Separator = s.Separator
	baseSchematics.Validators = s.Validators
	baseSchematics.Operators = s.Operators
	baseSchematics.Schema = *transformSchema(s.Schema)
	baseSchematics.Validators.BasicValidators()
	baseSchematics.Operators.LoadBasicOperations()
	return &baseSchematics
}

func transformSchema(schema Schema) *v0.Schema {
	var baseSchema v0.Schema
	baseSchema.Version = schema.Version
	baseSchema.Fields = make(map[v0.TargetKey]v0.Field)
	for _, field := range schema.Fields {
		baseSchema.Fields[v0.TargetKey(field.TargetKey)] = v0.Field{
			DependsOn:             field.DependsOn,
			Name:                  field.Name,
			Type:                  field.Name,
			IsRequired:            field.IsRequired,
			Description:           field.Description,
			Validators:            transformComponents(field.Validators),
			Operators:             transformComponents(field.Operators),
			L10n:                  field.L10n,
			AdditionalInformation: field.AdditionalInformation,
		}
	}

	return &baseSchema
}

func transformComponents(comp map[string]Component) map[string]v0.Constant {
	con := make(map[string]v0.Constant)
	for name, c := range comp {
		con[name] = v0.Constant{
			Attributes: c.Attributes,
			Error:      c.Error,
			L10n:       c.L10n,
		}
	}
	return con
}
