package resource_product

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"strings"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"os/exec"
	"time"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/common"
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/conn"
	"github.com/NaverCloudPlatform/terraform-plugin-codegen-framework/internal/util"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func ProductResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Description<br>Length(Min/Max): 0/300",
				MarkdownDescription: "Description<br>Length(Min/Max): 0/300",
			},
			"product": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"action_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Action Name",
						MarkdownDescription: "Action Name",
					},
					"disabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Disabled",
						MarkdownDescription: "Disabled",
					},
					"domain_code": schema.StringAttribute{
						Computed:            true,
						Description:         "Domain Code",
						MarkdownDescription: "Domain Code",
					},
					"invoke_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Invoke Id",
						MarkdownDescription: "Invoke Id",
					},
					"is_deleted": schema.BoolAttribute{
						Computed:            true,
						Description:         "Is Deleted",
						MarkdownDescription: "Is Deleted",
					},
					"is_published": schema.BoolAttribute{
						Computed:            true,
						Description:         "Is Published",
						MarkdownDescription: "Is Published",
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
					"permission": schema.StringAttribute{
						Computed:            true,
						Description:         "Permission",
						MarkdownDescription: "Permission",
					},
					"product_description": schema.StringAttribute{
						Computed:            true,
						Description:         "Product Description",
						MarkdownDescription: "Product Description",
					},
					"product_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Product Id",
						MarkdownDescription: "Product Id",
					},
					"product_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Product Name",
						MarkdownDescription: "Product Name",
					},
					"subscription_code": schema.StringAttribute{
						Computed:            true,
						Description:         "Subscription Code<br>Allowable values: PROTECTED, PUBLIC",
						MarkdownDescription: "Subscription Code<br>Allowable values: PROTECTED, PUBLIC",
					},
					"tenant_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Tenant Id",
						MarkdownDescription: "Tenant Id",
					},
				},
				CustomType: ProductType{
					ObjectType: types.ObjectType{
						AttrTypes: ProductValue{}.AttributeTypes(ctx),
					},
				},
				Computed: true,
			},
			"product_name": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Product Name<br>Length(Min/Max): 0/100",
				MarkdownDescription: "Product Name<br>Length(Min/Max): 0/100",
			},
			"productid": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "product-id",
				MarkdownDescription: "product-id",
			},
			"subscription_code": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Subscription Code<br>Allowable values: PROTECTED, PUBLIC",
				MarkdownDescription: "Subscription Code<br>Allowable values: PROTECTED, PUBLIC",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"PROTECTED",
						"PUBLIC",
					),
				},
			},
		},
	}
}

func NewProductResource() resource.Resource {
	return &productResource{}
}

type productResource struct {
	config *conn.ProviderConfig
}

func (a *productResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (a *productResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (a *productResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apigw_product"
}

func (a *productResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ProductResourceSchema(ctx)
}

func (a *productResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PostproductresponseModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := json.Marshal(map[string]string{
		"productName": clearDoubleQuote(plan.ProductName.String()),
"subscriptionCode": clearDoubleQuote(plan.SubscriptionCode.String()),

	})
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	tflog.Info(ctx, "CreateProduct reqParams="+strings.Replace(string(reqBody), `\"`, "", -1))

	response, err := util.MakeReqeust("POST", "/api/v1", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"products", strings.Replace(string(reqBody), `\"`, "", -1))
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}
	if response == nil {
		resp.Diagnostics.AddError("CREATING ERROR", "response invalid")
		return
	}

	err = waitResourceCreated(ctx, response["product"].(map[string]interface{})["productId"].(string), plan)
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	tflog.Info(ctx, "CreateProduct response="+common.MarshalUncheckedString(response))

	plan = *getAndRefresh(resp.Diagnostics, plan, response["product"].(map[string]interface{})["productId"].(string))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (a *productResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan PostproductresponseModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan = *getAndRefresh(resp.Diagnostics, plan, plan.ID.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (a *productResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PostproductresponseModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := json.Marshal(map[string]string{
		"productName": clearDoubleQuote(plan.ProductName.String()),
"subscriptionCode": clearDoubleQuote(plan.SubscriptionCode.String()),

	})
	if err != nil {
		resp.Diagnostics.AddError("CREATING ERROR", err.Error())
		return
	}

	tflog.Info(ctx, "UpdateProduct reqParams="+strings.Replace(string(reqBody), `\"`, "", -1))

	response, err := util.MakeReqeust("PATCH",  "/api/v1", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"products"+"/"+clearDoubleQuote(plan.Productid.String()), strings.Replace(string(reqBody), `\"`, "", -1))
	if err != nil {
		resp.Diagnostics.AddError("UPDATING ERROR", err.Error())
		return
	}
	if response == nil {
		resp.Diagnostics.AddError("UPDATING ERROR", "response invalid")
		return
	}

	tflog.Info(ctx, "UpdateProduct response="+common.MarshalUncheckedString(response))

	plan = *getAndRefresh(resp.Diagnostics, plan, plan.ID.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (a *productResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan PostproductresponseModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := util.MakeReqeust("DELETE", "/api/v1", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"products"+"/"+clearDoubleQuote(plan.Productid.String()), "")
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

type PostproductresponseModel struct {
    ID types.String `tfsdk:"id"`
    Description         types.String `tfsdk:"description"`
ProductName         types.String `tfsdk:"product_name"`
SubscriptionCode         types.String `tfsdk:"subscription_code"`
Product         types.Object `tfsdk:"product"`
Productid         types.String `tfsdk:"productid"`

}

func ConvertToFrameworkTypes(data map[string]interface{}, id string, rest []interface{}) (*PostproductresponseModel, error) {
	var dto PostproductresponseModel

	dto.ID = types.StringValue(id)

    dto.Description = types.StringValue(data["description"].(string))
dto.ProductName = types.StringValue(data["product_name"].(string))
dto.SubscriptionCode = types.StringValue(data["subscription_code"].(string))

			tempProduct := data["product"].(map[string]interface{})
			convertedTempProduct, err := convertMapToObject(context.TODO(), tempProduct)
			if err != nil {
				fmt.Println("ConvertMapToObject Error")
			}

			dto.Product = diagOff(types.ObjectValueFrom, context.TODO(), types.ObjectType{AttrTypes: map[string]attr.Type{
				"action_name": types.StringType,
"disabled": types.BoolType,
"domain_code": types.StringType,
"invoke_id": types.StringType,
"is_deleted": types.BoolType,
"is_published": types.BoolType,
"mod_time": types.StringType,
"modifier": types.StringType,
"permission": types.StringType,
"product_description": types.StringType,
"product_id": types.StringType,
"product_name": types.StringType,
"subscription_code": types.StringType,
"tenant_id": types.StringType,

			}}.AttributeTypes(), convertedTempProduct)
dto.Productid = types.StringValue(data["productid"].(string))


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

func getAndRefresh(diagnostics diag.Diagnostics, plan PostproductresponseModel, id string, rest ...interface{}) *PostproductresponseModel {
	response, err := util.MakeReqeust("GET", "/api/v1", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"products"+"/"+clearDoubleQuote(id), "")
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

func waitResourceCreated(ctx context.Context, id string, plan PostproductresponseModel) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING"},
		Target:  []string{"CREATED"},
		Refresh: func() (interface{}, string, error) {
			response, err := util.MakeReqeust("GET", "/api/v1", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"products"+"/"+clearDoubleQuote(id), "")
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

func waitResourceDeleted(ctx context.Context, id string, plan PostproductresponseModel) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			response, _ := util.MakeReqeust("GET", "/api/v1", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"products"+"/"+clearDoubleQuote(id), "")
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

