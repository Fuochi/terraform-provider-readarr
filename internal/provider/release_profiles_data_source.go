package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const releaseProfilesDataSourceName = "release_profiles"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ReleaseProfilesDataSource{}

func NewReleaseProfilesDataSource() datasource.DataSource {
	return &ReleaseProfilesDataSource{}
}

// ReleaseProfilesDataSource defines the release profiles implementation.
type ReleaseProfilesDataSource struct {
	client *readarr.APIClient
}

// ReleaseProfiles describes the release profiles data model.
type ReleaseProfiles struct {
	ReleaseProfiles types.Set    `tfsdk:"release_profiles"`
	ID              types.String `tfsdk:"id"`
}

func (d *ReleaseProfilesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + releaseProfilesDataSourceName
}

func (d *ReleaseProfilesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the release server.
		MarkdownDescription: "<!-- subcategory:Profiles -->List all available [Release Profiles](../resources/release_profile).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"release_profiles": schema.SetNestedAttribute{
				MarkdownDescription: "Release Profile list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Release Profile ID.",
							Computed:            true,
						},
						"enabled": schema.BoolAttribute{
							MarkdownDescription: "Enabled.",
							Computed:            true,
						},
						"include_preferred_when_renaming": schema.BoolAttribute{
							MarkdownDescription: "Include preferred when renaming flag.",
							Computed:            true,
						},
						"indexer_id": schema.Int64Attribute{
							MarkdownDescription: "Indexer ID. Set `0` for all.",
							Computed:            true,
						},
						"required": schema.StringAttribute{
							MarkdownDescription: "Required terms. Comma separated list. At least one of `required` and `ignored` must be set.",
							Computed:            true,
						},
						"ignored": schema.StringAttribute{
							MarkdownDescription: "Ignored terms. Comma separated list. At least one of `required` and `ignored` must be set.",
							Computed:            true,
						},
						"preferred": schema.SetNestedAttribute{
							MarkdownDescription: "Preferred terms.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"score": schema.Int64Attribute{
										MarkdownDescription: "Score.",
										Computed:            true,
									},
									"term": schema.StringAttribute{
										MarkdownDescription: "Term.",
										Computed:            true,
									},
								},
							},
						},
						"tags": schema.SetAttribute{
							MarkdownDescription: "List of associated tags.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
					},
				},
			},
		},
	}
}

func (d *ReleaseProfilesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *ReleaseProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ReleaseProfiles

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get releaseprofiles current value
	response, _, err := d.client.ReleaseProfileApi.ListReleaseProfile(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, releaseProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+releaseProfileResourceName)
	// Map response body to resource schema attribute
	profiles := make([]ReleaseProfile, len(response))
	for i, p := range response {
		profiles[i].write(ctx, p)
	}

	tfsdk.ValueFrom(ctx, profiles, data.ReleaseProfiles.Type(ctx), &data.ReleaseProfiles)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}