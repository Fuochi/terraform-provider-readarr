package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccHostResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccHostResourceConfig("readarr", "test") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccHostResourceConfig("readarr", "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("readarr_host.test", "port", "8989"),
					resource.TestCheckResourceAttrSet("readarr_host.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccHostResourceConfig("readarr", "test") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccHostResourceConfig("readarrTest", "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("readarr_host.test", "port", "8989"),
				),
			},
			// Update and Read testing
			{
				Config: testAccHostResourceConfig("readarrTest", "test1234"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("readarr_host.test", "port", "8989"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "readarr_host.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "test1234",
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccHostResourceConfig(name, pass string) string {
	return fmt.Sprintf(`
	resource "readarr_host" "test" {
		launch_browser = true
		port = 8989
		url_base = ""
		bind_address = "*"
		application_url =  ""
		instance_name = "%s"
		proxy = {
			enabled = false
		}
		ssl = {
			enabled = false
			certificate_validation = "enabled"
		}
		logging = {
			log_level = "info"
		}
		backup = {
			folder = "/backup"
			interval = 5
			retention = 10
		}
		authentication = {
			method = "basic"
			username = "test"
			password = "%s"
		}
		update = {
			mechanism = "docker"
			branch = "develop"
		}
	}`, name, pass)
}
