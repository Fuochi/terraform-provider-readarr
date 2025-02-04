package provider

import (
	"context"
	"os"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// needed for tf debug mode
// var stderr = os.Stderr

// Ensure provider defined types fully satisfy framework interfaces.
var _ provider.Provider = &ReadarrProvider{}

// ScaffoldingProvider defines the provider implementation.
type ReadarrProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Readarr describes the provider data model.
type Readarr struct {
	APIKey types.String `tfsdk:"api_key"`
	URL    types.String `tfsdk:"url"`
}

func (p *ReadarrProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "readarr"
	resp.Version = p.version
}

func (p *ReadarrProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Readarr provider is used to interact with any [Readarr](https://readarr.com/) installation. You must configure the provider with the proper credentials before you can use it. Use the left navigation to read about the available resources.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key for Readarr authentication. Can be specified via the `READARR_API_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Full Readarr URL with protocol and port (e.g. `https://test.readarr.lib:8787`). You should **NOT** supply any path (`/api`), the SDK will use the appropriate paths. Can be specified via the `READARR_URL` environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *ReadarrProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data Readarr
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// User must provide URL to the provider
	if data.URL.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as url",
		)

		return
	}

	var url string
	if data.URL.IsNull() {
		url = os.Getenv("READARR_URL")
	} else {
		url = data.URL.ValueString()
	}

	if url == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find URL",
			"URL cannot be an empty string",
		)

		return
	}

	// User must provide API key to the provider
	if data.APIKey.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as api_key",
		)

		return
	}

	var key string
	if data.APIKey.IsNull() {
		key = os.Getenv("READARR_API_KEY")
	} else {
		key = data.APIKey.ValueString()
	}

	if key == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find API key",
			"API key cannot be an empty string",
		)

		return
	}

	// Configuring client. API Key management could be changed once new options avail in sdk.
	config := readarr.NewConfiguration()
	config.AddDefaultHeader("X-API-Key", key)
	config.Servers[0].URL = url
	client := readarr.NewAPIClient(config)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ReadarrProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Author
		NewAuthorResource,

		// Download Clients
		NewDownloadClientConfigResource,
		NewDownloadClientResource,
		NewDownloadClientAria2Resource,
		NewDownloadClientDelugeResource,
		NewDownloadClientFloodResource,
		NewDownloadClientHadoukenResource,
		NewDownloadClientNzbgetResource,
		NewDownloadClientNzbvortexResource,
		NewDownloadClientPneumaticResource,
		NewDownloadClientQbittorrentResource,
		NewDownloadClientRtorrentResource,
		NewDownloadClientSabnzbdResource,
		NewDownloadClientTorrentBlackholeResource,
		NewDownloadClientTorrentDownloadStationResource,
		NewDownloadClientTransmissionResource,
		NewDownloadClientUsenetBlackholeResource,
		NewDownloadClientUsenetDownloadStationResource,
		NewDownloadClientUtorrentResource,
		NewDownloadClientVuzeResource,

		// Indexers
		NewIndexerConfigResource,
		NewIndexerResource,
		NewIndexerFilelistResource,
		NewIndexerGazelleResource,
		NewIndexerIptorrentsResource,
		NewIndexerNewznabResource,
		NewIndexerNyaaResource,
		NewIndexerTorrentRssResource,
		NewIndexerTorrentleechResource,
		NewIndexerTorznabResource,

		// Import Lists
		NewImportListExclusionResource,
		NewImportListResource,
		NewImportListGoodreadsBookshelfResource,
		NewImportListGoodreadsListResource,
		NewImportListGoodreadsOwnedBooksResource,
		NewImportListGoodreadsSeriesResource,
		NewImportListLazyLibrarianResource,
		NewImportListReadarrResource,

		// Media Management
		NewNamingResource,
		NewMediaManagementResource,
		NewRemotePathMappingResource,
		NewRootFolderResource,

		// Metadata
		NewMetadataConfigResource,

		// Notifications
		NewNotificationResource,
		NewNotificationBoxcarResource,
		NewNotificationCustomScriptResource,
		NewNotificationDiscordResource,
		NewNotificationEmailResource,
		NewNotificationGoodreadsBookshelvesResource,
		NewNotificationGoodreadsOwnedBooksResource,
		NewNotificationGotifyResource,
		NewNotificationJoinResource,
		NewNotificationKavitaResource,
		NewNotificationMailgunResource,
		NewNotificationNotifiarrResource,
		NewNotificationNtfyResource,
		NewNotificationProwlResource,
		NewNotificationPushbulletResource,
		NewNotificationPushoverResource,
		NewNotificationSendgridResource,
		NewNotificationSlackResource,
		NewNotificationSubsonicResource,
		NewNotificationSynologyResource,
		NewNotificationTelegramResource,
		NewNotificationTwitterResource,
		NewNotificationWebhookResource,

		// Profiles
		NewDelayProfileResource,
		NewMetadataProfileResource,
		NewQualityProfileResource,
		NewQualityDefinitionResource,
		NewReleaseProfileResource,
		NewCustomFormatResource,

		// System
		NewHostResource,

		// Tags
		NewTagResource,
	}
}

func (p *ReadarrProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Author
		NewAuthorDataSource,
		NewAuthorsDataSource,

		// Download Clients
		NewDownloadClientConfigDataSource,
		NewDownloadClientDataSource,
		NewDownloadClientsDataSource,

		// Indexers
		NewIndexerConfigDataSource,
		NewIndexerDataSource,
		NewIndexersDataSource,

		// Import Lists
		NewImportListExclusionDataSource,
		NewImportListExclusionsDataSource,
		NewImportListDataSource,
		NewImportListsDataSource,

		// Notifications
		NewNotificationDataSource,
		NewNotificationsDataSource,

		// Media Management
		NewNamingDataSource,
		NewMediaManagementDataSource,
		NewRemotePathMappingDataSource,
		NewRemotePathMappingsDataSource,
		NewRootFolderDataSource,
		NewRootFoldersDataSource,

		// Metadata
		NewMetadataConfigDataSource,

		// Profiles
		NewCustomFormatDataSource,
		NewCustomFormatsDataSource,
		NewDelayProfileDataSource,
		NewDelayProfilesDataSource,
		NewMetadataProfileDataSource,
		NewMetadataProfilesDataSource,
		NewQualityProfileDataSource,
		NewQualityProfilesDataSource,
		NewQualityDataSource,
		NewQualityDefinitionDataSource,
		NewQualityDefinitionsDataSource,
		NewReleaseProfileDataSource,
		NewReleaseProfilesDataSource,
		NewCustomFormatConditionDataSource,
		NewCustomFormatConditionReleaseGroupDataSource,
		NewCustomFormatConditionReleaseTitleDataSource,
		NewCustomFormatConditionSizeDataSource,

		// System
		NewHostDataSource,
		NewSystemStatusDataSource,

		// Tags
		NewTagDataSource,
		NewTagsDataSource,
	}
}

// New returns the provider with a specific version.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ReadarrProvider{
			version: version,
		}
	}
}
