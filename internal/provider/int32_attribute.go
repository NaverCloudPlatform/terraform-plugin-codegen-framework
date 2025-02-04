// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"fmt"

	"github.com/NaverCloudPlatform/terraform-plugin-codegen-spec/provider"

	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/convert"
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/model"
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/schema"
)

type GeneratorInt32Attribute struct {
	AssociatedExternalType *schema.AssocExtType
	OptionalRequired       convert.OptionalRequired
	CustomType             convert.CustomTypePrimitive
	DeprecationMessage     convert.DeprecationMessage
	Description            convert.Description
	Sensitive              convert.Sensitive
	Validators             convert.Validators
}

func NewGeneratorInt32Attribute(name string, a *provider.Int32Attribute) (GeneratorInt32Attribute, error) {
	if a == nil {
		return GeneratorInt32Attribute{}, fmt.Errorf("*provider.Int32Attribute is nil")
	}

	c := convert.NewOptionalRequired(a.OptionalRequired)

	ctp := convert.NewCustomTypePrimitive(a.CustomType, a.AssociatedExternalType, name)

	d := convert.NewDescription(a.Description)

	dm := convert.NewDeprecationMessage(a.DeprecationMessage)

	s := convert.NewSensitive(a.Sensitive)

	v := convert.NewValidators(convert.ValidatorTypeInt32, a.Validators.CustomValidators())

	return GeneratorInt32Attribute{
		AssociatedExternalType: schema.NewAssocExtType(a.AssociatedExternalType),
		OptionalRequired:       c,
		CustomType:             ctp,
		DeprecationMessage:     dm,
		Description:            d,
		Sensitive:              s,
		Validators:             v,
	}, nil
}

func (g GeneratorInt32Attribute) GeneratorSchemaType() schema.Type {
	return schema.GeneratorInt32Attribute
}

func (g GeneratorInt32Attribute) Imports() *schema.Imports {
	imports := schema.NewImports()

	imports.Append(g.CustomType.Imports())

	imports.Append(g.Validators.Imports())

	if g.AssociatedExternalType != nil {
		imports.Append(schema.AssociatedExternalTypeImports())
	}

	imports.Append(g.AssociatedExternalType.Imports())

	return imports
}

func (g GeneratorInt32Attribute) Equal(ga schema.GeneratorAttribute) bool {
	h, ok := ga.(GeneratorInt32Attribute)

	if !ok {
		return false
	}

	if !g.AssociatedExternalType.Equal(h.AssociatedExternalType) {
		return false
	}

	if !g.OptionalRequired.Equal(h.OptionalRequired) {
		return false
	}

	if !g.CustomType.Equal(h.CustomType) {
		return false
	}

	if !g.DeprecationMessage.Equal(h.DeprecationMessage) {
		return false
	}

	if !g.Description.Equal(h.Description) {
		return false
	}

	if !g.Sensitive.Equal(h.Sensitive) {
		return false
	}

	return g.Validators.Equal(h.Validators)
}

func (g GeneratorInt32Attribute) Schema(name schema.FrameworkIdentifier) (string, error) {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("%q: schema.Int32Attribute{\n", name))
	b.Write(g.CustomType.Schema())
	b.Write(g.OptionalRequired.Schema())
	b.Write(g.Sensitive.Schema())
	b.Write(g.Description.Schema())
	b.Write(g.DeprecationMessage.Schema())
	b.Write(g.Validators.Schema())
	b.WriteString("},")

	return b.String(), nil
}

func (g GeneratorInt32Attribute) ModelField(name schema.FrameworkIdentifier) (model.Field, error) {
	field := model.Field{
		Name:      name.ToPascalCase(),
		TfsdkName: name.ToString(),
		ValueType: model.Int32ValueType,
	}

	customValueType := g.CustomType.ValueType()

	if customValueType != "" {
		field.ValueType = customValueType
	}

	return field, nil
}

func (g GeneratorInt32Attribute) CustomTypeAndValue(name string) ([]byte, error) {
	if g.AssociatedExternalType == nil {
		return nil, nil
	}

	var buf bytes.Buffer

	int32Type := schema.NewCustomInt32Type(name)

	b, err := int32Type.Render()

	if err != nil {
		return nil, err
	}

	buf.Write(b)

	int32Value := schema.NewCustomInt32Value(name)

	b, err = int32Value.Render()

	if err != nil {
		return nil, err
	}

	buf.Write(b)

	return buf.Bytes(), nil
}

func (g GeneratorInt32Attribute) ToFromFunctions(name string) ([]byte, error) {
	if g.AssociatedExternalType == nil {
		return nil, nil
	}

	toFrom := schema.NewToFromInt32(name, g.AssociatedExternalType)

	b, err := toFrom.Render()

	if err != nil {
		return nil, err
	}

	return b, nil
}

// AttrType returns a string representation of a basetypes.Int32Typable type.
func (g GeneratorInt32Attribute) AttrType(name schema.FrameworkIdentifier) (string, error) {
	if g.AssociatedExternalType != nil {
		return fmt.Sprintf("%sType{}", name.ToPascalCase()), nil
	}

	return "basetypes.Int32Type{}", nil
}

// AttrValue returns a string representation of a basetypes.Int32Valuable type.
func (g GeneratorInt32Attribute) AttrValue(name schema.FrameworkIdentifier) string {
	if g.AssociatedExternalType != nil {
		return fmt.Sprintf("%sValue", name.ToPascalCase())
	}

	return "basetypes.Int32Value"
}

func (g GeneratorInt32Attribute) To() (schema.ToFromConversion, error) {
	if g.AssociatedExternalType != nil {
		return schema.ToFromConversion{
			AssocExtType: g.AssociatedExternalType,
		}, nil
	}

	return schema.ToFromConversion{
		Default: "ValueInt32Pointer",
	}, nil
}

func (g GeneratorInt32Attribute) From() (schema.ToFromConversion, error) {
	if g.AssociatedExternalType != nil {
		return schema.ToFromConversion{
			AssocExtType: g.AssociatedExternalType,
		}, nil
	}

	return schema.ToFromConversion{
		Default: "Int32PointerValue",
	}, nil
}
