package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientHadoukenResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDownloadClientHadoukenResourceConfig("resourceHadoukenTest", "hadouken") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccDownloadClientHadoukenResourceConfig("resourceHadoukenTest", "hadouken"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("readarr_download_client_hadouken.test", "host", "hadouken"),
					resource.TestCheckResourceAttr("readarr_download_client_hadouken.test", "url_base", "/hadouken/"),
					resource.TestCheckResourceAttrSet("readarr_download_client_hadouken.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDownloadClientHadoukenResourceConfig("resourceHadoukenTest", "hadouken") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientHadoukenResourceConfig("resourceHadoukenTest", "hadouken-host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("readarr_download_client_hadouken.test", "host", "hadouken-host"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "readarr_download_client_hadouken.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientHadoukenResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "readarr_download_client_hadouken" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "%s"
		url_base = "/hadouken/"
		port = 9091
		category = "readarr-tv"
		username = "username"
		password = "password"
	}`, name, host)
}
