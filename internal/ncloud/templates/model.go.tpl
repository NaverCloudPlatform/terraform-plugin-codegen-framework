{{ define "Model" }}
// Template for generating Terraform provider Model code
// Required data is as follows.
// RefreshObjectName string
// Model string

type {{.RefreshObjectName | ToPascalCase}}Model struct {
    ID types.String `tfsdk:"id"`
    {{.Model}}
}

{{ end }}