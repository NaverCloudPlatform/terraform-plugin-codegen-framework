package resource_resource

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"

	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/common"
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/conn"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func ResourceResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"apiid": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "api-id",
				MarkdownDescription: "api-id",
			},
			"cors_allow_credentials": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Cors Allow Credentials<br>Length(Min/Max): 0/45<br>Pattern: ^$|\\S|^\\S|\\S$|^\\S.*\\S$",
				MarkdownDescription: "Cors Allow Credentials<br>Length(Min/Max): 0/45<br>Pattern: ^$|\\S|^\\S|\\S$|^\\S.*\\S$",
			},
			"cors_allow_headers": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Cors Allow Headers<br>Length(Min/Max): 0/200<br>Pattern: ^$|\\S|^\\S|\\S$|^\\S.*\\S$",
				MarkdownDescription: "Cors Allow Headers<br>Length(Min/Max): 0/200<br>Pattern: ^$|\\S|^\\S|\\S$|^\\S.*\\S$",
			},
			"cors_allow_methods": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Cors Allow Methods<br>Length(Min/Max): 0/45<br>Pattern: ^$|\\S|^\\S|\\S$|^\\S.*\\S$",
				MarkdownDescription: "Cors Allow Methods<br>Length(Min/Max): 0/45<br>Pattern: ^$|\\S|^\\S|\\S$|^\\S.*\\S$",
			},
			"cors_allow_origin": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Cors Allow Origin<br>Length(Min/Max): 0/45<br>Pattern: ^$|\\S|^\\S|\\S$|^\\S.*\\S$",
				MarkdownDescription: "Cors Allow Origin<br>Length(Min/Max): 0/45<br>Pattern: ^$|\\S|^\\S|\\S$|^\\S.*\\S$",
			},
			"cors_expose_headers": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Cors Expose Headers<br>Length(Min/Max): 0/200<br>Pattern: ^$|\\S|^\\S|\\S$|^\\S.*\\S$",
				MarkdownDescription: "Cors Expose Headers<br>Length(Min/Max): 0/200<br>Pattern: ^$|\\S|^\\S|\\S$|^\\S.*\\S$",
			},
			"cors_max_age": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Cors Max Age<br>Length(Min/Max): 0/45<br>Pattern: ^$|\\S|^\\S|\\S$|^\\S.*\\S$",
				MarkdownDescription: "Cors Max Age<br>Length(Min/Max): 0/45<br>Pattern: ^$|\\S|^\\S|\\S$|^\\S.*\\S$",
			},
			"productid": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "product-id",
				MarkdownDescription: "product-id",
			},
			"resource": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"api_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Api Id",
						MarkdownDescription: "Api Id",
					},
					"cors_allow_credentials": schema.StringAttribute{
						Computed:            true,
						Description:         "Cors Allow Credentials",
						MarkdownDescription: "Cors Allow Credentials",
					},
					"cors_allow_headers": schema.StringAttribute{
						Computed:            true,
						Description:         "Cors Allow Headers",
						MarkdownDescription: "Cors Allow Headers",
					},
					"cors_allow_methods": schema.StringAttribute{
						Computed:            true,
						Description:         "Cors Allow Methods",
						MarkdownDescription: "Cors Allow Methods",
					},
					"cors_allow_origin": schema.StringAttribute{
						Computed:            true,
						Description:         "Cors Allow Origin",
						MarkdownDescription: "Cors Allow Origin",
					},
					"cors_expose_headers": schema.StringAttribute{
						Computed:            true,
						Description:         "Cors Expose Headers",
						MarkdownDescription: "Cors Expose Headers",
					},
					"cors_max_age": schema.StringAttribute{
						Computed:            true,
						Description:         "Cors Max Age",
						MarkdownDescription: "Cors Max Age",
					},
					"methods": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"method_code": schema.StringAttribute{
									Computed:            true,
									Description:         "Method Code",
									MarkdownDescription: "Method Code",
								},
								"method_name": schema.StringAttribute{
									Computed:            true,
									Description:         "Method Name<br>Allowable values: ANY, GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD",
									MarkdownDescription: "Method Name<br>Allowable values: ANY, GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD",
								},
							},
						},
						Computed:            true,
						Description:         "Methods",
						MarkdownDescription: "Methods",
					},
					"resource_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Resource Id",
						MarkdownDescription: "Resource Id",
					},
					"resource_path": schema.StringAttribute{
						Computed:            true,
						Description:         "Resource Path",
						MarkdownDescription: "Resource Path",
					},
				},
				Computed: true,
			},
			"resource_list": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"api_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Api Id",
							MarkdownDescription: "Api Id",
						},
						"cors_allow_credentials": schema.StringAttribute{
							Computed:            true,
							Description:         "Cors Allow Credentials",
							MarkdownDescription: "Cors Allow Credentials",
						},
						"cors_allow_headers": schema.StringAttribute{
							Computed:            true,
							Description:         "Cors Allow Headers",
							MarkdownDescription: "Cors Allow Headers",
						},
						"cors_allow_methods": schema.StringAttribute{
							Computed:            true,
							Description:         "Cors Allow Methods",
							MarkdownDescription: "Cors Allow Methods",
						},
						"cors_allow_origin": schema.StringAttribute{
							Computed:            true,
							Description:         "Cors Allow Origin",
							MarkdownDescription: "Cors Allow Origin",
						},
						"cors_expose_headers": schema.StringAttribute{
							Computed:            true,
							Description:         "Cors Expose Headers",
							MarkdownDescription: "Cors Expose Headers",
						},
						"cors_max_age": schema.StringAttribute{
							Computed:            true,
							Description:         "Cors Max Age",
							MarkdownDescription: "Cors Max Age",
						},
						"methods": schema.ListNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"method_code": schema.StringAttribute{
										Computed:            true,
										Description:         "Method Code",
										MarkdownDescription: "Method Code",
									},
									"method_name": schema.StringAttribute{
										Computed:            true,
										Description:         "Method Name<br>Allowable values: ANY, GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD",
										MarkdownDescription: "Method Name<br>Allowable values: ANY, GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD",
									},
								},
							},
							Computed:            true,
							Description:         "Methods",
							MarkdownDescription: "Methods",
						},
						"resource_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Resource Id",
							MarkdownDescription: "Resource Id",
						},
						"resource_path": schema.StringAttribute{
							Computed:            true,
							Description:         "Resource Path",
							MarkdownDescription: "Resource Path",
						},
					},
				},
				Computed:            true,
				Description:         "Resource List",
				MarkdownDescription: "Resource List",
			},
			"resource_path": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Resource Path<br>Length(Min/Max): 0/2,000<br>Pattern: /|(/([\\w+\\-]*\\w|\\{[\\w+\\-]*\\w+((})|(\\+}))))+",
				MarkdownDescription: "Resource Path<br>Length(Min/Max): 0/2,000<br>Pattern: /|(/([\\w+\\-]*\\w|\\{[\\w+\\-]*\\w+((})|(\\+}))))+",
			},
		},
	}
}

func NewResourceResource() resource.Resource {
	return &resourceResource{}
}

type resourceResource struct {
	config *conn.ProviderConfig
}

func (a *resourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*conn.ProviderConfig)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	a.config = config
}

func (a *resourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (a *resourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apigw_resource"
}

func (a *resourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceResourceSchema(ctx)
}

func (a *resourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ResourcedtoModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := json.Marshal(map[string]string{
		"resourcePath": clearDoubleQuote(plan.ResourcePath.String()),
	})
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	tflog.Info(ctx, "CreateResource reqParams="+strings.Replace(string(reqBody), `\"`, "", -1))

	execFunc := func(timestamp, accessKey, signature string) *exec.Cmd {
		return exec.Command("curl", "-s", "-X", "POST", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"products"+"/"+clearDoubleQuote(plan.Productid.String())+"/"+"apis"+"/"+clearDoubleQuote(plan.Apiid.String())+"/"+"resources",
			"-H", "Content-Type: application/json",
			"-H", "x-ncp-apigw-timestamp: "+timestamp,
			"-H", "x-ncp-iam-access-key: "+accessKey,
			"-H", "x-ncp-apigw-signature-v2: "+signature,
			"-H", "cache-control: no-cache",
			"-H", "pragma: no-cache",
			"-d", strings.Replace(string(reqBody), `\"`, "", -1),
		)
	}

	response, err := request(execFunc, "POST", "/api/v1"+"/"+"products"+"/"+clearDoubleQuote(plan.Productid.String())+"/"+"apis"+"/"+clearDoubleQuote(plan.Apiid.String())+"/"+"resources", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"), strings.Replace(string(reqBody), `\"`, "", -1))
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}
	if response == nil {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	err = waitResourceCreated(ctx, response["resource"].(map[string]interface{})["resourceid"].(string), plan)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	tflog.Info(ctx, "CreateResource response="+common.MarshalUncheckedString(response))

	plan = *getAndRefresh(resp.Diagnostics, plan, response["resource"].(map[string]interface{})["resourceid"].(string))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (a *resourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan ResourcedtoModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan = *getAndRefresh(resp.Diagnostics, plan, plan.ID.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (a *resourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ResourcedtoModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := json.Marshal(map[string]string{})
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	tflog.Info(ctx, "UpdateResource reqParams="+strings.Replace(string(reqBody), `\"`, "", -1))

	execFunc := func(timestamp, accessKey, signature string) *exec.Cmd {
		return exec.Command("curl", "-s", "-X", "PATCH", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"products"+"/"+clearDoubleQuote(plan.Productid.String())+"/"+"apis"+"/"+clearDoubleQuote(plan.Apiid.String())+"/"+"resources"+"/"+clearDoubleQuote(plan.Resourceid.String()),
			"-H", "Content-Type: application/json",
			"-H", "x-ncp-apigw-timestamp: "+timestamp,
			"-H", "x-ncp-iam-access-key: "+accessKey,
			"-H", "x-ncp-apigw-signature-v2: "+signature,
			"-H", "cache-control: no-cache",
			"-H", "pragma: no-cache",
			"-d", strings.Replace(string(reqBody), `\"`, "", -1),
		)
	}

	response, err := request(execFunc, "PATCH", "/api/v1"+"/"+"products"+"/"+clearDoubleQuote(plan.Productid.String())+"/"+"apis"+"/"+clearDoubleQuote(plan.Apiid.String())+"/"+"resources"+"/"+clearDoubleQuote(plan.Resourceid.String()), os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"), strings.Replace(string(reqBody), `\"`, "", -1))
	if err != nil {
		resp.Diagnostics.AddError("UPDATING ERROR", err.Error())
		return
	}
	if response == nil {
		resp.Diagnostics.AddError("UPDATING ERROR", "response invalid")
		return
	}

	tflog.Info(ctx, "UpdateResource response="+common.MarshalUncheckedString(response))

	plan = *getAndRefresh(resp.Diagnostics, plan, plan.ID.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (a *resourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan ResourcedtoModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	execFunc := func(timestamp, accessKey, signature string) *exec.Cmd {
		return exec.Command("curl", "-s", "-X", "DELETE", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"products"+"/"+clearDoubleQuote(plan.Productid.String())+"/"+"apis"+"/"+clearDoubleQuote(plan.Apiid.String())+"/"+"resources"+"/"+clearDoubleQuote(plan.Resourceid.String()),
			"-H", "Content-Type: application/json",
			"-H", "x-ncp-apigw-timestamp: "+timestamp,
			"-H", "x-ncp-iam-access-key: "+accessKey,
			"-H", "x-ncp-apigw-signature-v2: "+signature,
			"-H", "cache-control: no-cache",
			"-H", "pragma: no-cache",
		)
	}

	_, err := request(execFunc, "DELETE", "/api/v1"+"/"+"products"+"/"+clearDoubleQuote(plan.Productid.String())+"/"+"apis"+"/"+clearDoubleQuote(plan.Apiid.String())+"/"+"resources"+"/"+clearDoubleQuote(plan.Resourceid.String()), os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"), "")
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}

	err = waitResourceDeleted(ctx, clearDoubleQuote(plan.ID.String()), plan)
	if err != nil {
		resp.Diagnostics.AddError("DELETING ERROR", err.Error())
		return
	}
}

type ResourcedtoModel struct {
	ID                   types.String `tfsdk:"id"`
	CorsAllowCredentials types.String `tfsdk:"cors_allow_credentials"`
	CorsAllowHeaders     types.String `tfsdk:"cors_allow_headers"`
	CorsAllowMethods     types.String `tfsdk:"cors_allow_methods"`
	CorsAllowOrigin      types.String `tfsdk:"cors_allow_origin"`
	CorsExposeHeaders    types.String `tfsdk:"cors_expose_headers"`
	CorsMaxAge           types.String `tfsdk:"cors_max_age"`
	ResourcePath         types.String `tfsdk:"resource_path"`
	Resource             types.Object `tfsdk:"resource"`
	Resource_list        types.List   `tfsdk:"resource_list"`
	Productid            types.String `tfsdk:"productid"`
	Apiid                types.String `tfsdk:"apiid"`
}

func ConvertToFrameworkTypes(data map[string]interface{}, id string, rest []interface{}) (*ResourcedtoModel, error) {
	var dto ResourcedtoModel

	dto.ID = types.StringValue(id)

	dto.CorsAllowCredentials = types.StringValue(data["cors_allow_credentials"].(string))
	dto.CorsAllowHeaders = types.StringValue(data["cors_allow_headers"].(string))
	dto.CorsAllowMethods = types.StringValue(data["cors_allow_methods"].(string))
	dto.CorsAllowOrigin = types.StringValue(data["cors_allow_origin"].(string))
	dto.CorsExposeHeaders = types.StringValue(data["cors_expose_headers"].(string))
	dto.CorsMaxAge = types.StringValue(data["cors_max_age"].(string))
	dto.ResourcePath = types.StringValue(data["resource_path"].(string))

	tempResource := data["resource"].(map[string]interface{})
	convertedTempResource, err := convertMapToObject(context.TODO(), tempResource)
	if err != nil {
		fmt.Println("ConvertMapToObject Error")
	}

	dto.Resource = diagOff(types.ObjectValueFrom, context.TODO(), types.ObjectType{AttrTypes: map[string]attr.Type{
		"api_id":                 types.StringType,
		"cors_allow_credentials": types.StringType,
		"cors_allow_headers":     types.StringType,
		"cors_allow_methods":     types.StringType,
		"cors_allow_origin":      types.StringType,
		"cors_expose_headers":    types.StringType,
		"cors_max_age":           types.StringType,

		"methods": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

			"method_code": types.StringType,
			"method_name": types.StringType,
		},
		}},
		"resource_id":   types.StringType,
		"resource_path": types.StringType,
	}}.AttributeTypes(), convertedTempResource)

	tempResource_list := data["resource_list"].([]interface{})
	dto.Resource_list = diagOff(types.ListValueFrom, context.TODO(), types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{

		"api_id":                 types.StringType,
		"cors_allow_credentials": types.StringType,
		"cors_allow_headers":     types.StringType,
		"cors_allow_methods":     types.StringType,
		"cors_allow_origin":      types.StringType,
		"cors_expose_headers":    types.StringType,
		"cors_max_age":           types.StringType,
		"resource_id":            types.StringType,
		"resource_path":          types.StringType,
	},
	}}.ElementType(), tempResource_list)
	dto.Productid = types.StringValue(data["productid"].(string))
	dto.Apiid = types.StringValue(data["apiid"].(string))

	return &dto, nil
}

func diagOff[V, T interface{}](input func(ctx context.Context, elementType T, elements any) (V, diag.Diagnostics), ctx context.Context, elementType T, elements any) V {
	var emptyReturn V

	v, diags := input(ctx, elementType, elements)

	if diags.HasError() {
		diags.AddError("REFRESHING ERROR", "invalid diagOff operation")
		return emptyReturn
	}

	return v
}

func getAndRefresh(diagnostics diag.Diagnostics, plan ResourcedtoModel, id string, rest ...interface{}) *ResourcedtoModel {
	getExecFunc := func(timestamp, accessKey, signature string) *exec.Cmd {
		return exec.Command("curl", "-s", "-X", "GET", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"products"+"/"+clearDoubleQuote(plan.Productid.String())+"/"+"apis"+"/"+clearDoubleQuote(plan.Apiid.String())+"/"+clearDoubleQuote(id),
			"-H", "Content-Type: application/json",
			"-H", "x-ncp-apigw-timestamp: "+timestamp,
			"-H", "x-ncp-iam-access-key: "+accessKey,
			"-H", "x-ncp-apigw-signature-v2: "+signature,
			"-H", "cache-control: no-cache",
			"-H", "pragma: no-cache",
		)
	}

	response, _ := request(getExecFunc, "GET", "/api/v1"+"/"+"products"+"/"+clearDoubleQuote(plan.Productid.String())+"/"+"apis"+"/"+clearDoubleQuote(plan.Apiid.String())+"/"+clearDoubleQuote(id), os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"), "")
	if response == nil {
		diagnostics.AddError("UPDATING ERROR", "response invalid")
		return nil
	}

	newPlan, err := ConvertToFrameworkTypes(convertKeys(response).(map[string]interface{}), id, rest)
	if err != nil {
		diagnostics.AddError("CREATING ERROR", err.Error())
		return nil
	}

	return newPlan
}

func convertKeys(input interface{}) interface{} {
	switch v := input.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for key, value := range v {
			// Convert the key to snake_case
			newKey := camelToSnake(key)
			// Recursively convert nested values
			newMap[newKey] = convertKeys(value)
		}
		return newMap
	case []interface{}:
		newSlice := make([]interface{}, len(v))
		for i, value := range v {
			newSlice[i] = convertKeys(value)
		}
		return newSlice
	default:
		return v
	}
}

func convertMapToObject(ctx context.Context, data map[string]interface{}) (types.Object, error) {
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	for key, value := range data {
		attrType, attrValue, err := convertInterfaceToAttr(ctx, value)
		if err != nil {
			return types.Object{}, fmt.Errorf("error converting field %s: %v", key, err)
		}

		attrTypes[key] = attrType
		attrValues[key] = attrValue
	}

	r, _ := types.ObjectValue(attrTypes, attrValues)

	return r, nil
}

func camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

func convertInterfaceToAttr(ctx context.Context, value interface{}) (attr.Type, attr.Value, error) {
	switch v := value.(type) {
	case string:
		return types.StringType, types.StringValue(v), nil
	case float64:
		return types.Int64Type, types.Int64Value(int64(v)), nil
	case bool:
		return types.BoolType, types.BoolValue(v), nil
	case []interface{}:
		if len(v) == 0 {
			// Treat as array list in case of empty
			return types.ListType{ElemType: types.StringType},
				types.ListValueMust(types.StringType, []attr.Value{}),
				nil
		}
		// Determine type based on first element
		elemType, _, err := convertInterfaceToAttr(ctx, v[0])
		if err != nil {
			return nil, nil, err
		}

		values := make([]attr.Value, len(v))
		for i, item := range v {
			_, value, err := convertInterfaceToAttr(ctx, item)
			if err != nil {
				return nil, nil, err
			}
			values[i] = value
		}

		listType := types.ListType{ElemType: elemType}
		listValue, diags := types.ListValue(elemType, values)
		if diags.HasError() {
			return nil, nil, err
		}

		return listType, listValue, nil

	case map[string]interface{}:
		objValue, err := convertMapToObject(ctx, v)
		if err != nil {
			return nil, nil, err
		}
		return objValue.Type(ctx), objValue, nil
	case nil:
		return types.StringType, types.StringNull(), nil
	default:
		return nil, nil, fmt.Errorf("unsupported type: %T", value)
	}
}

func makeSignature(method, url, timestamp, accessKey, secretKey string) string {
	message := fmt.Sprintf("%s %s\n%s\n%s",
		method,
		url,
		timestamp,
		accessKey,
	)

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func request(command func(timestamp, accessKey, signature string) *exec.Cmd, method, url, accessKey, secretKey, requestBody string) (map[string]interface{}, error) {
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	signature := makeSignature(method, url, timestamp, accessKey, secretKey)

	cmd := command(timestamp, accessKey, signature)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, err
	}

	// code 200 but error occurs
	if result["error"] != nil {
		return result, fmt.Errorf("error with code 200: %s", result["error"])
	}

	return result, nil
}

func clearDoubleQuote(s string) string {
	return strings.Replace(strings.Replace(strings.Replace(s, "\\", "", -1), "\"", "", -1), `"`, "", -1)
}

func waitResourceCreated(ctx context.Context, id string, plan ResourcedtoModel) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING"},
		Target:  []string{"CREATED"},
		Refresh: func() (interface{}, string, error) {
			getExecFunc := func(timestamp, accessKey, signature string) *exec.Cmd {
				return exec.Command("curl", "-s", "-X", "GET", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"products"+"/"+clearDoubleQuote(plan.Productid.String())+"/"+"apis"+"/"+clearDoubleQuote(plan.Apiid.String())+"/"+clearDoubleQuote(id),
					"-H", "accept: application/json;charset=UTF-8",
					"-H", "Content-Type: application/json",
					"-H", "x-ncp-apigw-timestamp: "+timestamp,
					"-H", "x-ncp-iam-access-key: "+accessKey,
					"-H", "x-ncp-apigw-signature-v2: "+signature,
					"-H", "cache-control: no-cache",
					"-H", "pragma: no-cache",
				)
			}

			response, err := request(getExecFunc, "GET", "/api/v1"+"/"+"products"+"/"+clearDoubleQuote(plan.Productid.String())+"/"+"apis"+"/"+clearDoubleQuote(plan.Apiid.String())+"/"+clearDoubleQuote(id), os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"), "")
			if err != nil {
				return response, "CREATING", nil
			}
			if response != nil {
				return response, "CREATED", nil
			}

			return response, "CREATING", nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error occured while waiting for resource to be created: %s", err)
	}
	return nil
}

func waitResourceDeleted(ctx context.Context, id string, plan ResourcedtoModel) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			getExecFunc := func(timestamp, accessKey, signature string) *exec.Cmd {
				return exec.Command("curl", "-s", "-X", "GET", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"products"+"/"+clearDoubleQuote(plan.Productid.String())+"/"+"apis"+"/"+clearDoubleQuote(plan.Apiid.String())+"/"+clearDoubleQuote(id),
					"-H", "accept: application/json;charset=UTF-8",
					"-H", "Content-Type: application/json",
					"-H", "x-ncp-apigw-timestamp: "+timestamp,
					"-H", "x-ncp-iam-access-key: "+accessKey,
					"-H", "x-ncp-apigw-signature-v2: "+signature,
					"-H", "cache-control: no-cache",
					"-H", "pragma: no-cache",
				)
			}

			response, _ := request(getExecFunc, "GET", "/api/v1"+"/"+"products"+"/"+clearDoubleQuote(plan.Productid.String())+"/"+"apis"+"/"+clearDoubleQuote(plan.Apiid.String())+"/"+clearDoubleQuote(id), os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"), "")
			if response["error"] != nil {
				return response, "DELETED", nil
			}

			return response, "DELETING", nil
		},
		Timeout:    conn.DefaultTimeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error occured while waiting for resource to be deleted: %s", err)
	}
	return nil
}
