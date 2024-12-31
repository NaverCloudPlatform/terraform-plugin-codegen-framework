{{ define "Refresh" }}
package {{.PackageName}}

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/ncloudsdk"
)

// Refresh code to write down manually.
// Required data is as follows.
// PackageName string
// RefreshObjectName string
// Endpoint string
// CreateMethodName string
// ReadMethodName string
// ReadReqBody string

// Diagnostics might not be Required.
// Because response type of create operation is different from read operation, reload the read response to get unified refresh data.
func (a *{{.RefreshObjectName | ToPascalCase}}Model) refreshFromOutput_createOp(ctx context.Context, diagnostics *diag.Diagnostics, createRes map[string]interface{}) {

	// id :=

	// Indicate where to get resource id from create response
	// EX) err := a.waitResourceCreated(ctx, id)

	// if err != nil {
	//	diagnostics.AddError("CREATING ERROR", err.Error())
	//	return
	// }

	var postPlan {{.RefreshObjectName | ToPascalCase}}Model

	c := ncloudsdk.NewClient("{{.Endpoint}}", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))
	response, err := c.{{.ReadMethodName}}_TF(&ncloudsdk.{{.ReadMethodName}}Request{
			{{.ReadReqBody}}
	})

	// if err != nil {
	//	diagnostics.AddError("CREATING ERROR", err.Error())
	//	return
	// }

	// Fill required attributes
	// Copy(&postPlan, response)
	// postPlan.ID = types.StringValue(response.ApiKey.Attributes()["api_key_id"].String())

	*a = postPlan
}

func (a *{{.RefreshObjectName | ToPascalCase}}Model) refreshFromOutput(diagnostics *diag.Diagnostics, id string) {

	c := ncloudsdk.NewClient("{{.Endpoint}}", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))
	response, err := c.{{.ReadMethodName}}_TF(&ncloudsdk.{{.ReadMethodName}}Request{
			{{.ReadReqBody}}
	})

	// if err != nil {
	//	 diagnostics.AddError("CREATING ERROR", err.Error())
	//	 return
	// }

	var postPlan {{.RefreshObjectName | ToPascalCase}}Model

	// Fill required attributes
	// Copy(&postPlan, response)
	// postPlan.ID = types.StringValue(response.ApiKey.Attributes()["api_key_id"].String())
	// objectRes, _ := types.ObjectValueFrom(context.TODO(), a.ApiKey.AttributeTypes(context.TODO()), res)

	*a = postPlan
}

{{ end }}