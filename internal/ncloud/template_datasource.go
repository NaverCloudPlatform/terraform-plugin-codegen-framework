package ncloud

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
)

type DataSourceTemplate struct {
	spec              util.NcloudSpecification
	providerName      string
	resourceName      string
	importStateLogic  string
	refreshObjectName string
	model             string
	refreshLogic      string
	endpoint          string
	readPathParams    string
	readMethod        string
	readReqBody       string
	readMethodName    string
	idGetter          string
	funcMap           template.FuncMap
}

// RenderCreate implements BaseTemplate.
func (d *DataSourceTemplate) RenderCreate() []byte {
	panic("Data source doesn't provide this method. Please use the right method.")
}

// RenderDelete implements BaseTemplate.
func (d *DataSourceTemplate) RenderDelete() []byte {
	panic("Data source doesn't provide this method. Please use the right method.")
}

// RenderImportState implements BaseTemplate.
func (d *DataSourceTemplate) RenderImportState() []byte {
	panic("Data source doesn't provide this method. Please use the right method.")
}

// RenderInitial implements BaseTemplate.
func (d *DataSourceTemplate) RenderInitial() []byte {
	var b bytes.Buffer

	initialTemplate, err := template.New("").Funcs(d.funcMap).Parse(InitialTemplateDataSource)
	if err != nil {
		log.Fatalf("error occurred with baseTemplate at rendering initial: %v", err)
	}

	data := struct {
		ProviderName string
		ResourceName string
	}{
		ProviderName: d.providerName,
		ResourceName: d.resourceName,
	}

	err = initialTemplate.ExecuteTemplate(&b, "Initial_DataSource", data)
	if err != nil {
		log.Fatalf("error occurred with generating Initial template: %v", err)
	}

	return b.Bytes()
}

// RenderModel implements BaseTemplate.
func (d *DataSourceTemplate) RenderModel() []byte {
	var b bytes.Buffer

	modelTemplate, err := template.New("").Funcs(d.funcMap).Parse(ModelTemplateDataSource)
	if err != nil {
		log.Fatalf("error occurred with baseTemplate at rendering model: %v", err)
	}

	data := struct {
		RefreshObjectName string
		Model             string
	}{
		RefreshObjectName: d.refreshObjectName,
		Model:             d.model,
	}

	err = modelTemplate.ExecuteTemplate(&b, "Model_DataSource", data)
	if err != nil {
		log.Fatalf("error occurred with Generating Model: %v", err)
	}

	return b.Bytes()
}

// RenderRead implements BaseTemplate.
func (d *DataSourceTemplate) RenderRead() []byte {
	var b bytes.Buffer

	readTemplate, err := template.New("").Funcs(d.funcMap).Parse(ReadTemplateDataSource)
	if err != nil {
		log.Fatalf("error occurred with baseTemplate at rendering read: %v", err)
	}

	data := struct {
		ResourceName      string
		RefreshObjectName string
	}{
		ResourceName:      d.resourceName,
		RefreshObjectName: d.refreshObjectName,
	}

	err = readTemplate.ExecuteTemplate(&b, "Read_DataSource", data)
	if err != nil {
		log.Fatalf("error occurred with Generating Read: %v", err)
	}

	return b.Bytes()
}

// RenderRefresh implements BaseTemplate.
func (d *DataSourceTemplate) RenderRefresh() []byte {
	var b bytes.Buffer

	refreshTemplate, err := template.New("").Funcs(d.funcMap).Parse(RefreshTemplateDataSource)
	if err != nil {
		log.Fatalf("error occurred with baseTemplate at rendering refresh: %v", err)
	}

	data := struct {
		ResourceName      string
		RefreshObjectName string
		RefreshLogic      string
		ReadMethod        string
		Endpoint          string
		ReadPathParams    string
	}{
		ResourceName:      d.resourceName,
		RefreshObjectName: d.refreshObjectName,
		RefreshLogic:      d.refreshLogic,
		ReadMethod:        d.readMethod,
		Endpoint:          d.endpoint,
		ReadPathParams:    d.readPathParams,
	}

	err = refreshTemplate.ExecuteTemplate(&b, "Refresh_DataSource", data)
	if err != nil {
		log.Fatalf("error occurred with Generating Refresh: %v", err)
	}

	return b.Bytes()
}

// RenderTest implements BaseTemplate.
func (d *DataSourceTemplate) RenderTest() []byte {
	panic("Unimplemented yet.")
}

// RenderUpdate implements BaseTemplate.
func (d *DataSourceTemplate) RenderUpdate() []byte {
	panic("Data source doesn't provide this method. Please use the right method.")
}

// RenderWait implements BaseTemplate.
func (d *DataSourceTemplate) RenderWait() []byte {
	panic("Data source doesn't provide this method. Please use the right method.")
}

func NewDataSources(spec util.NcloudSpecification, datasourceName string) BaseTemplate {
	var b BaseTemplate
	var id string
	var attributes datasource.Attributes
	var readReqBody string
	var importStateOverride string
	var targetResourceRequest util.RequestWithRefreshObjectName

	d := &DataSourceTemplate{
		spec:         spec,
		resourceName: datasourceName,
	}

	funcMap := util.CreateFuncMap()

	for _, datasource := range spec.DataSources {
		if datasource.Name == datasourceName {
			id = datasource.Id
			attributes = datasource.Schema.Attributes
			importStateOverride = datasource.ImportStateOverride
		}
	}

	for _, val := range spec.Requests {
		if val.Name == datasourceName {
			targetResourceRequest = val
		}
	}

	_, model, err := Gen_ConvertOAStoTFTypes_Datasource(attributes)
	if err != nil {
		log.Fatalf("error occurred with Gen_ConvertOAStoTFTypes: %v", err)
	}

	for _, val := range targetResourceRequest.Read.Parameters {
		readReqBody = readReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueString(),`, util.PathToPascal(val), util.PathToPascal(val)) + "\n"
	}

	d.funcMap = funcMap
	d.providerName = spec.Provider.Name
	d.importStateLogic = MakeImportStateLogic(importStateOverride)
	d.model = model
	d.endpoint = spec.Provider.Endpoint
	d.readPathParams = extractReadPathParams(targetResourceRequest.Read.Path)
	d.readMethod = targetResourceRequest.Read.Method
	d.readReqBody = readReqBody
	d.readMethodName = strings.ToUpper(targetResourceRequest.Read.Method) + getMethodName(targetResourceRequest.Read.Path)
	d.idGetter = makeIdGetter(id)

	b = d

	return b
}
