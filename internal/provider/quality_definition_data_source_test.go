package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccQualityDefinitionDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccQualityDefinitionDataSourceConfig(999) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccQualityDefinitionDataSourceConfig(999),
				ExpectError: regexp.MustCompile("Unable to find quality_definition"),
			},
			// Read testing
			{
				Config: testAccQualityDefinitionDataSourceConfig(2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.readarr_quality_definition.test", "title"),
					resource.TestCheckResourceAttr("data.readarr_quality_definition.test", "id", "2")),
			},
		},
	})
}

func testAccQualityDefinitionDataSourceConfig(id int) string {
	return fmt.Sprintf(`
	data "readarr_quality_definition" "test" {
		id = %d
	}
	`, id)
}
