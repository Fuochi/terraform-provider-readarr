package provider

import (
	"context"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const metadataConfigDataSourceName = "metadata_config"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MetadataConfigDataSource{}

func NewMetadataConfigDataSource() datasource.DataSource {
	return &MetadataConfigDataSource{}
}

// MetadataConfigDataSource defines the metadata config implementation.
type MetadataConfigDataSource struct {
	client *readarr.APIClient
}

func (d *MetadataConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataConfigDataSourceName
}

func (d *MetadataConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Metadata -->[Metadata Config](../resources/metadata_config).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Metadata Config ID.",
				Computed:            true,
			},
			"write_audio_tags": schema.StringAttribute{
				MarkdownDescription: "Write audio tags.",
				Computed:            true,
			},
			"write_book_tags": schema.StringAttribute{
				MarkdownDescription: "Write book tags.",
				Computed:            true,
			},
			"scrub_audio_tags": schema.BoolAttribute{
				MarkdownDescription: "Scrub audio tags.",
				Computed:            true,
			},
			"update_covers": schema.BoolAttribute{
				MarkdownDescription: "Update covers.",
				Computed:            true,
			},
			"embed_metadata": schema.BoolAttribute{
				MarkdownDescription: "Embed metadata in book files.",
				Computed:            true,
			},
		},
	}
}

func (d *MetadataConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *MetadataConfigDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get metadata config current value
	response, _, err := d.client.MetadataProviderConfigAPI.GetMetadataProviderConfig(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataConfigDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataConfigDataSourceName)

	status := MetadataConfig{}
	status.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, status)...)
}
