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
	spec                 util.NcloudSpecification
	providerName         string
	dataSourceName       string
	packageName          string
	importStateLogic     string
	refreshObjectName    string
	model                string
	refreshLogic         string
	endpoint             string
	readPathParams       string
	readMethod           string
	readReqBody          string
	readMethodName       string
	readOpOptionalParams string
	idGetter             string
	funcMap              template.FuncMap
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
		ProviderName   string
		DataSourceName string
	}{
		ProviderName:   d.providerName,
		DataSourceName: d.dataSourceName,
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
		DataSourceName    string
		RefreshObjectName string
	}{
		DataSourceName:    d.dataSourceName,
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
		PackageName          string
		ResourceName         string
		RefreshObjectName    string
		RefreshLogic         string
		ReadMethodName       string
		ReadReqBody          string
		Endpoint             string
		ReadPathParams       string
		ReadOpOptionalParams string
	}{
		PackageName:          d.packageName,
		ResourceName:         d.dataSourceName,
		RefreshObjectName:    d.refreshObjectName,
		RefreshLogic:         d.refreshLogic,
		ReadMethodName:       d.readMethodName,
		ReadReqBody:          d.readReqBody,
		Endpoint:             d.endpoint,
		ReadPathParams:       d.readPathParams,
		ReadOpOptionalParams: d.readOpOptionalParams,
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

func NewDataSources(spec *util.NcloudSpecification, datasourceName, packageName string) BaseTemplate {
	var b BaseTemplate
	var targetResourceRequest util.RequestWithRefreshObjectName

	d := &DataSourceTemplate{
		spec:           *spec,
		dataSourceName: datasourceName,
		providerName:   spec.Provider.Name,
		packageName:    packageName,
		endpoint:       spec.Provider.Endpoint,
	}

	d.funcMap = util.CreateFuncMap()

	for _, val := range spec.Requests {
		if val.Name == datasourceName {
			targetResourceRequest = val
		}
	}

	if err := makeDataSourceIndividualValues(d, spec, datasourceName); err != nil {
		log.Fatalf("error occurred with MakeDataSourceIndividualValues: %v", err)
	}

	if err := makeDataSourceReadOperationLogics(d, &targetResourceRequest); err != nil {
		log.Fatalf("error occurred with MakeDataSourceReadOperationLogics: %v", err)
	}

	b = d
	return b
}

func makeDataSourceReadOperationLogics(d *DataSourceTemplate, t *util.RequestWithRefreshObjectName) error {
	var readOpOptionalParams strings.Builder
	var readReqBody strings.Builder

	for _, val := range t.Read.Parameters.Required {
		if _, err := readReqBody.WriteString(fmt.Sprintf(`%[1]s: plan.%[2]s.ValueString(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"); err != nil {
			return err
		}
	}

	for _, val := range t.Read.Parameters.Optional {

		switch val.Type {

		case "string":
			if _, err := readOpOptionalParams.WriteString(fmt.Sprintf(`
			if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
				r.%[1]s = plan.%[1]s.ValueString()
			}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"); err != nil {
				return err
			}

		case "integer":
			if _, err := readOpOptionalParams.WriteString(fmt.Sprintf(`
			if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
				r.%[1]s = strconv.Itoa(int(plan.%[1]s.ValueInt64()))
			}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"); err != nil {
				return err
			}

		case "boolean":
			if _, err := readOpOptionalParams.WriteString(fmt.Sprintf(`
			if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
				r.%[1]s = strconv.FormatBool(plan.%[1]s.ValueBool())
			}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"); err != nil {
				return err
			}

		case "array":
			if _, err := readOpOptionalParams.WriteString(fmt.Sprintf(`
			if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
				r.%[1]s = plan.%[1]s.ValueString()
			}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"); err != nil {
				return err
			}

		// Array and Object are treated as string with serialization
		default:
			if _, err := readOpOptionalParams.WriteString(fmt.Sprintf(`
			if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
				reqParams.%[1]s = plan.%[1]s.ValueString()
			}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"); err != nil {
				return err
			}
		}
	}

	d.readOpOptionalParams = readOpOptionalParams.String()
	d.readReqBody = readReqBody.String()
	d.readPathParams = extractReadPathParams(t.Read.Path)
	d.readMethod = t.Read.Method
	d.readMethodName = strings.ToUpper(t.Read.Method) + getMethodName(t.Read.Path)

	return nil
}

func makeDataSourceIndividualValues(d *DataSourceTemplate, spec *util.NcloudSpecification, datasourceName string) error {
	var attributes datasource.Attributes

	for _, datasource := range spec.DataSources {
		if datasource.Name == datasourceName {
			d.idGetter = makeIdGetter(datasource.Id)
			d.refreshObjectName = datasource.RefreshObjectName
			attributes = datasource.Schema.Attributes
			d.importStateLogic = makeImportStateLogic(datasource.ImportStateOverride)
		}
	}

	_, model, err := Gen_ConvertOAStoTFTypes_Datasource(attributes)
	if err != nil {
		return fmt.Errorf("error occurred with Gen_ConvertOAStoTFTypes: %v", err)
	}

	d.model = model
	return nil
}
