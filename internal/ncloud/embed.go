package ncloud

import (
	_ "embed"
)

//go:embed templates/initial_resource.go.tpl
var InitialTemplate string

//go:embed templates/create.go.tpl
var CreateTemplate string

//go:embed templates/read_resource.go.tpl
var ReadTemplate string

//go:embed templates/update.go.tpl
var UpdateTemplate string

//go:embed templates/delete.go.tpl
var DeleteTemplate string

//go:embed templates/model_resource.go.tpl
var ModelTemplate string

//go:embed templates/refresh_resource.go.tpl
var RefreshTemplate string

//go:embed templates/wait.go.tpl
var WaitTemplate string

//go:embed templates/test_resource.go.tpl
var TestTemplate string

//go:embed templates/import.go.tpl
var ImportStateTemplate string

//go:embed templates/initial_datasource.go.tpl
var InitialTemplateDataSource string

//go:embed templates/read_datasource.go.tpl
var ReadTemplateDataSource string

//go:embed templates/model_datasource.go.tpl
var ModelTemplateDataSource string

//go:embed templates/refresh_datasource.go.tpl
var RefreshTemplateDataSource string

//go:embed templates/test_datasource.go.tpl
var TestTemplateDataSource string
