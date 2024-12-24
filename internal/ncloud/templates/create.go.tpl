{{ define "Create" }}
// Template for generating Terraform provider Create operation code
// Needed data is as follows.
// ResourceName string
// RefreshObjectName string
// CreateReqBody string
// CreateMethodName string
// Endpoint string
// CreatePathParams string, optional

func (a *{{.ResourceName | ToCamelCase}}Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan {{.RefreshObjectName | ToPascalCase}}Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c := ncloudsdk.NewClient("{{.Endpoint}}", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))

	reqParams := &ncloudsdk.{{.CreateMethodName}}Request{
		{{.CreateReqBody}}
	}
	
	tflog.Info(ctx, "Create{{.ResourceName | ToPascalCase}} reqParams="+common.MarshalUncheckedString(reqParams))

	response, err := c.{{.CreateMethodName}}(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("Error with {{.CreateMethodName}}_TF", err.Error())
		return
	}

	tflog.Info(ctx, "Create{{.ResourceName | ToPascalCase}} response="+common.MarshalUncheckedString(response))

	plan.refreshFromOutput_createOp(ctx, &resp.Diagnostics, response)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

{{ end }}