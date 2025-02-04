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

const importListDataSourceName = "import_list"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ImportListDataSource{}

func NewImportListDataSource() datasource.DataSource {
	return &ImportListDataSource{}
}

// ImportListDataSource defines the import_list implementation.
type ImportListDataSource struct {
	client *readarr.APIClient
}

func (d *ImportListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListDataSourceName
}

func (d *ImportListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Import Lists -->Single [Import List](../resources/import_list).",
		Attributes: map[string]schema.Attribute{
			"enable_automatic_add": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic add flag.",
				Computed:            true,
			},
			"should_monitor_existing": schema.BoolAttribute{
				MarkdownDescription: "Should monitor existing flag.",
				Computed:            true,
			},
			"should_search": schema.BoolAttribute{
				MarkdownDescription: "Should search flag.",
				Computed:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality profile ID.",
				Computed:            true,
			},
			"metadata_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Metadata profile ID.",
				Computed:            true,
			},
			"list_order": schema.Int64Attribute{
				MarkdownDescription: "List order.",
				Computed:            true,
			},
			"root_folder_path": schema.StringAttribute{
				MarkdownDescription: "Root folder path.",
				Computed:            true,
			},
			"should_monitor": schema.StringAttribute{
				MarkdownDescription: "Should monitor.",
				Computed:            true,
			},
			"monitor_new_items": schema.StringAttribute{
				MarkdownDescription: "Monitor new items.",
				Computed:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "ImportList implementation name.",
				Computed:            true,
			},
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "ImportList configuration template.",
				Computed:            true,
			},
			"list_type": schema.StringAttribute{
				MarkdownDescription: "List type.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Import List name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Import List ID.",
				Computed:            true,
			},
			// Field values
			"series_id": schema.Int64Attribute{
				MarkdownDescription: "Series ID.",
				Computed:            true,
			},
			"list_id": schema.Int64Attribute{
				MarkdownDescription: "List ID.",
				Computed:            true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token.",
				Computed:            true,
				Sensitive:           true,
			},
			"access_token_secret": schema.StringAttribute{
				MarkdownDescription: "Access token secret.",
				Computed:            true,
				Sensitive:           true,
			},
			"request_token_secret": schema.StringAttribute{
				MarkdownDescription: "Request token secret.",
				Computed:            true,
				Sensitive:           true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Computed:            true,
				Sensitive:           true,
			},
			"user_id": schema.StringAttribute{
				MarkdownDescription: "User ID.",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Computed:            true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Computed:            true,
			},
			"profile_ids": schema.SetAttribute{
				MarkdownDescription: "Profile IDs.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"tag_ids": schema.SetAttribute{
				MarkdownDescription: "Tag IDs.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"bookshelf_ids": schema.SetAttribute{
				MarkdownDescription: "Bookshelf IDs.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *ImportListDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *ImportListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ImportList

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get importList current value
	response, _, err := d.client.ImportListAPI.ListImportList(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListDataSourceName, err))

		return
	}

	lists := make([]*readarr.ImportListResource, len(response))
	for i := range response {
		lists[i] = &response[i]
	}
	data.find(ctx, data.Name.ValueString(), lists, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+importListDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (i *ImportList) find(ctx context.Context, name string, importLists []*readarr.ImportListResource, diags *diag.Diagnostics) {
	for _, list := range importLists {
		if list.GetName() == name {
			i.write(ctx, list, diags)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(importListDataSourceName, "name", name))
}
