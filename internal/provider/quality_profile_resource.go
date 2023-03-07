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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const qualityProfileResourceName = "quality_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &QualityProfileResource{}
	_ resource.ResourceWithImportState = &QualityProfileResource{}
)

func NewQualityProfileResource() resource.Resource {
	return &QualityProfileResource{}
}

// QualityProfileResource defines the quality profile implementation.
type QualityProfileResource struct {
	client *readarr.APIClient
}

// QualityProfile describes the quality profile data model.
type QualityProfile struct {
	QualityGroups  types.Set    `tfsdk:"quality_groups"`
	Name           types.String `tfsdk:"name"`
	ID             types.Int64  `tfsdk:"id"`
	Cutoff         types.Int64  `tfsdk:"cutoff"`
	UpgradeAllowed types.Bool   `tfsdk:"upgrade_allowed"`
}

// QualityGroup is part of QualityProfile.
type QualityGroup struct {
	Qualities types.Set    `tfsdk:"qualities"`
	Name      types.String `tfsdk:"name"`
	ID        types.Int64  `tfsdk:"id"`
}

func (r *QualityProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + qualityProfileResourceName
}

func (r *QualityProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->Quality Profile resource.\nFor more information refer to [Quality Profile](https://wiki.servarr.com/readarr/settings#quality-profiles) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Quality Profile ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Quality Profile Name.",
				Required:            true,
			},
			"upgrade_allowed": schema.BoolAttribute{
				MarkdownDescription: "Upgrade allowed flag.",
				Optional:            true,
				Computed:            true,
			},
			"cutoff": schema.Int64Attribute{
				MarkdownDescription: "Quality ID to which cutoff.",
				Optional:            true,
				Computed:            true,
			},
			"quality_groups": schema.SetNestedAttribute{
				MarkdownDescription: "Quality groups.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.getQualityGroupSchema().Attributes,
				},
			},
		},
	}
}

func (r QualityProfileResource) getQualityGroupSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Quality group ID.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Quality group name.",
				Optional:            true,
				Computed:            true,
			},
			"qualities": schema.SetNestedAttribute{
				MarkdownDescription: "Qualities in group.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.getQualitySchema().Attributes,
				},
			},
		},
	}
}

func (r QualityProfileResource) getQualitySchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Quality ID.",
				Optional:            true,
				Computed:            true,
				// plan on uptate is unknown for 1 item array
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Quality name.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *QualityProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *QualityProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var profile *QualityProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	request := profile.read(ctx)

	// Create new QualityProfile
	response, _, err := r.client.QualityProfileApi.CreateQualityProfile(ctx).QualityProfileResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+qualityProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *QualityProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var profile *QualityProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get qualityprofile current value
	response, _, err := r.client.QualityProfileApi.GetQualityProfileById(ctx, int32(profile.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+qualityProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *QualityProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var profile *QualityProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	request := profile.read(ctx)

	// Update QualityProfile
	response, _, err := r.client.QualityProfileApi.UpdateQualityProfile(ctx, strconv.Itoa(int(request.GetId()))).QualityProfileResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+qualityProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *QualityProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var profile *QualityProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete qualityprofile current value
	_, err := r.client.QualityProfileApi.DeleteQualityProfile(ctx, int32(profile.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+qualityProfileResourceName+": "+strconv.Itoa(int(profile.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *QualityProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+qualityProfileResourceName+": "+req.ID)
}

func (p *QualityProfile) write(ctx context.Context, profile *readarr.QualityProfileResource) {
	p.UpgradeAllowed = types.BoolValue(profile.GetUpgradeAllowed())
	p.ID = types.Int64Value(int64(profile.GetId()))
	p.Name = types.StringValue(profile.GetName())
	p.Cutoff = types.Int64Value(int64(profile.GetCutoff()))
	p.QualityGroups = types.SetValueMust(QualityProfileResource{}.getQualityGroupSchema().Type(), nil)

	qualityGroups := make([]QualityGroup, len(profile.GetItems()))
	for n, g := range profile.GetItems() {
		qualityGroups[n].write(ctx, g)
	}

	tfsdk.ValueFrom(ctx, qualityGroups, p.QualityGroups.Type(ctx), &p.QualityGroups)
}

func (q *QualityGroup) write(ctx context.Context, group *readarr.QualityProfileQualityItemResource) {
	var (
		name      string
		id        int64
		qualities []Quality
	)

	if len(group.GetItems()) == 0 {
		qualities = []Quality{{
			ID:   types.Int64Value(int64(group.Quality.GetId())),
			Name: types.StringValue(group.Quality.GetName()),
		}}
	} else {
		name = group.GetName()
		id = int64(group.GetId())
		qualities = make([]Quality, len(group.GetItems()))
		for m, q := range group.GetItems() {
			qualities[m].write(q)
		}
	}

	q.Name = types.StringValue(name)
	q.ID = types.Int64Value(id)
	q.Qualities = types.SetValueMust(QualityProfileResource{}.getQualitySchema().Type(), nil)

	tfsdk.ValueFrom(ctx, qualities, q.Qualities.Type(ctx), &q.Qualities)
}

func (q *Quality) write(quality *readarr.QualityProfileQualityItemResource) {
	q.ID = types.Int64Value(int64(quality.Quality.GetId()))
	q.Name = types.StringValue(quality.Quality.GetName())
}

func (p *QualityProfile) read(ctx context.Context) *readarr.QualityProfileResource {
	groups := make([]QualityGroup, len(p.QualityGroups.Elements()))
	tfsdk.ValueAs(ctx, p.QualityGroups, &groups)
	qualities := make([]*readarr.QualityProfileQualityItemResource, len(groups))

	for n, g := range groups {
		q := make([]Quality, len(g.Qualities.Elements()))
		tfsdk.ValueAs(ctx, g.Qualities, &q)

		if len(q) == 1 {
			quality := readarr.NewQuality()
			quality.SetId(int32(q[0].ID.ValueInt64()))
			quality.SetName(q[0].Name.ValueString())

			item := readarr.NewQualityProfileQualityItemResource()
			item.SetAllowed(true)
			item.SetQuality(*quality)

			qualities[n] = item

			continue
		}

		items := make([]*readarr.QualityProfileQualityItemResource, len(q))
		for m, q := range q {
			items[m] = q.read()
		}

		quality := readarr.NewQualityProfileQualityItemResource()
		quality.SetId(int32(g.ID.ValueInt64()))
		quality.SetName(g.Name.ValueString())
		quality.SetAllowed(true)
		quality.SetItems(items)
		qualities[n] = quality
	}

	profile := readarr.NewQualityProfileResource()
	profile.SetUpgradeAllowed(p.UpgradeAllowed.ValueBool())
	profile.SetId(int32(p.ID.ValueInt64()))
	profile.SetCutoff(int32(p.Cutoff.ValueInt64()))
	profile.SetName(p.Name.ValueString())
	profile.SetItems(qualities)

	return profile
}

func (q *Quality) read() *readarr.QualityProfileQualityItemResource {
	quality := readarr.NewQuality()
	quality.SetName(q.Name.ValueString())
	quality.SetId(int32(q.ID.ValueInt64()))

	item := readarr.NewQualityProfileQualityItemResource()
	item.SetAllowed(true)
	item.SetQuality(*quality)

	return item
}