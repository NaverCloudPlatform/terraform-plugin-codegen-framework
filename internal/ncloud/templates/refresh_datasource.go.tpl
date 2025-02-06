{{ define "Refresh_DataSource" }}
/* =================================================================================
 * Refresh Template
 * Required data are as follows
 *
		PackageName          string
		ResourceName         string
		RefreshObjectName    string
		RefreshLogic         string
		ReadMethodName       string
		ReadReqBody          string
		Endpoint             string
		ReadPathParams       string
		ReadOpOptionalParams string
		IdGetter             string
 * ================================================================================= */

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

func (plan *{{.RefreshObjectName | ToPascalCase}}Model) refreshFromOutput(ctx context.Context, diagnostics *diag.Diagnostics) {

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

	id := {{.IdGetter}}

	var postPlan {{.RefreshObjectName | ToPascalCase}}Model

	// Fill required attributes
	*plan = postPlan
}


{{ end }}