package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationWebhookResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationWebhookResourceConfig("resourceWebhookTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("readarr_notification_webhook.test", "on_upgrade", "false"),
					resource.TestCheckResourceAttrSet("readarr_notification_webhook.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccNotificationWebhookResourceConfig("resourceWebhookTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("readarr_notification_webhook.test", "on_upgrade", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "readarr_notification_webhook.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationWebhookResourceConfig(name, upgrade string) string {
	return fmt.Sprintf(`
	resource "readarr_notification_webhook" "test" {
		on_grab                            = false
		on_download_failure                = true
		on_upgrade                         = %s
		on_rename                          = false
		on_import_failure                  = false
		on_book_delete                    = false
		on_book_file_delete               = false
		on_book_file_delete_for_upgrade   = true
		on_health_issue                   = false
		on_book_retag 					  = false
		on_author_delete                  = false
		on_release_import                 = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		url = "http://transmission:9091"
		method = 1
	}`, upgrade, name)
}