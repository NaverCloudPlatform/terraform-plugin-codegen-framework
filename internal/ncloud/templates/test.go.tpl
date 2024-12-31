{{ define "Test" }}
/* =================================================================================
 * Test Template
 * Required data are as follows
 *
		ProviderName      string
		ResourceName      string
		PackageName       string
		RefreshObjectName string
		ReadMethod        string
		ReadMethodName    string
		ReadReqBody       string
		Endpoint          string
		ReadPathParams    string
 * ================================================================================= */

package {{.PackageName}}_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/ncloudsdk"
)

func TestAccResourceNcloud{{.ProviderName | ToPascalCase}}_{{.ResourceName | ToLowerCase}}_basic(t *testing.T) {
	{{.ResourceName | ToCamelCase}}Name := fmt.Sprintf("tf-{{.ResourceName | ToCamelCase}}-%s", acctest.RandString(5))

	resourceName := "ncloud_{{.ProviderName | ToLowerCase}}_{{.ResourceName | ToLowerCase}}.testing_{{.ResourceName | ToLowerCase}}"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheck{{.ResourceName | ToPascalCase}}Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAcc{{.ResourceName | ToLowerCase}}Config({{.ResourceName | ToCamelCase}}Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheck{{.ResourceName | ToLowerCase}}Exists(resourceName, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceName, "{{.ResourceName | ToCamelCase}}_name", {{.ResourceName | ToCamelCase}}Name),

                    // check all the other attributes
				),
			},
		},
	})
}

func testAccCheck{{.ResourceName | ToLowerCase}}Exists(n string, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		c := ncloudsdk.NewClient("https://apigateway.apigw.ntruss.com/api/v1", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))

		response, err := c.{{.ReadMethodName}}_TF(&ncloudsdk.{{.ReadMethodName}}Request{
            // change value with "resource.Primary.ID"
            {{.ReadReqBody}}
		})
		if response == nil {
			return err
		}
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheck{{.ResourceName | ToPascalCase}}Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_{{.ProviderName | ToLowerCase}}_{{.ResourceName | ToLowerCase}}.testing_{{.ResourceName | ToLowerCase}}" {
			continue
		}

		c := ncloudsdk.NewClient("https://apigateway.apigw.ntruss.com/api/v1", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))
		_, err := c.{{.ReadMethodName}}_TF(&ncloudsdk.{{.ReadMethodName}}Request{
            // change value with "rs.Primary.ID"
            {{.ReadReqBody}}
		})
		if err != nil {
			return nil
		}
	}

	return nil
}

func testAcc{{.ResourceName | ToLowerCase}}Config({{.ResourceName | ToCamelCase}}Name string) string {
	return fmt.Sprintf(`
	resource "ncloud_{{.ProviderName | ToLowerCase}}_{{.ResourceName | ToLowerCase}}" "testing_{{.ResourceName | ToLowerCase}}" {
		{{.ResourceName | ToCamelCase}}_name			= "%[1]s"

        // fill the other required attributes
	}`, {{.ResourceName | ToCamelCase}}Name)
}

{{ end }}