package resource_api_keys

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/common"
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/conn"
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func ApiKeysResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"api_key": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"api_key_description": schema.StringAttribute{
						Computed:            true,
						Description:         "Api Key Description",
						MarkdownDescription: "Api Key Description",
					},
					"api_key_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Api Key Id",
						MarkdownDescription: "Api Key Id",
					},
					"api_key_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Api Key Name",
						MarkdownDescription: "Api Key Name",
					},
					"domain_code": schema.StringAttribute{
						Computed:            true,
						Description:         "Domain Code",
						MarkdownDescription: "Domain Code",
					},
					"is_enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Is Enabled",
						MarkdownDescription: "Is Enabled",
					},
					"mod_time": schema.StringAttribute{
						Computed:            true,
						Description:         "Mod Time",
						MarkdownDescription: "Mod Time",
					},
					"modifier": schema.StringAttribute{
						Computed:            true,
						Description:         "Modifier",
						MarkdownDescription: "Modifier",
					},
					"primary_key": schema.StringAttribute{
						Computed:            true,
						Description:         "Primary Key",
						MarkdownDescription: "Primary Key",
					},
					"secondary_key": schema.StringAttribute{
						Computed:            true,
						Description:         "Secondary Key",
						MarkdownDescription: "Secondary Key",
					},
					"tenant_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Tenant Id",
						MarkdownDescription: "Tenant Id",
					},
				},
				Computed: true,
			},
			"api_key_description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Api Key Description<br>Length(Min/Max): 0/50",
				MarkdownDescription: "Api Key Description<br>Length(Min/Max): 0/50",
			},
			"api_key_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Api Key Id",
				MarkdownDescription: "Api Key Id",
			},
			"api_key_name": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Api Key Name<br>Length(Min/Max): 0/20",
				MarkdownDescription: "Api Key Name<br>Length(Min/Max): 0/20",
			},
			"apikeyid": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "api-key-id",
				MarkdownDescription: "api-key-id",
			},
			"domain_code": schema.StringAttribute{
				Computed:            true,
				Description:         "Domain Code",
				MarkdownDescription: "Domain Code",
			},
			"is_enabled": schema.BoolAttribute{
				Computed:            true,
				Description:         "Is Enabled",
				MarkdownDescription: "Is Enabled",
			},
			"mod_time": schema.StringAttribute{
				Computed:            true,
				Description:         "Mod Time",
				MarkdownDescription: "Mod Time",
			},
			"modifier": schema.StringAttribute{
				Computed:            true,
				Description:         "Modifier",
				MarkdownDescription: "Modifier",
			},
			"primary_key": schema.StringAttribute{
				Computed:            true,
				Description:         "Primary Key",
				MarkdownDescription: "Primary Key",
			},
			"secondary_key": schema.StringAttribute{
				Computed:            true,
				Description:         "Secondary Key",
				MarkdownDescription: "Secondary Key",
			},
			"tenant_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Tenant Id",
				MarkdownDescription: "Tenant Id",
			},
		},
	}
}

func NewApiKeysResource() resource.Resource {
	return &apiKeysResource{}
}

type apiKeysResource struct {
	config *conn.ProviderConfig
}

func (a *apiKeysResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (a *apiKeysResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (a *apiKeysResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apigw_api_keys"
}

func (a *apiKeysResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ApiKeysResourceSchema(ctx)
}

func (a *apiKeysResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ApikeydtoModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := json.Marshal(map[string]string{
		"apiKeyName": clearDoubleQuote(plan.ApiKeyName.String()),
	})
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	tflog.Info(ctx, "CreateApiKeys reqParams="+strings.Replace(string(reqBody), `\"`, "", -1))

	response, err := util.MakeReqeust("POST", "/api/v1", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"api-keys", strings.Replace(string(reqBody), `\"`, "", -1))
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}
	if response == nil {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	err = waitResourceCreated(ctx, response["apiKey"].(map[string]interface{})["apiKeyId"].(string), plan)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	tflog.Info(ctx, "CreateApiKeys response="+common.MarshalUncheckedString(response))

	plan = *getAndRefresh(resp.Diagnostics, plan, response["apiKey"].(map[string]interface{})["apiKeyId"].(string))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (a *apiKeysResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan ApikeydtoModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan = *getAndRefresh(resp.Diagnostics, plan, plan.ID.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (a *apiKeysResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ApikeydtoModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := json.Marshal(map[string]string{
		"apiKeyName": clearDoubleQuote(plan.ApiKeyName.String()),
		"isEnabled":  clearDoubleQuote(plan.IsEnabled.String()),
	})
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	tflog.Info(ctx, "UpdateApiKeys reqParams="+strings.Replace(string(reqBody), `\"`, "", -1))

	response, err := util.MakeReqeust("PUT", "/api/v1", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"api-keys"+"/"+clearDoubleQuote(plan.Apikeyid.String()), strings.Replace(string(reqBody), `\"`, "", -1))
	if err != nil {
		resp.Diagnostics.AddError("UPDATING ERROR", err.Error())
		return
	}
	if response == nil {
		resp.Diagnostics.AddError("UPDATING ERROR", "response invalid")
		return
	}

	tflog.Info(ctx, "UpdateApiKeys response="+common.MarshalUncheckedString(response))

	plan = *getAndRefresh(resp.Diagnostics, plan, plan.ID.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (a *apiKeysResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan ApikeydtoModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := util.MakeReqeust("DELETE", "/api/v1", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"api-keys"+"/"+clearDoubleQuote(plan.Apikeyid.String()), "")
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

type ApikeydtoModel struct {
	ID                types.String `tfsdk:"id"`
	ApiKeyDescription types.String `tfsdk:"api_key_description"`
	ApiKeyName        types.String `tfsdk:"api_key_name"`
	Api_key           types.Object `tfsdk:"api_key"`
	ApiKeyId          types.String `tfsdk:"api_key_id"`
	DomainCode        types.String `tfsdk:"domain_code"`
	IsEnabled         types.Bool   `tfsdk:"is_enabled"`
	ModTime           types.String `tfsdk:"mod_time"`
	Modifier          types.String `tfsdk:"modifier"`
	PrimaryKey        types.String `tfsdk:"primary_key"`
	SecondaryKey      types.String `tfsdk:"secondary_key"`
	TenantId          types.String `tfsdk:"tenant_id"`
	Apikeyid          types.String `tfsdk:"apikeyid"`
}

func ConvertToFrameworkTypes(data map[string]interface{}, id string, rest []interface{}) (*ApikeydtoModel, error) {
	var dto ApikeydtoModel

	dto.ID = types.StringValue(id)

	dto.ApiKeyDescription = types.StringValue(data["api_key_description"].(string))
	dto.ApiKeyName = types.StringValue(data["api_key_name"].(string))

	tempApi_key := data["api_key"].(map[string]interface{})
	convertedTempApi_key, err := convertMapToObject(context.TODO(), tempApi_key)
	if err != nil {
		fmt.Println("ConvertMapToObject Error")
	}

	dto.Api_key = diagOff(types.ObjectValueFrom, context.TODO(), types.ObjectType{AttrTypes: map[string]attr.Type{
		"api_key_description": types.StringType,
		"api_key_id":          types.StringType,
		"api_key_name":        types.StringType,
		"domain_code":         types.StringType,
		"is_enabled":          types.BoolType,
		"mod_time":            types.StringType,
		"modifier":            types.StringType,
		"primary_key":         types.StringType,
		"secondary_key":       types.StringType,
		"tenant_id":           types.StringType,
	}}.AttributeTypes(), convertedTempApi_key)
	dto.ApiKeyId = types.StringValue(data["api_key_id"].(string))
	dto.DomainCode = types.StringValue(data["domain_code"].(string))
	dto.IsEnabled = types.BoolValue(data["is_enabled"].(bool))
	dto.ModTime = types.StringValue(data["mod_time"].(string))
	dto.Modifier = types.StringValue(data["modifier"].(string))
	dto.PrimaryKey = types.StringValue(data["primary_key"].(string))
	dto.SecondaryKey = types.StringValue(data["secondary_key"].(string))
	dto.TenantId = types.StringValue(data["tenant_id"].(string))
	dto.Apikeyid = types.StringValue(data["apikeyid"].(string))

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

func getAndRefresh(diagnostics diag.Diagnostics, plan ApikeydtoModel, id string, rest ...interface{}) *ApikeydtoModel {
	response, err := util.MakeReqeust("GET", "/api/v1", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"api-keys"+"/"+clearDoubleQuote(id), "")
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

func clearDoubleQuote(s string) string {
	return strings.Replace(strings.Replace(strings.Replace(s, "\\", "", -1), "\"", "", -1), `"`, "", -1)
}

func waitResourceCreated(ctx context.Context, id string, plan ApikeydtoModel) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING"},
		Target:  []string{"CREATED"},
		Refresh: func() (interface{}, string, error) {
			response, err := util.MakeReqeust("GET", "/api/v1", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"api-keys"+"/"+clearDoubleQuote(id), "")
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

func waitResourceDeleted(ctx context.Context, id string, plan ApikeydtoModel) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			response, _ := util.MakeReqeust("GET", "/api/v1", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"api-keys"+"/"+clearDoubleQuote(id), "")
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
