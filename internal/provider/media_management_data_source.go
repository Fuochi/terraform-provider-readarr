package provider

import (
	"context"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const mediaManagementDataSourceName = "media_management"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MediaManagementDataSource{}

func NewMediaManagementDataSource() datasource.DataSource {
	return &MediaManagementDataSource{}
}

// MediaManagementDataSource defines the media management implementation.
type MediaManagementDataSource struct {
	client *readarr.APIClient
}

func (d *MediaManagementDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + mediaManagementDataSourceName
}

func (d *MediaManagementDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Media Management -->[Media Management](../resources/media_management).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Media Management ID.",
				Computed:            true,
			},
			"unmonitor_previous_books": schema.BoolAttribute{
				MarkdownDescription: "Unmonitor deleted files.",
				Computed:            true,
			},
			"hardlinks_copy": schema.BoolAttribute{
				MarkdownDescription: "Use hardlinks instead of copy.",
				Computed:            true,
			},
			"create_empty_author_folders": schema.BoolAttribute{
				MarkdownDescription: "Create empty author directories.",
				Computed:            true,
			},
			"delete_empty_folders": schema.BoolAttribute{
				MarkdownDescription: "Delete empty directories.",
				Computed:            true,
			},
			"watch_ibrary_for_changes": schema.BoolAttribute{
				MarkdownDescription: "Watch library for changes.",
				Computed:            true,
			},
			"import_extra_files": schema.BoolAttribute{
				MarkdownDescription: "Import extra files. If enabled it will leverage 'extra_file_extensions'.",
				Computed:            true,
			},
			"set_permissions": schema.BoolAttribute{
				MarkdownDescription: "Set permission for imported files.",
				Computed:            true,
			},
			"skip_free_space_check": schema.BoolAttribute{
				MarkdownDescription: "Skip free space check before importing.",
				Computed:            true,
			},
			"minimum_free_space": schema.Int64Attribute{
				MarkdownDescription: "Minimum free space in MB to allow import.",
				Computed:            true,
			},
			"recycle_bin_days": schema.Int64Attribute{
				MarkdownDescription: "Recyle bin days of retention.",
				Computed:            true,
			},
			"chmod_folder": schema.StringAttribute{
				MarkdownDescription: "Permission in linux format.",
				Computed:            true,
			},
			"chown_group": schema.StringAttribute{
				MarkdownDescription: "Group used for permission.",
				Computed:            true,
			},
			"download_propers_repacks": schema.StringAttribute{
				MarkdownDescription: "Download proper and repack policy. valid inputs are: 'preferAndUpgrade', 'doNotUpgrade', and 'doNotPrefer'.",
				Computed:            true,
			},
			"allow_fingerprinting": schema.StringAttribute{
				MarkdownDescription: "Allow fingerprinting. valid inputs are: 'newFiles', 'allFiles' and 'never'.",
				Computed:            true,
			},
			"extra_file_extensions": schema.StringAttribute{
				MarkdownDescription: "Comma separated list of extra files to import (.nfo will be imported as .nfo-orig).",
				Computed:            true,
			},
			"file_date": schema.StringAttribute{
				MarkdownDescription: "Define the file date modification. valid inputs are: 'none', and 'bookReleaseDate'.",
				Computed:            true,
			},
			"recycle_bin_path": schema.StringAttribute{
				MarkdownDescription: "Recycle bin absolute path.",
				Computed:            true,
			},
			"rescan_after_refresh": schema.StringAttribute{
				MarkdownDescription: "Rescan after refresh policy. valid inputs are: 'always', 'afterManual' and 'never'.",
				Computed:            true,
			},
		},
	}
}

func (d *MediaManagementDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *MediaManagementDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get indexer config current value
	response, _, err := d.client.MediaManagementConfigApi.GetMediaManagementConfig(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, mediaManagementDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+mediaManagementDataSourceName)

	state := MediaManagement{}
	state.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
