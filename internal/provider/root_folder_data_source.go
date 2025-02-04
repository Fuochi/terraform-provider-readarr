package provider

import (
	"context"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const rootFolderDataSourceName = "root_folder"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RootFolderDataSource{}

func NewRootFolderDataSource() datasource.DataSource {
	return &RootFolderDataSource{}
}

// RootFolderDataSource defines the root folders implementation.
type RootFolderDataSource struct {
	client *readarr.APIClient
}

func (d *RootFolderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + rootFolderDataSourceName
}

func (d *RootFolderDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Media Management -->Single [Root Folder](../resources/root_folder).",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				MarkdownDescription: "Root Folder absolute path.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Root Folder name.",
				Computed:            true,
			},
			"default_monitor_option": schema.StringAttribute{
				MarkdownDescription: "Default monitor option.",
				Computed:            true,
			},
			"default_monitor_new_item_option": schema.StringAttribute{
				MarkdownDescription: "Default monitor new item option.",
				Computed:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Calibre host.",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Calibre username.",
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Calibre password.",
				Computed:            true,
				Sensitive:           true,
			},
			"library": schema.StringAttribute{
				MarkdownDescription: "Calibre library.",
				Computed:            true,
			},
			"output_profile": schema.StringAttribute{
				MarkdownDescription: "Calibre output profile.",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Calibre Port.",
				Computed:            true,
			},
			"default_metadata_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Default metadata profile ID.",
				Computed:            true,
			},
			"default_quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Default metadata profile ID.",
				Computed:            true,
			},
			"accessible": schema.BoolAttribute{
				MarkdownDescription: "Access flag.",
				Computed:            true,
			},
			"is_calibre_library": schema.BoolAttribute{
				MarkdownDescription: "Is calibre library flag.",
				Computed:            true,
			},
			"use_ssl": schema.BoolAttribute{
				MarkdownDescription: "Use SSL for calibre connection.",
				Computed:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Root Folder ID.",
				Computed:            true,
			},
			"default_tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (d *RootFolderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *RootFolderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var folder *RootFolder

	resp.Diagnostics.Append(req.Config.Get(ctx, &folder)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get rootfolders current value
	response, _, err := d.client.RootFolderAPI.ListRootFolder(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, rootFolderDataSourceName, err))

		return
	}

	folders := make([]*readarr.RootFolderResource, len(response))
	for i := range response {
		folders[i] = &response[i]
	}
	folder.find(ctx, folder.Path.ValueString(), folders, &resp.Diagnostics)

	tflog.Trace(ctx, "read "+rootFolderDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &folder)...)
}

func (r *RootFolder) find(ctx context.Context, path string, folders []*readarr.RootFolderResource, diags *diag.Diagnostics) {
	for _, folder := range folders {
		if folder.GetPath() == path {
			r.write(ctx, folder, diags)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(rootFolderDataSourceName, "path", path))
}
