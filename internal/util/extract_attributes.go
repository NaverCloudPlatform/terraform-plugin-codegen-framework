package util

import (
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-spec/spec"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type OptionalRequestBody struct {
	Name     string                   `json:"name,omitempty"`
	Required []string                 `json:"required,omitempty"`
	Optional []*RequestParametersInfo `json:"optional,omitempty"`
}

type RequestParametersInfo struct {
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
	Format string `json:"format,omitempty"`
}

type IndividualOperationInfo struct {
	Parameters  []string             `json:"parameters,omitempty"`
	RequestBody *OptionalRequestBody `json:"request_body,omitempty"`
	Response    string               `json:"response,omitempty"`
	Method      string               `json:"method,omitempty"`
	Path        string               `json:"path,omitempty"`
}

type CRUDParameters struct {
	Create *IndividualOperationInfo   `json:"create,omitempty"`
	Read   *IndividualOperationInfo   `json:"read"`
	Update []*IndividualOperationInfo `json:"update"`
	Delete *IndividualOperationInfo   `json:"delete"`
}

type RequestParameters struct {
	Required []*RequestParametersInfo `json:"required,omitempty"`
	Optional []*RequestParametersInfo `json:"optional,omitempty"`
}

type DetailedRequestType struct {
	spec.RequestType
	Parameters  *RequestParameters `json:"parameters,omitempty"`
	RequestBody *NcloudRequestBody `json:"request_body,omitempty"`
}

type NcloudRequestBody struct {
	spec.RequestBody
	Required []*RequestParametersInfo `json:"required,omitempty"`
	Optional []*RequestParametersInfo `json:"optional,omitempty"`
}

type NcloudProvider struct {
	provider.Provider
	Endpoint string `json:"endpoint,omitempty"`
}

type NcloudSpecification struct {
	spec.Specification
	Provider    *NcloudProvider    `json:"provider"`
	Resources   []NcloudResource   `json:"resources"`
	DataSources []NcloudDataSource `json:"datasources"`
}

type NcloudResource struct {
	resource.Resource
	CRUDParameters      CRUDParameters `json:"crud_parameters"`
	RefreshObjectName   string         `json:"refresh_object_name"`
	ImportStateOverride string         `json:"import_state_override"`
	Id                  string         `json:"id"`
}

type NcloudDataSource struct {
	datasource.DataSource
	CRUDParameters      CRUDParameters `json:"crud_parameters"`
	RefreshObjectName   string         `json:"refresh_object_name"`
	ImportStateOverride string         `json:"import_state_override"`
	Id                  string         `json:"id"`
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
