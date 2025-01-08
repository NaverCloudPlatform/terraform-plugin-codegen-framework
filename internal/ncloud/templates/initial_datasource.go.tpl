{{ define "Initial_DataSource"}}
// Template for generating Terraform provider Initial code for Data Source
// Needed data is as follows.
// DataSourceName string
// ProviderName string

var (
	_ datasource.DataSource              = &{{.DataSourceName | ToCamelCase}}DataSource{}
	_ datasource.DataSourceWithConfigure = &{{.DataSourceName | ToCamelCase}}DataSource{}
)

func New{{.DataSourceName | ToPascalCase}}DataSource() datasource.DataSource {
	return &{{.DataSourceName | ToCamelCase}}DataSource{}
}

type {{.DataSourceName | ToCamelCase}}DataSource struct {
	config *conn.ProviderConfig
}

func (b *{{.DataSourceName | ToCamelCase}}DataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*conn.ProviderConfig)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	b.config = config
}

func (b *{{.DataSourceName | ToCamelCase}}DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_{{.ProviderName}}_{{.DataSourceName}}"
}

func (b *{{.DataSourceName | ToCamelCase}}DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = {{.DataSourceName | ToPascalCase}}DataSourceSchema(ctx)
}

{{ end }}