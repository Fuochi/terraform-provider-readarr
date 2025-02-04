package provider

import (
	"context"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const namingDataSourceName = "naming"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &NamingDataSource{}

func NewNamingDataSource() datasource.DataSource {
	return &NamingDataSource{}
}

// NamingDataSource defines the naming implementation.
type NamingDataSource struct {
	client *readarr.APIClient
}

func (d *NamingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + namingDataSourceName
}

func (d *NamingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Media Management -->[Naming](../resources/naming).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Naming ID.",
				Computed:            true,
			},
			"rename_books": schema.BoolAttribute{
				MarkdownDescription: "Readarr will use the existing file name if false.",
				Computed:            true,
			},
			"replace_illegal_characters": schema.BoolAttribute{
				MarkdownDescription: "Replace illegal characters. They will be removed if false.",
				Computed:            true,
			},
			"colon_replacement_format": schema.Int64Attribute{
				MarkdownDescription: "Change how Readarr handles colon replacement. '0' Delete, '1' Dash, '2' Space Dash, '3' Space Dash Space, '4' Smart.",
				Computed:            true,
			},
			"author_folder_format": schema.StringAttribute{
				MarkdownDescription: "Author folder format.",
				Computed:            true,
			},
			"standard_book_format": schema.StringAttribute{
				MarkdownDescription: "Standard book formatss.",
				Computed:            true,
			},
		},
	}
}

func (d *NamingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *NamingDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get naming current value
	response, _, err := d.client.NamingConfigAPI.GetNamingConfig(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, namingDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+namingDataSourceName)

	state := Naming{}
	state.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
