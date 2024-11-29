package util

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/NaverCloudPlatform/terraform-plugin-codegen-spec/resource"
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-spec/spec"
)

type RequestTypeWithMethodAndPath struct {
	spec.RequestType
	Method string `json:"method"`
	Path   string `json:"path"`
}

type RequestWithMethodAndPath struct {
	Create RequestTypeWithMethodAndPath    `json:"create,omitempty"`
	Read   RequestTypeWithMethodAndPath    `json:"read"`
	Update []*RequestTypeWithMethodAndPath `json:"update"`
	Delete RequestTypeWithMethodAndPath    `json:"delete"`
}

type RequestWithRefreshObjectName struct {
	Create RequestTypeWithMethodAndPath     `json:"create,omitempty"`
	Read   RequestTypeWithRefreshObjectName `json:"read"`
	Update []*RequestTypeWithMethodAndPath  `json:"update"`
	Delete RequestTypeWithMethodAndPath     `json:"delete"`
	Name   string                           `json:"name"`
	Id     string                           `json:"id"`
}

type RequestTypeWithRefreshObjectName struct {
	RequestTypeWithMethodAndPath
	Response string `json:"response"`
}

type RequestWithResponse struct {
	Requests []spec.Request
}

type CodeSpec struct {
	Provider    map[string]interface{}         `json:"provider"`
	Resources   []Resource                     `json:"resources"`
	DataSources []Resource                     `json:"datasources"`
	Requests    []RequestWithRefreshObjectName `json:"requests"`
	Version     string                         `json:"version"`
}

type Resource struct {
	Name                string `json:"name"`
	Schema              Schema `json:"schema"`
	RefreshObjectName   string `json:"dto_name"`
	ImportStateOverride string `json:"import_state_override"`
	Id                  string `json:"id"`
}

type Schema struct {
	Attributes resource.Attributes `json:"attributes"`
}

type Attribute struct {
	resource.Attribute
	Computed bool `json:"computed"`
	Optional bool `json:"optional"`
	Required bool `json:"required"`
}

type SingleNestedAttributeType struct {
	ComputedOptionalRequired string               `json:"computed_optional_required"`
	Attributes               []resource.Attribute `json:"attributes"`
}

type ListAttributeType struct {
	ComputedOptionalRequired string                   `json:"computed_optional_required"`
	ElementType              ListAttributeElementType `json:"element_type"`
}

type ListAttributeElementType struct {
	String interface{} `json:"string,omitempty"`
	Bool   interface{} `json:"bool,omitempty"`
	Int64  interface{} `json:"int64,omitempty"`
}

type ListNestedAttributeType struct {
	ComputedOptionalRequired string           `json:"computed_optional_required"`
	NestedObject             NestedObjectType `json:"nested_object"`
}

type NestedObjectType struct {
	Attributes []resource.Attribute `json:"attributes"`
}

func ExtractAttribute(file string) *CodeSpec {
	jsonFile, err := os.Open(file)
	if err != nil {
		return nil
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil
	}

	var result CodeSpec
	if err := json.Unmarshal(byteValue, &result); err != nil {
		return nil
	}

	return &result
}

func ExtractRequest(file, resourceName string) RequestWithRefreshObjectName {

	jsonFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("Failed to open JSON req file: %v", err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	var result CodeSpec
	if err := json.Unmarshal(byteValue, &result); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(result.Requests) == 0 {
		log.Fatalf("No requests found in the JSON file")
	}

	for _, val := range result.Requests {
		if val.Name == resourceName {
			return val
		}
	}

	return result.Requests[0]
}
