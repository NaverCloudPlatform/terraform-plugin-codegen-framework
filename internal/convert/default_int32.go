// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package convert

import (
	"fmt"

	"github.com/NaverCloudPlatform/terraform-plugin-codegen-spec/code"
	specschema "github.com/NaverCloudPlatform/terraform-plugin-codegen-spec/schema"

	generatorschema "github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/schema"
)

const defaultInt32Import = "github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"

type DefaultInt32 struct {
	int32Default *specschema.Int32Default
}

func NewDefaultInt32(b *specschema.Int32Default) DefaultInt32 {
	return DefaultInt32{
		int32Default: b,
	}
}

func (d DefaultInt32) Equal(other DefaultInt32) bool {
	return d.int32Default.Equal(other.int32Default)
}

func (d DefaultInt32) Imports() *generatorschema.Imports {
	imports := generatorschema.NewImports()

	if d.int32Default == nil {
		return imports
	}

	if d.int32Default.Static != nil {
		imports.Add(code.Import{
			Path: defaultInt32Import,
		})
	}

	if d.int32Default.Custom != nil {
		for _, i := range d.int32Default.Custom.Imports {
			if len(i.Path) > 0 {
				imports.Add(i)
			}
		}
	}

	return imports
}

func (d DefaultInt32) Schema() []byte {
	if d.int32Default == nil {
		return nil
	}

	if d.int32Default.Static != nil {
		return []byte(fmt.Sprintf("Default: int32default.StaticInt32(%d),\n", *d.int32Default.Static))
	}

	if d.int32Default.Custom != nil && d.int32Default.Custom.SchemaDefinition != "" {
		return []byte(fmt.Sprintf("Default: %s,\n", d.int32Default.Custom.SchemaDefinition))
	}

	return nil
}
