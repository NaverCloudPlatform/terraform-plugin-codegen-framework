{{ define "Test_DataSource" }}
/* =================================================================================
 * Test Template
 * Required data are as follows
 *
		ProviderName   string
		DataSourceName string
		PackageName    string
		ConfigParams   string
 * ================================================================================= */

package {{.PackageName}}_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	. "github.com/NaverCloudPlatform/terraform-codegen-poc/internal/test"
)

func TestAccDataSourceNcloud{{.ProviderName | ToPascalCase}}_{{.DataSourceName | ToLowerCase}}_basic(t *testing.T) {
	{{.DataSourceName | ToCamelCase}}Name := fmt.Sprintf("tf-{{.DataSourceName | ToCamelCase}}-%s", acctest.RandString(5))

	datasourceName := "ncloud_{{.ProviderName | ToLowerCase}}_{{.DataSourceName | ToLowerCase}}.testing_{{.DataSourceName | ToLowerCase}}"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: test.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAcc{{.DataSourceName | ToLowerCase}}Config({{.DataSourceName | ToCamelCase}}Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "{{.DataSourceName | ToCamelCase}}_name", {{.DataSourceName | ToCamelCase}}Name),
                    // check all the other attributes
				),
			},
		},
	})
}

func testAcc{{.DataSourceName | ToLowerCase}}Config({{.DataSourceName | ToCamelCase}}Name string) string {
	return fmt.Sprintf(`
	resource "ncloud_{{.ProviderName | ToLowerCase}}_{{.DataSourceName | ToLowerCase}}" "testing_{{.DataSourceName | ToLowerCase}}" {
		{{.ConfigParams}}
	}`, {{.DataSourceName | ToCamelCase}}Name)
}

{{ end }}