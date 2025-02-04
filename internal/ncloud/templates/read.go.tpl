{{ define "Read" }}
/* =================================================================================
 * Read Template
 * Required data are as follows
 *
		ResourceName      string
		RefreshObjectName string
 * ================================================================================= */

func (a *{{.ResourceName | ToCamelCase}}Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan {{.RefreshObjectName | ToPascalCase}}Model

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.refreshFromOutput(ctx, &resp.Diagnostics, plan.ID.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

{{ end }}