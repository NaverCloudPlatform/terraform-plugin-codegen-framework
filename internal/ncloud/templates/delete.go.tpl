{{ define "Delete" }}
/* =================================================================================
 * Delete Template
 * Required data are as follows
 *
		ResourceName      string
		RefreshObjectName string
		DeleteMethod      string
		DeleteReqBody     string
		DeleteMethodName  string
		Endpoint          string
		DeletePathParams  string
		IdGetter          string
 * ================================================================================= */

func (a *{{.ResourceName | ToCamelCase}}Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan {{.RefreshObjectName | ToPascalCase}}Model

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

 	reqParams := &ncloudsdk.{{.DeleteMethodName}}Request{
		{{.DeleteReqBody}}
	}

	tflog.Info(ctx, "Update{{.DeleteMethodName}} reqParams="+common.MarshalUncheckedString(reqParams))

	c := ncloudsdk.NewClient("{{.Endpoint}}", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))

	_, err := c.{{.DeleteMethodName}}_TF(reqParams)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}

	err = plan.waitResourceDeleted(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}
}

{{ end }}