package jsonschematics

import (
	v2 "github.com/ashbeelghouri/jsonschematics/data/v2"
	"github.com/ashbeelghouri/jsonschematics/utils"
	"log"
	"os"
	"testing"
)

func TestV2Validate(t *testing.T) {
	schematics, err := v2.LoadJsonSchemaFile("test-data/schema/direct/v2/example-1.json")
	if err != nil {
		t.Error(err)
	}
	content, err := os.ReadFile("test-data/data/direct/v2/example-2.json")
	if err != nil {
		t.Error(err)
	}
	jsonData, err := utils.BytesToMap(content)
	if err != nil {
		t.Error(err)
	}
	errs := schematics.Validate(jsonData)
	log.Println(errs.GetStrings("en", "%data\n"))
}
