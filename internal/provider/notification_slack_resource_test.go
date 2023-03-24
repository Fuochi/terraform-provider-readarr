package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationSlackResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationSlackResourceConfig("resourceSlackTest", "test") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationSlackResourceConfig("resourceSlackTest", "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("readarr_notification_slack.test", "channel", "test"),
					resource.TestCheckResourceAttrSet("readarr_notification_slack.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationSlackResourceConfig("resourceSlackTest", "test") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationSlackResourceConfig("resourceSlackTest", "test1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("readarr_notification_slack.test", "channel", "test1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "readarr_notification_slack.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationSlackResourceConfig(name, channel string) string {
	return fmt.Sprintf(`
	resource "readarr_notification_slack" "test" {
		on_grab                           = false
		on_download_failure               = false
		on_upgrade                        = false
		on_rename                         = false
		on_import_failure                 = false
		on_book_delete                    = false
		on_book_file_delete               = false
		on_book_file_delete_for_upgrade   = false
		on_health_issue                   = false
		on_book_retag 					  = false
		on_author_delete                  = false
		on_release_import                 = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		web_hook_url = "http://my.slack.com/test"
		username = "user"
		channel = "%s"
	}`, name, channel)
}