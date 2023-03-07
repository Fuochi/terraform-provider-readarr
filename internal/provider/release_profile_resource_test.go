package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccReleaseProfileResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccReleaseProfileResourceConfig("test1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccReleaseProfileResourceConfig("test1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("readarr_release_profile.test", "required", "test1"),
					resource.TestCheckResourceAttrSet("readarr_release_profile.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccReleaseProfileResourceConfig("test1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccReleaseProfileResourceConfig("test2,test3"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("readarr_release_profile.test", "required", "test2,test3"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "readarr_release_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccReleaseProfileResourceConfig(required string) string {
	return fmt.Sprintf(`
	resource "readarr_release_profile" "test" {
		enabled = true
		indexer_id = 0
		required = "%s"

		preferred = [
			{
				term = "test"
				score = 100
			}
		] 
	}`, required)
}
