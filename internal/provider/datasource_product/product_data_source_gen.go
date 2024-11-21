package datasource_product

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"strings"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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

func ProductDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
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
			"productid": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "product-id",
				MarkdownDescription: "product-id",
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

func (a *productResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan PostproductresponseModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan = *getAndRefresh(resp.Diagnostics, plan.ID.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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
			convertedTempProduct, err := util.ConvertMapToObject(context.TODO(), tempProduct)
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

func getAndRefresh(diagnostics diag.Diagnostics, id string, rest ...interface{}) *PostproductresponseModel {
	getExecFunc := func(timestamp, accessKey, signature string) *exec.Cmd {
		return exec.Command("curl", "-s", "-X", "GET", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"products"+"/"+util.ClearDoubleQuote(id),
			"-H", "Content-Type: application/json",
			"-H", "x-ncp-apigw-timestamp: "+timestamp,
			"-H", "x-ncp-iam-access-key: "+accessKey,
			"-H", "x-ncp-apigw-signature-v2: "+signature,
			"-H", "cache-control: no-cache",
			"-H", "pragma: no-cache",
		)
	}

	response, _ := util.Request(getExecFunc, "GET", "/api/v1"+"/"+"products"+"/"+util.ClearDoubleQuote(id), os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"), "")
	if response == nil {
		diagnostics.AddError("UPDATING ERROR", "response invalid")
		return nil
	}

	newPlan, err := ConvertToFrameworkTypes(util.ConvertKeys(response).(map[string]interface{}), id, rest)
	if err != nil {
		diagnostics.AddError("CREATING ERROR", err.Error())
		return nil
	}

	return newPlan
}

func waitResourceCreated(ctx context.Context, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING"},
		Target:  []string{"CREATED"},
		Refresh: func() (interface{}, string, error) {
			getExecFunc := func(timestamp, accessKey, signature string) *exec.Cmd {
			return exec.Command("curl", "-s", "-X", "GET",  "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"products"+"/"+util.ClearDoubleQuote(id),
					"-H", "accept: application/json;charset=UTF-8",
					"-H", "Content-Type: application/json",
					"-H", "x-ncp-apigw-timestamp: "+timestamp,
					"-H", "x-ncp-iam-access-key: "+accessKey,
					"-H", "x-ncp-apigw-signature-v2: "+signature,
					"-H", "cache-control: no-cache",
					"-H", "pragma: no-cache",
				)
			}

			response, err := util.Request(getExecFunc, "GET","/api/v1"+"/"+"products"+"/"+util.ClearDoubleQuote(id), os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"), "")
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

func waitResourceDeleted(ctx context.Context, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			getExecFunc := func(timestamp, accessKey, signature string) *exec.Cmd {
			return exec.Command("curl", "-s", "-X", "GET", "https://apigateway.apigw.ntruss.com/api/v1"+"/"+"products"+"/"+util.ClearDoubleQuote(id),
					"-H", "accept: application/json;charset=UTF-8",
					"-H", "Content-Type: application/json",
					"-H", "x-ncp-apigw-timestamp: "+timestamp,
					"-H", "x-ncp-iam-access-key: "+accessKey,
					"-H", "x-ncp-apigw-signature-v2: "+signature,
					"-H", "cache-control: no-cache",
					"-H", "pragma: no-cache",
				)
			}

			response, _ := util.Request(getExecFunc, "GET", "/api/v1"+"/"+"products"+"/"+util.ClearDoubleQuote(id), os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"), "")
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
