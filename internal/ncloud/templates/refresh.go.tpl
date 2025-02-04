{{ define "Refresh" }}
/* =================================================================================
 * Refresh Template
 * Required data are as follows
 *
		PackageName       string
		RefreshObjectName string
		Endpoint          string
		CreateMethodName  string
		ReadMethodName    string
		ReadReqBody       string
		IdGetter          string
		RefreshWithResponse string
 * ================================================================================= */

package {{.PackageName}}

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/ncloudsdk"
)

// Diagnostics might not be Required.
// Because response type of create operation is different from read operation, reload the read response to get unified refresh data.
func (plan *{{.RefreshObjectName | ToPascalCase}}Model) refreshFromOutput_createOp(ctx context.Context, diagnostics *diag.Diagnostics, createRes map[string]interface{}) {

	// Allocate resource id from create response
	id := {{.IdGetter}}

	// Indicate where to get resource id from create response
	err := plan.waitResourceCreated(ctx, id)

	if err != nil {
		diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	var postPlan {{.RefreshObjectName | ToPascalCase}}Model

	c := ncloudsdk.NewClient("{{.Endpoint}}", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))
	response, err := c.{{.ReadMethodName}}_TF(ctx, &ncloudsdk.Primitive{{.ReadMethodName}}Request{
			{{.ReadReqBody}}
	})

	if err != nil {
		diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	// Fill required attributes
	{{.RefreshWithResponse}}

	*plan = postPlan
}

func (plan *{{.RefreshObjectName | ToPascalCase}}Model) refreshFromOutput(ctx context.Context, diagnostics *diag.Diagnostics, id string) {

	c := ncloudsdk.NewClient("{{.Endpoint}}", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))
	response, err := c.{{.ReadMethodName}}_TF(ctx, &ncloudsdk.Primitive{{.ReadMethodName}}Request{
			{{.ReadReqBody}}
	})

	if err != nil {
		 diagnostics.AddError("CREATING ERROR", err.Error())
		 return
	}

	var postPlan {{.RefreshObjectName | ToPascalCase}}Model

	// Fill required attributes
	{{.RefreshWithResponse}}

	*plan = postPlan
}

{{ end }}