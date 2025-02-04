package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImportListLazyLibrarianResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListLazyLibrarianResourceConfig("resourceLazyLibrarianTest", "entireAuthor") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListLazyLibrarianResourceConfig("resourceLazyLibrarianTest", "entireAuthor"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("readarr_import_list_lazy_librarian.test", "should_monitor", "entireAuthor"),
					resource.TestCheckResourceAttrSet("readarr_import_list_lazy_librarian.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListLazyLibrarianResourceConfig("resourceLazyLibrarianTest", "entireAuthor") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListLazyLibrarianResourceConfig("resourceLazyLibrarianTest", "specificBook"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("readarr_import_list_lazy_librarian.test", "should_monitor", "specificBook"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "readarr_import_list_lazy_librarian.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListLazyLibrarianResourceConfig(name, folder string) string {
	return fmt.Sprintf(`
	resource "readarr_import_list_lazy_librarian" "test" {
		enable_automatic_add = false
		should_monitor = "%s"
		monitor_new_items = "none"
		should_search = false
		root_folder_path = "/config"
		quality_profile_id = 1
		metadata_profile_id = 1
		name = "%s"
		base_url = "http://localhost:5299"
		api_key = "APIKey"
	}`, folder, name)
}
