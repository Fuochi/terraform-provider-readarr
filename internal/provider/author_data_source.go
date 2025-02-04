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

const authorDataSourceName = "author"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AuthorDataSource{}

func NewAuthorDataSource() datasource.DataSource {
	return &AuthorDataSource{}
}

// AuthorDataSource defines the author implementation.
type AuthorDataSource struct {
	client *readarr.APIClient
}

func (d *AuthorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + authorDataSourceName
}

func (d *AuthorDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Authors -->Single [Author](../resources/author).",
		Attributes: map[string]schema.Attribute{
			"monitored": schema.BoolAttribute{
				MarkdownDescription: "Monitored flag.",
				Computed:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality profile ID.",
				Computed:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Author ID.",
				Computed:            true,
			},
			"author_name": schema.StringAttribute{
				MarkdownDescription: "Author name.",
				Computed:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Full author path.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Author status.",
				Computed:            true,
			},
			"overview": schema.StringAttribute{
				MarkdownDescription: "Overview.",
				Computed:            true,
			},
			"foreign_author_id": schema.StringAttribute{
				MarkdownDescription: "Foreign author ID.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"genres": schema.SetAttribute{
				MarkdownDescription: "List genres.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *AuthorDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *AuthorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Author

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get authors current value
	response, _, err := d.client.AuthorAPI.ListAuthor(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, authorDataSourceName, err))

		return
	}

	pointerResponse := make([]*readarr.AuthorResource, len(response))
	for i := range response {
		pointerResponse[i] = &response[i]
	}
	data.find(ctx, data.ForeignAuthorID.ValueString(), pointerResponse, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+authorDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (a *Author) find(ctx context.Context, ID string, authors []*readarr.AuthorResource, diags *diag.Diagnostics) {
	for _, author := range authors {
		if author.GetForeignAuthorId() == ID {
			a.write(ctx, author, diags)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(authorDataSourceName, "Foreign author ID", ID))
}
