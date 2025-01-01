package ncloud

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

// To generate actual data, extract data from config.yml and code-spec.json, and render code for each receiver based on that data.
// New(): Extracts the data Required for code generation. Currently, it extracts data from config.yml and code-spec.json, but it is planned to unify everything into code-spec.json in the future.
// RenderInitial(): Generates small code blocks Required initially.
// RenderCreate(): Generates the Create function.
// RenderRead(): Generates the Read function.
// RenderUpdate(): Generates the Update function.
// RenderDelete(): Generates the Delete function.
// Calculates the necessary data during initialization and performs rendering for each method.
type Template struct {
	spec                  util.NcloudSpecification
	providerName          string
	resourceName          string
	packageName           string
	importStateLogic      string
	refreshObjectName     string
	model                 string
	refreshLogic          string
	refreshWithResponse   string
	endpoint              string
	deletePathParams      string
	updatePathParams      string
	readPathParams        string
	createPathParams      string
	deleteMethod          string
	updateMethod          string
	readMethod            string
	createMethod          string
	createReqBody         string
	updateReqBody         string
	readReqBody           string
	deleteReqBody         string
	createMethodName      string
	readMethodName        string
	updateMethodName      string
	deleteMethodName      string
	idGetter              string
	funcMap               template.FuncMap
	createOpOptionalParam string
	updateOpOptionalParam string
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
		CreateReqOptionalParam: t.createOpOptionalParam,
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
		ResourceName           string
		RefreshObjectName      string
		UpdateReqBody          string
		UpdateReqOptioanlParam string
		UpdateMethod           string
		UpdateMethodName       string
		Endpoint               string
		UpdatePathParams       string
		ReadPathParams         string
	}{
		ResourceName:           t.resourceName,
		RefreshObjectName:      t.refreshObjectName,
		UpdateReqBody:          t.updateReqBody,
		UpdateReqOptioanlParam: t.updateOpOptionalParam,
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
		ProviderName      string
		ResourceName      string
		PackageName       string
		RefreshObjectName string
		ReadMethod        string
		ReadMethodName    string
		ReadReqBody       string
		Endpoint          string
		ReadPathParams    string
	}{
		ProviderName:      t.providerName,
		ResourceName:      t.resourceName,
		PackageName:       t.packageName,
		RefreshObjectName: t.refreshObjectName,
		ReadMethod:        t.readMethod,
		ReadMethodName:    t.readMethodName,
		ReadReqBody:       t.readReqBody,
		Endpoint:          t.endpoint,
		ReadPathParams:    t.readPathParams,
	}

	err = testTemplate.ExecuteTemplate(&b, "Test", data)
	if err != nil {
		log.Fatalf("error occurred with Generating test: %v", err)
	}

	return b.Bytes()
}

type RequestType struct {
	Parameters  []string                 `json:"parameters,omitempty"`
	RequestBody *RequestBodyWithOptional `json:"request_body,omitempty"`
	Response    string                   `json:"response,omitempty"`
}

type RequestBodyWithOptional struct {
	Name     string   `json:"name,omitempty"`
	Required []string `json:"required,omitempty"`
	Optional []string `json:"optional,omitempty"`
}

// Allocate prerequisite values
func New(spec util.NcloudSpecification, resourceName, packageName string) *Template {
	var refreshObjectName string
	var id string
	var attributes resource.Attributes
	var createReqBody string
	var updateReqBody string
	var readReqBody string
	var deleteReqBody string
	var createOpOptionalParam string
	var updateOpOptionalParam string
	var importStateOverride string
	var targetResourceRequest util.RequestWithRefreshObjectName

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

	refreshLogic, model, err := Gen_ConvertOAStoTFTypes(attributes)
	if err != nil {
		log.Fatalf("error occurred with Gen_ConvertOAStoTFTypes: %v", err)
	}

	for _, val := range targetResourceRequest.Create.RequestBody.Required {
		createReqBody = createReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueString(),`, util.FirstAlphabetToUpperCase(val), util.FirstAlphabetToUpperCase(val)) + "\n"
	}

	for _, val := range targetResourceRequest.Update[0].RequestBody.Required {
		updateReqBody = updateReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueString(),`, util.FirstAlphabetToUpperCase(val), util.FirstAlphabetToUpperCase(val)) + "\n"
	}

	for _, val := range targetResourceRequest.Read.Parameters {
		readReqBody = readReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueString(),`, util.PathToPascal(val), util.PathToPascal(val)) + "\n"
	}

	for _, val := range targetResourceRequest.Create.Parameters {
		createReqBody = createReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueString(),`, util.PathToPascal(val), util.PathToPascal(val)) + "\n"
	}

	for _, val := range targetResourceRequest.Update[0].Parameters {
		updateReqBody = updateReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueString(),`, util.PathToPascal(val), util.PathToPascal(val)) + "\n"
	}

	for _, val := range targetResourceRequest.Delete.Parameters {
		deleteReqBody = deleteReqBody + fmt.Sprintf(`%[1]s: plan.%[2]s.ValueString(),`, util.PathToPascal(val), util.PathToPascal(val)) + "\n"
	}

	for _, val := range targetResourceRequest.Create.RequestBody.Optional {

		// TODO - Check types and cast it to right type
		// in case of int64 : use strconv.Itoa(int())
		// in case of bool : use strconv.FormatBool()
		createOpOptionalParam = createOpOptionalParam + fmt.Sprintf(`
		if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
			reqParams.%[1]s = plan.%[1]s.ValueString()
		},`, util.ToPascalCase(val)) + "\n"
	}

	for _, val := range targetResourceRequest.Update[0].RequestBody.Optional {

		// TODO - Check types and cast it to right type
		// in case of int64 : use strconv.Itoa(int())
		// in case of bool : use strconv.FormatBool()
		updateOpOptionalParam = updateOpOptionalParam + fmt.Sprintf(`
		if !plan.%[1]s.IsNull() && !plan.%[1]s.IsUnknown() {
			reqParams.%[1]s = plan.%[1]s.ValueString()
		},`, util.ToPascalCase(val)) + "\n"
	}

	t.funcMap = funcMap
	t.providerName = spec.Provider.Name
	t.packageName = packageName
	t.refreshObjectName = refreshObjectName
	t.importStateLogic = MakeImportStateLogic(importStateOverride)
	t.model = model
	t.refreshLogic = refreshLogic
	t.refreshWithResponse = MakeRefreshFromResponse(attributes, resourceName)
	t.endpoint = spec.Provider.Endpoint
	t.deletePathParams = extractPathParams(targetResourceRequest.Delete.Path)
	t.updatePathParams = extractPathParams(targetResourceRequest.Update[0].Path)
	t.readPathParams = extractReadPathParams(targetResourceRequest.Read.Path)
	t.createPathParams = extractPathParams(targetResourceRequest.Create.Path)
	t.deleteMethod = targetResourceRequest.Delete.Method
	t.updateMethod = targetResourceRequest.Update[0].Method
	t.readMethod = targetResourceRequest.Read.Method
	t.createMethod = targetResourceRequest.Create.Method
	t.createReqBody = createReqBody
	t.updateReqBody = updateReqBody
	t.readReqBody = readReqBody
	t.deleteReqBody = deleteReqBody
	t.createOpOptionalParam = createOpOptionalParam
	t.updateOpOptionalParam = updateOpOptionalParam
	t.createMethodName = strings.ToUpper(targetResourceRequest.Create.Method) + getMethodName(targetResourceRequest.Create.Path)
	t.readMethodName = strings.ToUpper(targetResourceRequest.Read.Method) + getMethodName(targetResourceRequest.Read.Path)
	t.updateMethodName = strings.ToUpper(targetResourceRequest.Update[0].Method) + getMethodName(targetResourceRequest.Update[0].Path)
	t.deleteMethodName = strings.ToUpper(targetResourceRequest.Delete.Method) + getMethodName(targetResourceRequest.Delete.Path)
	t.idGetter = makeIdGetter(id)

	return t
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

func MakeImportStateLogic(target string) string {
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
					diagnostics.AddError("CONVERSION ERROR", fmt.Sprintf("Failed to convert %[1]s to bool: \%v", err))
					return
				}
				postPlan.%[3]s = types.BoolValue(tempVal)
			}`, val.Name, util.ToPascalCase(resourceName), util.ToPascalCase(val.Name)) + "\n")
		} else if val.Int64 != nil {
			s.WriteString(fmt.Sprintf(`
			if !response.%[2]s.Attributes()["%[1]s"].IsNull() || !response.%[2]s.Attributes()["%[1]s"].IsUnknown() {
				tempVal, err := strconv.Atoi(response.%[2]s.Attributes()["%[1]s"].String())
				if err != nil {
					diagnostics.AddError("CONVERSION ERROR", fmt.Sprintf("Failed to convert %[1]s to int: \%v", err))
					return
				}
				postPlan.%[3]s = types.Int64Value(int64(tempVal))
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
