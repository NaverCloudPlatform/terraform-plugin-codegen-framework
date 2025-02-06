package ncloud

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/util"
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

// To generate actual data, extract data from config.yml and code-spec.json, and render code for each receiver based on that data.
type BaseTemplate interface {

	// RenderInitial generates small code blocks needed initially.
	RenderInitial() []byte

	// RenderCreate generates the Create function.
	RenderCreate() []byte

	// RenderRead generates the Read function.
	RenderRead() []byte

	// RenderUpdate generates the Update function.
	RenderUpdate() []byte

	// RenderDelete generates the Delete function.
	RenderDelete() []byte

	// RenderModel generates the model.
	RenderModel() []byte

	// RenderRefresh generates the Refresh function.
	RenderRefresh() []byte

	// RenderWait generates the Waiting Logic.
	// Will be Rendered in refresh file.
	RenderWait() []byte

	// RenderTest generates the Test logic.
	RenderTest() []byte

	// RenderImportState generates the ImportState function.
	RenderImportState() []byte
}

type Template struct {
	spec                       util.NcloudSpecification
	providerName               string
	resourceName               string
	packageName                string
	importStateLogic           string
	refreshObjectName          string
	model                      string
	refreshLogic               string
	refreshWithResponse        string
	endpoint                   string
	deletePathParams           string
	updatePathParams           string
	readPathParams             string
	createPathParams           string
	deleteMethod               string
	updateMethod               string
	readMethod                 string
	createMethod               string
	createReqBody              string
	createReqListParams        string
	createReqObjectParams      string
	updateReqBody              string
	updateReqListParams        string
	updateReqObjectParams      string
	readReqBody                string
	deleteReqBody              string
	createMethodName           string
	readMethodName             string
	updateMethodName           string
	deleteMethodName           string
	idGetter                   string
	funcMap                    template.FuncMap
	configParams               string
	readReqBodyForCheckExist   string
	readReqBodyForCheckDestroy string
	isUpdateExists             bool
	createOpOptionalParams     string
	updateOpOptionalParams     string
	readOpOptionalParams       string
}

func (t *Template) RenderInitial() []byte {
	var b bytes.Buffer

	initialTemplate, err := template.New("").Funcs(t.funcMap).Parse(InitialTemplate)
	if err != nil {
		log.Fatalf("error occurred with baseTemplate at rendering initial: %v", err)
	}

	data := struct {
		ProviderName string
		ResourceName string
	}{
		ProviderName: t.providerName,
		ResourceName: t.resourceName,
	}

	err = initialTemplate.ExecuteTemplate(&b, "Initial", data)
	if err != nil {
		log.Fatalf("error occurred with generating Initial template: %v", err)
	}

	return b.Bytes()
}

func (t *Template) RenderImportState() []byte {
	var b bytes.Buffer

	initialTemplate, err := template.New("").Funcs(t.funcMap).Parse(ImportStateTemplate)
	if err != nil {
		log.Fatalf("error occurred with baseTemplate at rendering initial: %v", err)
	}

	data := struct {
		ResourceName     string
		ImportStateLogic string
	}{
		ResourceName:     t.resourceName,
		ImportStateLogic: t.importStateLogic,
	}

	err = initialTemplate.ExecuteTemplate(&b, "ImportState", data)
	if err != nil {
		log.Fatalf("error occurred with generating ImportState template: %v", err)
	}

	return b.Bytes()
}

func (t *Template) RenderCreate() []byte {
	var b bytes.Buffer

	createTemplate, err := template.New("").Funcs(t.funcMap).Parse(CreateTemplate)
	if err != nil {
		log.Fatalf("error occurred with baseTemplate at rendering create: %v", err)
	}

	data := struct {
		ResourceName           string
		RefreshObjectName      string
		CreateReqBody          string
		CreateReqListParam     string
		CreateReqObjectParam   string
		CreateReqOptionalParam string
		CreateMethod           string
		CreateMethodName       string
		Endpoint               string
		CreatePathParams       string
		IdGetter               string
	}{
		ResourceName:           t.resourceName,
		RefreshObjectName:      t.refreshObjectName,
		CreateReqBody:          t.createReqBody,
		CreateReqListParam:     t.createReqListParams,
		CreateReqObjectParam:   t.createReqObjectParams,
		CreateReqOptionalParam: t.createOpOptionalParams,
		CreateMethod:           t.createMethod,
		CreateMethodName:       t.createMethodName,
		Endpoint:               t.endpoint,
		CreatePathParams:       t.createPathParams,
		IdGetter:               t.idGetter,
	}

	err = createTemplate.ExecuteTemplate(&b, "Create", data)
	if err != nil {
		log.Fatalf("error occurred with Generating Create: %v", err)
	}

	return b.Bytes()
}

func (t *Template) RenderRead() []byte {
	var b bytes.Buffer

	readTemplate, err := template.New("").Funcs(t.funcMap).Parse(ReadTemplate)
	if err != nil {
		log.Fatalf("error occurred with baseTemplate at rendering read: %v", err)
	}

	data := struct {
		ResourceName      string
		RefreshObjectName string
	}{
		ResourceName:      t.resourceName,
		RefreshObjectName: t.refreshObjectName,
	}

	err = readTemplate.ExecuteTemplate(&b, "Read", data)
	if err != nil {
		log.Fatalf("error occurred with Generating Read: %v", err)
	}

	return b.Bytes()
}

func (t *Template) RenderUpdate() []byte {
	var b bytes.Buffer

	updateTemplate, err := template.New("").Funcs(t.funcMap).Parse(UpdateTemplate)
	if err != nil {
		log.Fatalf("error occurred with baseTemplate at rendering update: %v", err)
	}

	data := struct {
		IsUpdateExists         bool
		ResourceName           string
		RefreshObjectName      string
		UpdateReqBody          string
		UpdateReqListParam     string
		UpdateReqObjectParam   string
		UpdateReqOptionalParam string
		UpdateMethod           string
		UpdateMethodName       string
		Endpoint               string
		UpdatePathParams       string
		ReadPathParams         string
	}{
		IsUpdateExists:         t.isUpdateExists,
		ResourceName:           t.resourceName,
		RefreshObjectName:      t.refreshObjectName,
		UpdateReqBody:          t.updateReqBody,
		UpdateReqListParam:     t.updateReqListParams,
		UpdateReqObjectParam:   t.updateReqObjectParams,
		UpdateReqOptionalParam: t.updateOpOptionalParams,
		UpdateMethod:           t.updateMethod,
		UpdateMethodName:       t.updateMethodName,
		Endpoint:               t.endpoint,
		UpdatePathParams:       t.updatePathParams,
		ReadPathParams:         t.readPathParams,
	}

	err = updateTemplate.ExecuteTemplate(&b, "Update", data)
	if err != nil {
		log.Fatalf("error occurred with Generating Update: %v", err)
	}

	return b.Bytes()
}

func (t *Template) RenderDelete() []byte {
	var b bytes.Buffer

	deleteTemplate, err := template.New("").Funcs(t.funcMap).Parse(DeleteTemplate)
	if err != nil {
		log.Fatalf("error occurred with baseTemplate at rendering delete: %v", err)
	}

	data := struct {
		ResourceName      string
		RefreshObjectName string
		DeleteMethod      string
		DeleteReqBody     string
		DeleteMethodName  string
		Endpoint          string
		DeletePathParams  string
		IdGetter          string
	}{
		ResourceName:      t.resourceName,
		RefreshObjectName: t.refreshObjectName,
		DeleteMethod:      t.deleteMethod,
		DeleteReqBody:     t.deleteReqBody,
		DeleteMethodName:  t.deleteMethodName,
		Endpoint:          t.endpoint,
		DeletePathParams:  t.deletePathParams,
		IdGetter:          t.idGetter,
	}

	err = deleteTemplate.ExecuteTemplate(&b, "Delete", data)
	if err != nil {
		log.Fatalf("error occurred with Generating delete: %v", err)
	}

	return b.Bytes()
}

func (t *Template) RenderModel() []byte {
	var b bytes.Buffer

	modelTemplate, err := template.New("").Funcs(t.funcMap).Parse(ModelTemplate)
	if err != nil {
		log.Fatalf("error occurred with baseTemplate at rendering model: %v", err)
	}

	data := struct {
		RefreshObjectName string
		Model             string
	}{
		RefreshObjectName: t.refreshObjectName,
		Model:             t.model,
	}

	err = modelTemplate.ExecuteTemplate(&b, "Model", data)
	if err != nil {
		log.Fatalf("error occurred with Generating Model: %v", err)
	}

	return b.Bytes()
}

func (t *Template) RenderRefresh() []byte {
	var b bytes.Buffer

	refreshTemplate, err := template.New("").Funcs(t.funcMap).Parse(RefreshTemplate)
	if err != nil {
		log.Fatalf("error occurred with baseTemplate at rendering refresh: %v", err)
	}

	data := struct {
		PackageName         string
		RefreshObjectName   string
		RefreshWithResponse string
		Endpoint            string
		CreateMethodName    string
		ReadMethodName      string
		ReadReqBody         string
		IdGetter            string
	}{
		PackageName:         t.packageName,
		RefreshObjectName:   t.refreshObjectName,
		RefreshWithResponse: t.refreshWithResponse,
		Endpoint:            t.endpoint,
		CreateMethodName:    t.createMethodName,
		ReadMethodName:      t.readMethodName,
		ReadReqBody:         t.readReqBody,
		IdGetter:            t.idGetter,
	}

	err = refreshTemplate.ExecuteTemplate(&b, "Refresh", data)
	if err != nil {
		log.Fatalf("error occurred with Generating Refresh: %v", err)
	}

	return b.Bytes()
}

func (t *Template) RenderWait() []byte {
	var b bytes.Buffer

	waitTemplate, err := template.New("").Funcs(t.funcMap).Parse(WaitTemplate)
	if err != nil {
		log.Fatalf("error occurred with baseTemplate at rendering wait: %v", err)
	}

	data := struct {
		ReadMethod        string
		ReadMethodName    string
		Endpoint          string
		ReadPathParams    string
		RefreshObjectName string
		ReadReqBody       string
	}{
		ReadMethod:        t.readMethod,
		ReadMethodName:    t.readMethodName,
		Endpoint:          t.endpoint,
		ReadPathParams:    t.readPathParams,
		RefreshObjectName: t.refreshObjectName,
		ReadReqBody:       t.readReqBody,
	}

	err = waitTemplate.ExecuteTemplate(&b, "Wait", data)
	if err != nil {
		log.Fatalf("error occurred with Generating wait: %v", err)
	}

	return b.Bytes()
}

func (t *Template) RenderTest() []byte {
	var b bytes.Buffer

	testTemplate, err := template.New("").Funcs(t.funcMap).Parse(TestTemplate)
	if err != nil {
		log.Fatalf("error occurred with baseTemplate at rendering test: %v", err)
	}

	data := struct {
		ProviderName               string
		ResourceName               string
		PackageName                string
		RefreshObjectName          string
		ReadMethod                 string
		ReadMethodName             string
		ReadReqBody                string
		Endpoint                   string
		ReadPathParams             string
		ConfigParams               string
		ReadReqBodyForCheckExist   string
		ReadReqBodyForCheckDestroy string
	}{
		ProviderName:               t.providerName,
		ResourceName:               t.resourceName,
		PackageName:                t.packageName,
		RefreshObjectName:          t.refreshObjectName,
		ReadMethod:                 t.readMethod,
		ReadMethodName:             t.readMethodName,
		ReadReqBody:                t.readReqBody,
		Endpoint:                   t.endpoint,
		ReadPathParams:             t.readPathParams,
		ConfigParams:               t.configParams,
		ReadReqBodyForCheckExist:   t.readReqBodyForCheckExist,
		ReadReqBodyForCheckDestroy: t.readReqBodyForCheckDestroy,
	}

	err = testTemplate.ExecuteTemplate(&b, "Test", data)
	if err != nil {
		log.Fatalf("error occurred with Generating test: %v", err)
	}

	return b.Bytes()
}

type RequestType struct {
	Parameters  []string             `json:"parameters,omitempty"`
	RequestBody *OptionalRequestBody `json:"request_body,omitempty"`
	Response    string               `json:"response,omitempty"`
}

type OptionalRequestBody struct {
	Name     string   `json:"name,omitempty"`
	Required []string `json:"required,omitempty"`
	Optional []string `json:"optional,omitempty"`
}

// Extracts the data needed for code generation. Currently, it extracts data from config.yml and code-spec.json, but it is planned to unify everything into code-spec.json in the future.
func NewResource(spec util.NcloudSpecification, resourceName, packageName string) BaseTemplate {
	var b BaseTemplate
	var refreshObjectName string
	var id string
	var attributes resource.Attributes
	var createReqBody string
	var createReqListParams string
	var createReqObjectParams string
	var createOpOptionalParams string
	var updateReqBody string
	var updateReqListParams string
	var updateReqObjectParams string
	var updateOpOptionalParams string
	var readReqBody string
	var readReqBodyForCheckExist string
	var readReqBodyForCheckDestroy string
	var deleteReqBody string
	var readOpOptionalParams string
	var importStateOverride string
	var targetResourceRequest util.RequestInfo

	t := &Template{
		spec:         spec,
		resourceName: resourceName,
	}

	funcMap := util.CreateFuncMap()

	for _, resource := range spec.Resources {
		if resource.Name == resourceName {
			refreshObjectName = resource.RefreshObjectName
			id = resource.Id
			attributes = resource.Schema.Attributes
			importStateOverride = resource.ImportStateOverride
		}
	}

	for _, val := range spec.Requests {
		if val.Name == resourceName {
			targetResourceRequest = val
		}
	}

	refreshLogic, model, err := Gen_ConvertOAStoTFTypes_Resource(attributes)
	if err != nil {
		log.Fatalf("error occurred with Gen_ConvertOAStoTFTypes: %v", err)
	}

	// Address Request > Create
	if targetResourceRequest.Create != nil {
		// Address Request > Create > RequestBody
		if targetResourceRequest.Create.RequestBody != nil {

			if targetResourceRequest.Create.RequestBody.Required != nil {
				for _, val := range targetResourceRequest.Create.RequestBody.Required {

					switch val.Type {
					case "string":
						createReqBody = createReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueString(),`, util.FirstAlphabetToUpperCase(val.Name), util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					case "integer":
						switch val.Format {
						case "int64":
							createReqBody = createReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueInt64(),`, util.FirstAlphabetToUpperCase(val.Name), util.FirstAlphabetToUpperCase(val.Name)) + "\n"

						case "int32":
							createReqBody = createReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueInt32(),`, util.FirstAlphabetToUpperCase(val.Name), util.FirstAlphabetToUpperCase(val.Name)) + "\n"
						}

					case "number":
						createReqBody = createReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueFloat64(),`, util.FirstAlphabetToUpperCase(val.Name), util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					case "boolean":
						createReqBody = createReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueBool(),`, util.FirstAlphabetToUpperCase(val.Name), util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					case "array":
						createReqListParams = createReqListParams + fmt.Sprintf(`
					list%[1]s, diags := types.ListValue(
						plan.%[1]s.ElementType(ctx),
						plan.%[1]s.Elements(),
					)
					resp.Diagnostics.Append(diags...)
					if resp.Diagnostics.HasError() {
						return
					}
	
					reqParams.%[1]s = list%[1]s
					}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					case "object":
						createReqObjectParams = createReqObjectParams + fmt.Sprintf(`
					obj%[1]s, diags := types.ObjectValue(
						plan.%[1]s.AttributeTypes(ctx),
						plan.%[1]s.Attributes(),
					)
					resp.Diagnostics.Append(diags...)
					if resp.Diagnostics.HasError() {
						return
					}
	
					reqParams.%[1]s= obj%[1]s
					`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"
					}

				}
			}

			if targetResourceRequest.Create.RequestBody.Optional != nil {
				for _, val := range targetResourceRequest.Create.RequestBody.Optional {

					switch val.Type {
					case "string":
						createOpOptionalParams = createOpOptionalParams + fmt.Sprintf(`
					if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
						reqParams.%[1]s = plan.%[1]s.ValueString()
					}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					case "integer":

						switch val.Format {
						case "int64":
							createOpOptionalParams = createOpOptionalParams + fmt.Sprintf(`
							if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
								reqParams.%[1]s = plan.%[1]s.ValueInt64()
							}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"

						case "int32":
							createOpOptionalParams = createOpOptionalParams + fmt.Sprintf(`
							if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
								reqParams.%[1]s = plan.%[1]s.ValueInt32()
							}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"
						}

					case "number":
						createOpOptionalParams = createOpOptionalParams + fmt.Sprintf(`
					if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
						reqParams.%[1]s = plan.%[1]s.ValueFloat64()
					}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					case "boolean":
						createOpOptionalParams = createOpOptionalParams + fmt.Sprintf(`
					if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
						reqParams.%[1]s = plan.%[1]s.ValueBool()
					}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					case "array":
						createOpOptionalParams = createOpOptionalParams + fmt.Sprintf(`
					if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
						list%[1]s, diags := types.ListValue(
							plan.%[1]s.ElementType(ctx),
							plan.%[1]s.Elements(),
						)
						resp.Diagnostics.Append(diags...)
						if resp.Diagnostics.HasError() {
							return
						}

						reqParams.%[1]s = list%[1]s
					}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					case "object":
						createOpOptionalParams = createOpOptionalParams + fmt.Sprintf(`
					if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
						obj%[1]s, diags := types.ObjectValue(
							plan.%[1]s.AttributeTypes(ctx),
							plan.%[1]s.Attributes(),
						)
						resp.Diagnostics.Append(diags...)
						if resp.Diagnostics.HasError() {
							return
						}

						reqParams.%[1]s= obj%[1]s
					}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					// Array and Object are treated as string with serialization
					default:
						createOpOptionalParams = createOpOptionalParams + fmt.Sprintf(`
					if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
						reqParams.%[1]s = plan.%[1]s.ValueString()
					}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"
					}
				}
			}

			if targetResourceRequest.Create.Parameters != nil {
				t.configParams = MakeTestTFConfig(targetResourceRequest.Create)
			}

			t.createPathParams = extractPathParams(targetResourceRequest.Create.Path)
			t.createMethod = targetResourceRequest.Create.Method
			t.createMethodName = strings.ToUpper(targetResourceRequest.Create.Method) + getMethodName(targetResourceRequest.Create.Path)
		}

		if targetResourceRequest.Create.Parameters != nil {

			if targetResourceRequest.Create.Parameters.Required != nil {
				for _, val := range targetResourceRequest.Create.Parameters.Required {
					switch val.Type {
					case "string":
						createReqBody = createReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueString(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"

					case "integer":
						if val.Format != "" {
							switch val.Format {
							case "int64":
								createReqBody = createReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueInt64(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"

							case "int32":
								createReqBody = createReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueInt32(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"
							}
						}

					case "number":
						createReqBody = createReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueFloat64(,`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"

					case "boolean":
						createReqBody = createReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueBool(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"
					}
				}
			}

			if targetResourceRequest.Create.Parameters.Optional != nil {
				// NOTE - CREATE does not have optional parameters
				panic("CREATE does not have optional Parameters in the previous cases. Please notify the developer to implement.")
			}

		}

		t.createReqListParams = createReqListParams
		t.createReqObjectParams = createReqObjectParams
	}

	// Address Request > Read
	if targetResourceRequest.Read != nil {
		// Address Request > Read > Parameters (READ does not have RequestBody)
		if targetResourceRequest.Read.Parameters != nil {

			for _, val := range targetResourceRequest.Read.Parameters.Required {

				switch val.Type {
				case "string":
					readReqBody = readReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueString(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"

				case "integer":
					if val.Format != "" {
						switch val.Format {
						case "int64":
							readReqBody = readReqBody + fmt.Sprintf(`%[1]s: int(plan.%[2]s.ValueInt64(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"

						case "int32":
							readReqBody = readReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueInt32(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"
						}
					}

				case "number":
					readReqBody = readReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueFloat64(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"

				case "boolean":
					readReqBody = readReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueBool(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"
				}

				readReqBodyForCheckExist = readReqBodyForCheckExist + fmt.Sprintf(`		%[1]s: resource.Primary.Attributes["%[2]s"],`, util.PathToPascal(val.Name), util.FirstAlphabetToLowerCase(util.PathToPascal(val.Name))) + "\n"
				readReqBodyForCheckDestroy = readReqBodyForCheckDestroy + fmt.Sprintf(`		%[1]s: rs.Primary.Attributes["%[2]s"],`, util.PathToPascal(val.Name), util.FirstAlphabetToLowerCase(util.PathToPascal(val.Name))) + "\n"
			}

			if targetResourceRequest.Read.Parameters.Optional != nil {
				for _, val := range targetResourceRequest.Read.Parameters.Optional {

					switch val.Type {
					case "string":
						readOpOptionalParams = readOpOptionalParams + fmt.Sprintf(`
					if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
						reqParams.%[1]s = plan.%[1]s.ValueString()
					}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					case "integer":
						if val.Format != "" {
							switch val.Format {
							case "int64":
								readOpOptionalParams = readOpOptionalParams + fmt.Sprintf(`
							if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
								reqParams.%[1]s = plan.%[1]s.ValueInt64()
							}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"

							case "int32":
								readOpOptionalParams = readOpOptionalParams + fmt.Sprintf(`
							if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
								reqParams.%[1]s = plan.%[1]s.ValueInt32()
							}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"
							}
						}

					case "number":
						readOpOptionalParams = readOpOptionalParams + fmt.Sprintf(`
					if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
						reqParams.%[1]s = plan.%[1]s.ValueFloat64()
					}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					case "boolean":
						readOpOptionalParams = readOpOptionalParams + fmt.Sprintf(`
					if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
						reqParams.%[1]s = plan.%[1]s.ValueBool()
					}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					case "array":
						readOpOptionalParams = readOpOptionalParams + fmt.Sprintf(`
					if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
						reqParams.%[1]s = plan.%[1]s.ValueString()
					}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"

						// Array and Object are treated as string with serialization
					default:
						readOpOptionalParams = readOpOptionalParams + fmt.Sprintf(`
					if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
						reqParams.%[1]s = plan.%[1]s.ValueString()
					}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"
					}
				}
			}
		}

		t.readPathParams = extractReadPathParams(targetResourceRequest.Read.Path)
		t.readMethod = targetResourceRequest.Read.Method
		t.readMethodName = strings.ToUpper(targetResourceRequest.Read.Method) + getMethodName(targetResourceRequest.Read.Path)
	}

	// Address Request > Update
	if len(targetResourceRequest.Update) > 0 {

		// Address Request > Update > Parameters
		if targetResourceRequest.Update[0].Parameters != nil {

			if targetResourceRequest.Update[0].Parameters.Required != nil {
				for _, val := range targetResourceRequest.Update[0].Parameters.Required {
					switch val.Type {
					case "string":
						updateReqBody = updateReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueString(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"

					case "integer":
						if val.Format != "" {
							switch val.Format {
							case "int64":
								updateReqBody = updateReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueInt64(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"

							case "int32":
								updateReqBody = updateReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueInt32(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"
							}
						}

					case "number":
						updateReqBody = updateReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueFloat64()),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"

					case "boolean":
						updateReqBody = updateReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueBool(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"

					case "array":
						updateReqListParams = updateReqListParams + fmt.Sprintf(`
						list%[1]s, diags := types.ListValue(
							plan.%[1]s.ElementType(ctx),
							plan.%[1]s.Elements(),
						)
						resp.Diagnostics.Append(diags...)
						if resp.Diagnostics.HasError() {
							return
						}
		
						reqParams.%[1]s = list%[1]s
						}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					case "object":
						updateReqObjectParams = updateReqObjectParams + fmt.Sprintf(`
						obj%[1]s, diags := types.ObjectValue(
							plan.%[1]s.AttributeTypes(ctx),
							plan.%[1]s.Attributes(),
						)
						resp.Diagnostics.Append(diags...)
						if resp.Diagnostics.HasError() {
							return
						}
		
						reqParams.%[1]s= obj%[1]s
						`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"
					}
				}
			}

			if targetResourceRequest.Update[0].Parameters.Optional != nil {
				// NOTE - UPDATE does not have optional parameters
				panic("UPDATE does not have optional Parameters in the previous cases. Please notify the developer to implement.")
			}

		}

		if targetResourceRequest.Update[0].RequestBody != nil {

			if targetResourceRequest.Update[0].RequestBody.Required != nil {
				for _, val := range targetResourceRequest.Update[0].RequestBody.Required {
					switch val.Type {
					case "string":
						updateReqBody = updateReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueString(),`, util.FirstAlphabetToUpperCase(val.Name), util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					case "integer":
						if val.Format != "" {
							switch val.Format {
							case "int64":
								updateReqBody = updateReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueInt64(),`, util.FirstAlphabetToUpperCase(val.Name), util.FirstAlphabetToUpperCase(val.Name)) + "\n"

							case "int32":
								updateReqBody = updateReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueInt32(),`, util.FirstAlphabetToUpperCase(val.Name), util.FirstAlphabetToUpperCase(val.Name)) + "\n"
							}
						}

					case "number":
						updateReqBody = updateReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueFloat64(),`, util.FirstAlphabetToUpperCase(val.Name), util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					case "boolean":
						updateReqBody = updateReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueBool(),`, util.FirstAlphabetToUpperCase(val.Name), util.FirstAlphabetToUpperCase(val.Name)) + "\n"
					}
				}
			}

			if targetResourceRequest.Update[0].RequestBody.Optional != nil {
				for _, val := range targetResourceRequest.Update[0].RequestBody.Optional {

					switch val.Type {
					case "string":
						updateOpOptionalParams = updateOpOptionalParams + fmt.Sprintf(`
						if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
							reqParams.%[1]s = plan.%[1]s.ValueString()
						}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					case "integer":
						if val.Format != "" {
							switch val.Format {
							case "int64":
								updateOpOptionalParams = updateOpOptionalParams + fmt.Sprintf(`
								if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
									reqParams.%[1]s = plan.%[1]s.ValueInt64()
								}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"

							case "int32":
								updateOpOptionalParams = updateOpOptionalParams + fmt.Sprintf(`
								if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
									reqParams.%[1]s = plan.%[1]s.ValueInt32()
								}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"
							}
						}

					case "number":
						updateOpOptionalParams = updateOpOptionalParams + fmt.Sprintf(`
								if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
									reqParams.%[1]s = plan.%[1]s.ValueFloat64()
								}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					case "boolean":
						updateOpOptionalParams = updateOpOptionalParams + fmt.Sprintf(`
						if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
							reqParams.%[1]s = plan.%[1]s.ValueBool()
						}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					case "array":
						updateOpOptionalParams = updateOpOptionalParams + fmt.Sprintf(`
						if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
							reqParams.%[1]s = plan.%[1]s.ValueString()
						}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"

					// Array and Object are treated as string with serialization
					default:
						updateOpOptionalParams = updateOpOptionalParams + fmt.Sprintf(`
						if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
							reqParams.%[1]s = plan.%[1]s.ValueString()
						}`, util.FirstAlphabetToUpperCase(val.Name)) + "\n"
					}
				}
			}

		}

		t.updateReqListParams = updateReqListParams
		t.updateReqObjectParams = updateReqObjectParams
		t.isUpdateExists = true
		t.updatePathParams = extractPathParams(targetResourceRequest.Update[0].Path)
		t.updateMethod = targetResourceRequest.Update[0].Method
		t.updateReqBody = updateReqBody
		t.updateOpOptionalParams = updateOpOptionalParams
		t.updateMethodName = strings.ToUpper(targetResourceRequest.Update[0].Method) + getMethodName(targetResourceRequest.Update[0].Path)
	}

	// Address Request > Delete
	if targetResourceRequest.Delete != nil {
		if targetResourceRequest.Delete.Parameters != nil {

			if targetResourceRequest.Delete.Parameters.Required != nil {
				for _, val := range targetResourceRequest.Delete.Parameters.Required {
					switch val.Type {
					case "string":
						deleteReqBody = deleteReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueString(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"

					case "integer":
						if val.Format != "" {
							switch val.Format {
							case "int64":
								deleteReqBody = deleteReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueInt64(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"

							case "int32":
								deleteReqBody = deleteReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueInt32(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"
							}
						}

					case "number":
						deleteReqBody = deleteReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueFloat64(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"

					case "boolean":
						deleteReqBody = deleteReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueBool(),`, util.PathToPascal(val.Name), util.PathToPascal(val.Name)) + "\n"
					}
				}
			}

			if targetResourceRequest.Delete.Parameters.Optional != nil {
				// NOTE - DELETE does not have optional parameters
				panic("DELETE does not have optional Parameters in the previous cases. Please notify the developer to implement.")
			}

		}

		if targetResourceRequest.Delete.RequestBody != nil {
			// NOTE - DELETE does not have RequestBody
			panic("DELETE does not have RequestBody in the previous cases. Please notify the developer to implement.")
		}

		t.deletePathParams = extractPathParams(targetResourceRequest.Delete.Path)
		t.deleteMethod = targetResourceRequest.Delete.Method
		t.deleteMethodName = strings.ToUpper(targetResourceRequest.Delete.Method) + getMethodName(targetResourceRequest.Delete.Path)
	}

	t.funcMap = funcMap
	t.providerName = spec.Provider.Name
	t.packageName = packageName
	t.refreshObjectName = refreshObjectName
	t.importStateLogic = makeImportStateLogic(importStateOverride)
	t.model = model
	t.refreshLogic = refreshLogic
	t.refreshWithResponse = MakeRefreshFromResponse(attributes, resourceName)
	t.endpoint = spec.Provider.Endpoint
	t.createReqBody = createReqBody
	t.readReqBody = readReqBody
	t.deleteReqBody = deleteReqBody
	t.readReqBodyForCheckExist = readReqBodyForCheckExist
	t.readReqBodyForCheckDestroy = readReqBodyForCheckDestroy
	t.createOpOptionalParams = createOpOptionalParams
	t.readOpOptionalParams = readOpOptionalParams
	t.idGetter = makeIdGetter(id)

	b = t

	return b
}

func getMethodName(s string) string {
	parts := strings.Split(s, "/")
	var result []string

	for _, part := range parts {
		if part == "" {
			continue
		}

		// Remove curly braces if present
		part = strings.TrimPrefix(part, "{")
		part = strings.TrimSuffix(part, "}")

		// Remove hyphens and convert to uppercase
		part = strings.ReplaceAll(part, "-", "")
		part = util.FirstAlphabetToUpperCase(part)

		result = append(result, part)
	}

	return strings.Join(result, "")
}

func extractPathParams(path string) string {
	parts := strings.Split(path, "/")
	s := ``

	for idx, val := range parts {

		if len(val) < 1 {
			continue
		}

		s = s + `+"/"+`

		start := strings.Index(val, "{")

		// if val doesn't wrapped with curly brace
		if start == -1 {
			s = s + fmt.Sprintf(`"%s"`, val)
		} else {
			if idx == len(parts)-1 {
				s = s + `plan.ID.ValueString()`
			} else {
				s = s + fmt.Sprintf(`plan.%s.ValueString())`, util.PathToPascal(val))
			}
		}
	}

	return s
}

func extractReadPathParams(path string) string {
	parts := strings.Split(path, "/")
	s := ``

	for idx, val := range parts {
		if len(val) < 1 {
			continue
		}

		if idx == len(parts)-1 {
			continue
		}

		s = s + `+"/"+`

		start := strings.Index(val, "{")

		// if val doesn't wrapped with curly brace
		if start == -1 {
			s = s + fmt.Sprintf(`"%s"`, val)
		} else {
			s = s + fmt.Sprintf(`plan.%s.ValueString()`, util.PathToPascal(val))
		}
	}

	return s
}

func makeIdGetter(target string) string {
	parts := strings.Split(target, ".")
	s := "createRes"

	for idx, val := range parts {
		if idx == len(parts)-1 {
			s = s + fmt.Sprintf(`["%s"].(string)`, util.ToCamelCase(val))
			continue
		}

		s = s + fmt.Sprintf(`["%s"].(map[string]interface{})`, util.ToCamelCase(val))
	}

	return s
}

func makeImportStateLogic(target string) string {
	parts := strings.Split(target, ".")

	if len(parts) < 2 {
		return `resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)` + "\n"
	}

	s := `parts := strings.Split(req.ID, ".")` + "\n"
	for idx, val := range parts {
		s = s + fmt.Sprintf(`resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("%s"), parts[%d])...)`, util.ToLowerCase(util.PathToPascal(val)), idx) + "\n"
	}

	return s
}

func MakeRefreshFromResponse(attr resource.Attributes, resourceName string) string {
	var s strings.Builder

	for _, val := range attr {
		if val.String != nil {
			s.WriteString(fmt.Sprintf(`
			if !response.%[2]s.Attributes()["%[1]s"].IsNull() || !response.%[2]s.Attributes()["%[1]s"].IsUnknown() {
				postPlan.%[3]s = types.StringValue(response.%[2]s.Attributes()["%[1]s"].String())
			}`, val.Name, util.ToPascalCase(resourceName), util.ToPascalCase(val.Name)) + "\n")
		} else if val.Bool != nil {
			s.WriteString(fmt.Sprintf(`
			if !response.%[2]s.Attributes()["%[1]s"].IsNull() || !response.%[2]s.Attributes()["%[1]s"].IsUnknown() {
				tempVal, err := strconv.ParseBool(response.%[2]s.Attributes()["%[1]s"].String())
				if err != nil {
					diagnostics.AddError("CONVERSION ERROR", fmt.Sprintf("Failed to convert %[1]s to bool: %%v", err))
					return
				}
				postPlan.%[3]s = types.BoolValue(tempVal)
			}`, val.Name, util.ToPascalCase(resourceName), util.ToPascalCase(val.Name)) + "\n")
		} else if val.Int32 != nil {
			s.WriteString(fmt.Sprintf(`
			if !response.%[2]s.Attributes()["%[1]s"].IsNull() || !response.%[2]s.Attributes()["%[1]s"].IsUnknown() {
				tempVal, err := strconv.Atoi(response.%[2]s.Attributes()["%[1]s"].String())
				if err != nil {
					diagnostics.AddError("CONVERSION ERROR", fmt.Sprintf("Failed to convert %[1]s to int: %%v", err))
					return
				}
				postPlan.%[3]s = types.Int32Value(int32(tempVal))
			}`, val.Name, util.ToPascalCase(resourceName), util.ToPascalCase(val.Name)) + "\n")
		} else if val.Int64 != nil {
			s.WriteString(fmt.Sprintf(`
			if !response.%[2]s.Attributes()["%[1]s"].IsNull() || !response.%[2]s.Attributes()["%[1]s"].IsUnknown() {
				tempVal, err := strconv.Atoi(response.%[2]s.Attributes()["%[1]s"].String())
				if err != nil {
					diagnostics.AddError("CONVERSION ERROR", fmt.Sprintf("Failed to convert %[1]s to int: %%v", err))
					return
				}
				postPlan.%[3]s = types.Int64Value(int64(tempVal))
			}`, val.Name, util.ToPascalCase(resourceName), util.ToPascalCase(val.Name)) + "\n")
		} else if val.Float64 != nil || val.Float32 != nil {
			s.WriteString(fmt.Sprintf(`
			if !response.%[2]s.Attributes()["%[1]s"].IsNull() || !response.%[2]s.Attributes()["%[1]s"].IsUnknown() {
				tempVal, err := strconv.ParseFloat(response.%[2]s.Attributes()["%[1]s"].String(), 64)
				if err != nil {
					diagnostics.AddError("CONVERSION ERROR", fmt.Sprintf("Failed to convert %[1]s to float64: %%v", err))
					return
				}
				postPlan.%[3]s = types.Float64Value(tempVal)
			}`, val.Name, util.ToPascalCase(resourceName), util.ToPascalCase(val.Name)) + "\n")
		} else if val.List != nil {
			s.WriteString(fmt.Sprintf(`
			if !response.%[2]s.Attributes()["%[1]s"].IsNull() || !response.%[2]s.Attributes()["%[1]s"].IsUnknown() {
				listRes, diag := types.ListValueFrom(ctx, postPlan.%[2]s.AttributeTypes(ctx), response.%[2]s)
				if diag.HasError() {
					diagnostics.AddError("CONVERSION ERROR", "Error occured while getting object value: %[1]s")
					return
				}
				postPlan.%[3]s = listRes
			}`, val.Name, util.ToPascalCase(resourceName), util.ToPascalCase(val.Name)) + "\n")
		} else if val.ListNested != nil {
			s.WriteString(fmt.Sprintf(`
			if !response.%[2]s.Attributes()["%[1]s"].IsNull() || !response.%[2]s.Attributes()["%[1]s"].IsUnknown() {
				listRes, diag := types.ListValueFrom(ctx, postPlan.%[2]s.AttributeTypes(ctx), response.%[2]s)
				if diag.HasError() {
					diagnostics.AddError("CONVERSION ERROR", "Error occured while getting object value: %[1]s")
					return
				}
				postPlan.%[3]s = listRes
			}`, val.Name, util.ToPascalCase(resourceName), util.ToPascalCase(val.Name)) + "\n")
		} else if val.SingleNested != nil {
			s.WriteString(fmt.Sprintf(`
			if !response.%[2]s.Attributes()["%[1]s"].IsNull() || !response.%[2]s.Attributes()["%[1]s"].IsUnknown() {
				objectRes, diag := types.ObjectValueFrom(ctx, postPlan.%[2]s.AttributeTypes(ctx), response.%[2]s)
				if diag.HasError() {
					diagnostics.AddError("CONVERSION ERROR", "Error occured while getting object value: %[1]s")
					return
				}
				postPlan.%[3]s = objectRes
			}`, val.Name, util.ToPascalCase(resourceName), util.ToPascalCase(val.Name)) + "\n")
		}
	}

	return s.String()
}

func MakeTestTFConfig(c *util.NcloudCommonRequestType) string {
	var t strings.Builder

	if c.RequestBody != nil {
		for _, val := range c.RequestBody.Required {
			t.WriteString(fmt.Sprintf(`		%[1]s = "%[2]s"`, util.FirstAlphabetToLowerCase(util.PathToPascal(val.Name)), "tf-"+acctest.RandString(5)) + "\n")
		}
	}

	if c.Parameters != nil {
		for _, val := range c.Parameters.Required {
			t.WriteString(fmt.Sprintf(`		%[1]s = "%[2]s"`, util.FirstAlphabetToLowerCase(util.PathToPascal(val.Name)), "tf-"+acctest.RandString(5)) + "\n")
		}
	}

	return t.String()
}
