{{ define "ImportState" }}
/* =================================================================================
 * Import Template
 * Required data are as follows
 *
 		ResourceName     string
		ImportStateLogic string
 * ================================================================================= */

func (a *{{.ResourceName | ToCamelCase}}Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	{{.ImportStateLogic}}
}

{{ end }}