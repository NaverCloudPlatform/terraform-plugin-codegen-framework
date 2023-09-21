// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasource_generate

import (
	"strings"
	"text/template"

	specschema "github.com/hashicorp/terraform-plugin-codegen-spec/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-codegen-framework/internal/model"
	generatorschema "github.com/hashicorp/terraform-plugin-codegen-framework/internal/schema"
)

type GeneratorSetNestedBlock struct {
	schema.SetNestedBlock

	// The "specschema" types are used instead of the types within the attribute
	// because support for extracting custom import information is required.
	CustomType   *specschema.CustomType
	NestedObject GeneratorNestedBlockObject
	Validators   specschema.SetValidators
}

func (g GeneratorSetNestedBlock) AssocExtType() *generatorschema.AssocExtType {
	return g.NestedObject.AssociatedExternalType
}

func (g GeneratorSetNestedBlock) GeneratorSchemaType() generatorschema.Type {
	return generatorschema.GeneratorSetNestedBlock
}

func (g GeneratorSetNestedBlock) Imports() *generatorschema.Imports {
	imports := generatorschema.NewImports()

	customTypeImports := generatorschema.CustomTypeImports(g.CustomType)
	imports.Append(customTypeImports)

	for _, v := range g.Validators {
		customValidatorImports := generatorschema.CustomValidatorImports(v.Custom)
		imports.Append(customValidatorImports)
	}

	customTypeImports = generatorschema.CustomTypeImports(g.NestedObject.CustomType)
	imports.Append(customTypeImports)

	for _, v := range g.NestedObject.Validators {
		customValidatorImports := generatorschema.CustomValidatorImports(v.Custom)
		imports.Append(customValidatorImports)
	}

	for _, v := range g.NestedObject.Attributes {
		imports.Append(v.Imports())
	}

	for _, v := range g.NestedObject.Blocks {
		imports.Append(v.Imports())
	}

	// TODO: This should only be added if custom types (models) are being generated.
	imports.Append(generatorschema.AttrImports())

	imports.Append(g.NestedObject.AssociatedExternalType.Imports())

	return imports
}

func (g GeneratorSetNestedBlock) Equal(ga generatorschema.GeneratorBlock) bool {
	h, ok := ga.(GeneratorSetNestedBlock)

	if !ok {
		return false
	}

	if !g.CustomType.Equal(h.CustomType) {
		return false
	}

	if !g.Validators.Equal(h.Validators) {
		return false
	}

	if !g.NestedObject.Equal(h.NestedObject) {
		return false
	}

	return g.SetNestedBlock.Equal(h.SetNestedBlock)
}

func (g GeneratorSetNestedBlock) Schema(name string) (string, error) {
	type block struct {
		Name                    string
		TypeValueName           string
		Attributes              string
		Blocks                  string
		GeneratorSetNestedBlock GeneratorSetNestedBlock
	}

	attributesStr, err := g.NestedObject.Attributes.Schema()

	if err != nil {
		return "", err
	}

	blocksStr, err := g.NestedObject.Blocks.Schema()

	if err != nil {
		return "", err
	}

	b := block{
		Name:                    name,
		TypeValueName:           model.SnakeCaseToCamelCase(name),
		Attributes:              attributesStr,
		Blocks:                  blocksStr,
		GeneratorSetNestedBlock: g,
	}

	t, err := template.New("set_nested_block").Parse(setNestedBlockGoTemplate)
	if err != nil {
		return "", err
	}

	if _, err = addCommonBlockTemplate(t); err != nil {
		return "", err
	}

	var buf strings.Builder

	err = t.Execute(&buf, b)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (g GeneratorSetNestedBlock) ModelField(name string) (model.Field, error) {
	field := model.Field{
		Name:      model.SnakeCaseToCamelCase(name),
		TfsdkName: name,
		ValueType: model.SetValueType,
	}

	if g.CustomType != nil {
		field.ValueType = g.CustomType.ValueType
	}

	return field, nil
}

func (g GeneratorSetNestedBlock) GetAttributes() generatorschema.GeneratorAttributes {
	return g.NestedObject.Attributes
}

func (g GeneratorSetNestedBlock) GetBlocks() generatorschema.GeneratorBlocks {
	return g.NestedObject.Blocks
}
