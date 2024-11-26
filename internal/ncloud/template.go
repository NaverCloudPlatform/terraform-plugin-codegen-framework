package ncloud

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/util"
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-spec/resource"
)

// 실제 데이터를 생성하기 위하여 config.yml과 code-spec.json에서 데이터를 추출하고, 해당 데이터를 기반으로 리시버 별로 코드를 렌더링.
// New(): 코드 생성에 필요한 데이터들을 추출합니다. 현재는 config.yml과 code-spec.json에서 데이터를 추출하지만 추후 code-spec.json으로 전부 통일 예정.
// RenderInitial(): 초기에 필요한 작은 코드 블록들을 생성한다.
// RenderCreate(): Create 함수를 생성한다.
// RenderRead()): Read 함수를 생성한다.
// RenderUpdate(): Update 함수를 생성한다.
// RenderDelete(): Delete 함수를 생성한다.
// 필요 데이터들을 초기화 시 계산하고, 각 메서드 별 렌더링을 수행한다..
type Template struct {
	configPath       string
	codeSpecPath     string
	providerName     string
	resourceName     string
	dtoName          string
	model            string
	refreshLogic     string
	endpoint         string
	deletePathParams string
	updatePathParams string
	readPathParams   string
	createPathParams string
	deleteMethod     string
	updateMethod     string
	readMethod       string
	createMethod     string
	createReqBody    string
	updateReqBody    string
	idGetter         string
	funcMap          template.FuncMap
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

func (t *Template) RenderCreate() []byte {
	var b bytes.Buffer

	createTemplate, err := template.New("").Funcs(t.funcMap).Parse(CreateTemplate)
	if err != nil {
		log.Fatalf("error occurred with baseTemplate at rendering create: %v", err)
	}

	data := struct {
		ResourceName     string
		DtoName          string
		CreateReqBody    string
		CreateMethod     string
		Endpoint         string
		CreatePathParams string
		IdGetter         string
	}{
		ResourceName:     t.resourceName,
		DtoName:          t.dtoName,
		CreateReqBody:    t.createReqBody,
		CreateMethod:     t.createMethod,
		Endpoint:         t.endpoint,
		CreatePathParams: t.createPathParams,
		IdGetter:         t.idGetter,
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
		ResourceName   string
		DtoName        string
		ReadPathParams string
	}{
		ResourceName:   t.resourceName,
		DtoName:        t.dtoName,
		ReadPathParams: t.readPathParams,
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
		ResourceName     string
		DtoName          string
		UpdateReqBody    string
		UpdateMethod     string
		Endpoint         string
		UpdatePathParams string
		ReadPathParams   string
	}{
		ResourceName:     t.resourceName,
		DtoName:          t.dtoName,
		UpdateReqBody:    t.updateReqBody,
		UpdateMethod:     t.updateMethod,
		Endpoint:         t.endpoint,
		UpdatePathParams: t.updatePathParams,
		ReadPathParams:   t.readPathParams,
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
		ResourceName     string
		DtoName          string
		DeleteMethod     string
		Endpoint         string
		DeletePathParams string
		IdGetter         string
	}{
		ResourceName:     t.resourceName,
		DtoName:          t.dtoName,
		DeleteMethod:     t.deleteMethod,
		Endpoint:         t.endpoint,
		DeletePathParams: t.deletePathParams,
		IdGetter:         t.idGetter,
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
		DtoName string
		Model   string
	}{
		DtoName: t.dtoName,
		Model:   t.model,
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
		DtoName        string
		RefreshLogic   string
		ReadMethod     string
		Endpoint       string
		ReadPathParams string
	}{
		DtoName:        t.dtoName,
		RefreshLogic:   t.refreshLogic,
		ReadMethod:     t.readMethod,
		Endpoint:       t.endpoint,
		ReadPathParams: t.readPathParams,
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
		ReadMethod     string
		Endpoint       string
		ReadPathParams string
		DtoName        string
	}{
		ReadMethod:     t.readMethod,
		Endpoint:       t.endpoint,
		ReadPathParams: t.readPathParams,
		DtoName:        t.dtoName,
	}

	err = waitTemplate.ExecuteTemplate(&b, "Wait", data)
	if err != nil {
		log.Fatalf("error occurred with Generating wait: %v", err)
	}

	return b.Bytes()
}

// 초기화를 통해 필요한 데이터들을 미리 계산한다.
func New(configPath, codeSpecPath, resourceName string) *Template {
	var dtoName string
	var id string
	var attributes resource.Attributes
	var createReqBody string
	var updateReqBody string

	t := &Template{
		configPath:   configPath,
		codeSpecPath: codeSpecPath,
		resourceName: resourceName,
	}

	funcMap := util.CreateFuncMap()

	t.funcMap = funcMap

	codeSpec := util.ExtractAttribute(codeSpecPath)

	for _, resource := range codeSpec.Resources {
		if resource.Name == resourceName {
			dtoName = resource.DtoName
			id = resource.Id
			attributes = resource.Schema.Attributes
		}
	}

	refreshLogic, model, err := Gen_ConvertOAStoTFTypes(attributes)
	if err != nil {
		log.Fatalf("error occurred with Gen_ConvertOAStoTFTypes: %v", err)
	}

	targetResourceRequest := util.ExtractRequest(codeSpecPath, resourceName)

	for _, val := range targetResourceRequest.Create.RequestBody.Required {
		createReqBody = createReqBody + fmt.Sprintf(`"%[1]s": clearDoubleQuote(plan.%[2]s.String()),`, val, util.FirstAlphabetToUpperCase(val)) + "\n"
	}

	for _, val := range targetResourceRequest.Update[0].RequestBody.Required {
		updateReqBody = updateReqBody + fmt.Sprintf(`"%[1]s": clearDoubleQuote(plan.%[2]s.String()),`, val, util.FirstAlphabetToUpperCase(val)) + "\n"
	}

	t.providerName = codeSpec.Provider["name"].(string)
	t.dtoName = dtoName
	t.model = model
	t.refreshLogic = refreshLogic
	t.endpoint = codeSpec.Provider["endpoint"].(string)
	t.deletePathParams = extractPathParams(targetResourceRequest.Delete.Path)
	t.updatePathParams = extractPathParams(targetResourceRequest.Update[0].Path)
	t.readPathParams = extractReadPathParams(targetResourceRequest.Read.Path)
	t.createPathParams = extractCreatePathParams(targetResourceRequest.Create.Path)
	t.deleteMethod = targetResourceRequest.Delete.Method
	t.updateMethod = targetResourceRequest.Update[0].Method
	t.readMethod = targetResourceRequest.Read.Method
	t.createMethod = targetResourceRequest.Create.Method
	t.createReqBody = createReqBody
	t.updateReqBody = updateReqBody
	t.idGetter = makeIdGetter(id)

	return t
}

func extractPathParams(path string) string {
	parts := strings.Split(path, "/")
	s := ``

	for _, val := range parts {

		if len(val) < 1 {
			continue
		}

		s = s + `+"/"+`

		start := strings.Index(val, "{")

		// if val doesn't wrapped with curly brace
		if start == -1 {
			s = s + fmt.Sprintf(`"%s"`, val)
		} else {
			s = s + fmt.Sprintf(`clearDoubleQuote(plan.%s.String())`, util.PathToPascal(val))
		}
	}

	return s
}

func extractCreatePathParams(path string) string {
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
				s = s + `clearDoubleQuote(plan.ID.String())`
			} else {
				s = s + fmt.Sprintf(`clearDoubleQuote(plan.%s.String())`, util.PathToPascal(val))
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
			s = s + fmt.Sprintf(`clearDoubleQuote(plan.%s.String())`, util.PathToPascal(val))
		}
	}

	return s
}

func makeIdGetter(target string) string {
	s := "response"
	parts := strings.Split(target, ".")

	for idx, val := range parts {
		if idx == len(parts)-1 {
			s = s + fmt.Sprintf(`["%s"].(string)`, util.ToCamelCase(val))
			continue
		}

		s = s + fmt.Sprintf(`["%s"].(map[string]interface{})`, util.ToCamelCase(val))
	}

	return s
}
