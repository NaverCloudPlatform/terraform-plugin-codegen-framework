{{ define "Read_DataSource" }}
/* =================================================================================
 * Read Template
 * Required data are as follows
 *
		DataSourceName    string
		RefreshObjectName string
 * ================================================================================= */

func (a *{{.DataSourceName | ToCamelCase}}DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan {{.RefreshObjectName | ToPascalCase}}Model

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.refreshFromOutput(ctx, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

{{ end }}