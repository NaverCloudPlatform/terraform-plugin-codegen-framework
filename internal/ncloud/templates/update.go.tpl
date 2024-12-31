{{ define "Update" }}
// Template for generating Terraform provider Update operation code
// Required data is as follows.
// ResourceName string
// RefreshObjectName string
// UpdateReqBody string
// UpdateMethodName string
// Endpoint string
// UpdatePathParams string, optional

func (a *{{.ResourceName | ToCamelCase}}Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan {{.RefreshObjectName | ToPascalCase}}Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

 	reqParams := &ncloudsdk.{{.UpdateMethodName}}Request{
		{{.UpdateReqBody}}
	}

	tflog.Info(ctx, "Update{{.UpdateMethodName}} reqParams="+common.MarshalUncheckedString(reqParams))

	c := ncloudsdk.NewClient("{{.Endpoint}}", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))

	response, err := c.{{.UpdateMethodName}}_TF(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("UPDATING ERROR", err.Error())
		return
	}
		if response == nil {
		resp.Diagnostics.AddError("UPDATING ERROR", "response invalid")
		return
	}

	tflog.Info(ctx, "Update{{.UpdateMethodName}} response="+common.MarshalUncheckedString(response))

	plan.refreshFromOutput(&resp.Diagnostics, plan.ID.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

{{ end }}