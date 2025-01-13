{{ define "Refresh_DataSource" }}
// Template for generating Terraform provider Refresh operation code
// Needed data is as follows.
// RefreshObjectName string
// PackageName string
// RefreshLogic string
// ReadMethodName string
// ReadReqBody string
// Endpoint string
// ReadPathParams string, optional

package {{.PackageName}}

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (plan *{{.RefreshObjectName | ToPascalCase}}Model) refreshFromOutput(diagnostics *diag.Diagnostics, id string) {

	c := ncloudsdk.NewClient("{{.Endpoint}}", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))

	r := &ncloudsdk.{{.ReadMethodName}}Request{
		{{.ReadReqBody}}
	}

	{{.ReadOpOptionalParams}}

	response, err := c.{{.ReadMethodName}}_TF(r)

	if err != nil {
		 diagnostics.AddError("CREATING ERROR", err.Error())
		 return
	}

	var postPlan {{.RefreshObjectName | ToPascalCase}}Model

	// Fill required attributes
	ncloudsdk.Copy(&postPlan, response)

	*plan = postPlan
}


{{ end }}