package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const metadataProfileResourceName = "metadata_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MetadataProfileResource{}
	_ resource.ResourceWithImportState = &MetadataProfileResource{}
)

func NewMetadataProfileResource() resource.Resource {
	return &MetadataProfileResource{}
}

// MetadataProfileResource defines the metadata profile implementation.
type MetadataProfileResource struct {
	client *readarr.APIClient
}

// MetadataProfile describes the metadata profile data model.
type MetadataProfile struct {
	Name                types.String  `tfsdk:"name"`
	AllowedLanguages    types.String  `tfsdk:"allowed_languages"`
	Ignored             types.String  `tfsdk:"ignored"`
	ID                  types.Int64   `tfsdk:"id"`
	MinPages            types.Int64   `tfsdk:"min_pages"`
	MinPopularity       types.Float64 `tfsdk:"min_popularity"`
	SkipMissingDate     types.Bool    `tfsdk:"skip_missing_date"`
	SkipMissingIsbn     types.Bool    `tfsdk:"skip_missing_isbn"`
	SkipPartsAndSets    types.Bool    `tfsdk:"skip_parts_and_sets"`
	SkipSeriesSecondary types.Bool    `tfsdk:"skip_series_secondary"`
}

func (r *MetadataProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->Metadata Profile resource.\nFor more information refer to [Metadata Profile](https://wiki.servarr.com/readarr/settings#metadata-profiles) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Metadata Profile ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Metadata Profile name.",
				Required:            true,
			},
			"allowed_languages": schema.StringAttribute{
				MarkdownDescription: "Allowed languages. Comma separated list of ISO 639-3 language codes.",
				Optional:            true,
				Computed:            true,
			},
			"ignored": schema.StringAttribute{
				MarkdownDescription: "Terms to ignore. Comma separated list.",
				Optional:            true,
				Computed:            true,
			},
			"min_popularity": schema.Float64Attribute{
				MarkdownDescription: "Minimum popularity.",
				Optional:            true,
				Computed:            true,
			},
			"min_pages": schema.Int64Attribute{
				MarkdownDescription: "Minimum pages.",
				Optional:            true,
				Computed:            true,
			},
			"skip_missing_date": schema.BoolAttribute{
				MarkdownDescription: "Skip missing date.",
				Optional:            true,
				Computed:            true,
			},
			"skip_missing_isbn": schema.BoolAttribute{
				MarkdownDescription: "Skip missing ISBN.",
				Optional:            true,
				Computed:            true,
			},
			"skip_parts_and_sets": schema.BoolAttribute{
				MarkdownDescription: "Skip parts and sets.",
				Optional:            true,
				Computed:            true,
			},
			"skip_series_secondary": schema.BoolAttribute{
				MarkdownDescription: "Skip secondary series books.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *MetadataProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataProfileResourceName
}

func (r *MetadataProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *MetadataProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var profile *MetadataProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new MetadataProfile
	request := profile.read()

	response, _, err := r.client.MetadataProfileApi.CreateMetadataProfile(ctx).MetadataProfileResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, metadataProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+metadataProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	profile.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *MetadataProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var profile *MetadataProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get metadataProfile current value
	response, _, err := r.client.MetadataProfileApi.GetMetadataProfileById(ctx, int32(profile.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	profile.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *MetadataProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var profile *MetadataProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update MetadataProfile
	request := profile.read()

	response, _, err := r.client.MetadataProfileApi.UpdateMetadataProfile(ctx, strconv.Itoa(int(request.GetId()))).MetadataProfileResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, metadataProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+metadataProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	profile.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *MetadataProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var profile *MetadataProfile

	diags := req.State.Get(ctx, &profile)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete metadataProfile current value
	_, err := r.client.MetadataProfileApi.DeleteMetadataProfile(ctx, int32(profile.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+metadataProfileResourceName+": "+strconv.Itoa(int(profile.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *MetadataProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+metadataProfileResourceName+": "+req.ID)
}

func (p *MetadataProfile) write(profile *readarr.MetadataProfileResource) {
	p.ID = types.Int64Value(int64(profile.GetId()))
	p.Name = types.StringValue(profile.GetName())
	p.MinPopularity = types.Float64Value(profile.GetMinPopularity())
	p.MinPages = types.Int64Value(int64(profile.GetMinPages()))
	p.AllowedLanguages = types.StringValue(profile.GetAllowedLanguages())
	p.Ignored = types.StringValue(profile.GetIgnored())
	p.SkipMissingDate = types.BoolValue(profile.GetSkipMissingDate())
	p.SkipMissingIsbn = types.BoolValue(profile.GetSkipMissingIsbn())
	p.SkipPartsAndSets = types.BoolValue(profile.GetSkipPartsAndSets())
	p.SkipSeriesSecondary = types.BoolValue(profile.GetSkipSeriesSecondary())
}

func (p *MetadataProfile) read() *readarr.MetadataProfileResource {
	profile := readarr.NewMetadataProfileResource()
	profile.SetName(p.Name.ValueString())
	profile.SetId(int32(p.ID.ValueInt64()))
	profile.SetMinPopularity(p.MinPopularity.ValueFloat64())
	profile.SetMinPages(int32(p.MinPages.ValueInt64()))
	profile.SetAllowedLanguages(p.AllowedLanguages.ValueString())
	profile.SetIgnored(p.Ignored.ValueString())
	profile.SetSkipMissingDate(p.SkipMissingDate.ValueBool())
	profile.SetSkipMissingIsbn(p.SkipMissingIsbn.ValueBool())
	profile.SetSkipPartsAndSets(p.SkipPartsAndSets.ValueBool())
	profile.SetSkipSeriesSecondary(p.SkipSeriesSecondary.ValueBool())

	return profile
}