package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRootFolderDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccRootFolderDataSourceConfig("/error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccRootFolderDataSourceConfig("/error"),
				ExpectError: regexp.MustCompile("Unable to find root_folder"),
			},
			// Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccRootFolderDataSourceConfig("/config"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.readarr_root_folder.test", "id"),
					resource.TestCheckResourceAttr("data.readarr_root_folder.test", "path", "/config")),
			},
		},
	})
}

func testAccRootFolderDataSourceConfig(path string) string {
	return fmt.Sprintf(`
	data "readarr_root_folder" "test" {
  			path = "%s"
		}
	`, path)
}

func rootFolderDSInit() {
	// ensure a /config root path is configured
	client := testAccAPIClient()
	folder := readarr.NewRootFolderResource()
	folder.SetPath("/config")
	folder.SetDefaultMetadataProfileId(1)
	folder.SetDefaultQualityProfileId(1)
	folder.SetDefaultMonitorOption("all")
	folder.SetOutputProfile("default")
	folder.SetHost("")
	folder.SetIsCalibreLibrary(false)
	folder.SetName("Config")
	tags := make([]int32, 0)
	folder.SetDefaultTags(tags)
	_, _, _ = client.RootFolderAPI.CreateRootFolder(context.TODO()).RootFolderResource(*folder).Execute()
}
