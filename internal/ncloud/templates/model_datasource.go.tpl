{{ define "Model_DataSource" }}
/* =================================================================================
 * Model Template
 * Required data are as follows
 *
		RefreshObjectName string
		Model             string
 * ================================================================================= */

type {{.RefreshObjectName | ToPascalCase}}Model struct {
    ID types.String `tfsdk:"id"`
    {{.Model}}
}

{{ end }}