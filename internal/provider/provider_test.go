package provider

import (
	"os"
	"testing"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"readarr": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	if v := os.Getenv("READARR_URL"); v == "" {
		t.Skip("READARR_URL must be set for acceptance tests")
	}

	if v := os.Getenv("READARR_API_KEY"); v == "" {
		t.Skip("READARR_API_KEY must be set for acceptance tests")
	}
}

func testAccAPIClient() *readarr.APIClient {
	config := readarr.NewConfiguration()
	config.AddDefaultHeader("X-API-Key", os.Getenv("READARR_API_KEY"))
	config.Servers[0].URL = os.Getenv("READARR_URL")

	return readarr.NewAPIClient(config)
}

const testUnauthorizedProvider = `
provider "readarr" {
	url = "http://localhost:8787"
	api_key = "ErrorAPIKey"
  }
`
