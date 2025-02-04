package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const authorsDataSourceName = "authors"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AuthorsDataSource{}

func NewAuthorsDataSource() datasource.DataSource {
	return &AuthorsDataSource{}
}

// AuthorsDataSource defines the authors implementation.
type AuthorsDataSource struct {
	client *readarr.APIClient
}

// Authors describes the authors data model.
type Authors struct {
	Authors types.Set    `tfsdk:"authors"`
	ID      types.String `tfsdk:"id"`
}

func (d *AuthorsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + authorsDataSourceName
}

func (d *AuthorsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Authors -->List all available [Authors](../resources/author).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"authors": schema.SetNestedAttribute{
				MarkdownDescription: "Author list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
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
							Computed:            true,
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
				},
			},
		},
	}
}

func (d *AuthorsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *AuthorsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get authors current value
	response, _, err := d.client.AuthorAPI.ListAuthor(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.List, authorsDataSourceName, err))
		return
	}

	tflog.Trace(ctx, "read "+authorsDataSourceName)

	pointerResponse := make([]*readarr.AuthorResource, len(response))
	for i := range response {
		pointerResponse[i] = &response[i]
	}

	authors := make([]Author, len(pointerResponse))
	for i, m := range pointerResponse {
		authors[i].write(ctx, m, &resp.Diagnostics)
	}

	authorList, diags := types.SetValueFrom(ctx, Author{}.getType(), authors)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, Authors{Authors: authorList, ID: types.StringValue(strconv.Itoa(len(response)))})...)
}
