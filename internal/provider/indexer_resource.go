package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const indexerResourceName = "indexer"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerResource{}
	_ resource.ResourceWithImportState = &IndexerResource{}
)

var indexerFields = helpers.Fields{
	IntSlices:        []string{"categories"},
	Bools:            []string{"rankedOnly", "allowZeroSize"},
	Ints:             []string{"earlyReleaseLimit", "delay", "minimumSeeders", "seedTime", "discographySeedTime"},
	IntsExceptions:   []string{"seedCriteria.seedTime", "seedCriteria.discographySeedTime"},
	Strings:          []string{"apiKey", "apiPath", "baseUrl", "username", "passkey", "password", "additionalParameters", "captchaToken", "cookie"},
	Floats:           []string{"seedRatio"},
	FloatsExceptions: []string{"seedCriteria.seedRatio"},
}

func NewIndexerResource() resource.Resource {
	return &IndexerResource{}
}

// IndexerResource defines the indexer implementation.
type IndexerResource struct {
	client *readarr.APIClient
}

// Indexer describes the indexer data model.
type Indexer struct {
	Tags                    types.Set     `tfsdk:"tags"`
	Categories              types.Set     `tfsdk:"categories"`
	APIUser                 types.String  `tfsdk:"api_user"`
	AdditionalParameters    types.String  `tfsdk:"additional_parameters"`
	Name                    types.String  `tfsdk:"name"`
	Implementation          types.String  `tfsdk:"implementation"`
	CaptchaToken            types.String  `tfsdk:"captcha_token"`
	Protocol                types.String  `tfsdk:"protocol"`
	ConfigContract          types.String  `tfsdk:"config_contract"`
	APIKey                  types.String  `tfsdk:"api_key"`
	APIPath                 types.String  `tfsdk:"api_path"`
	Cookie                  types.String  `tfsdk:"cookie"`
	BaseURL                 types.String  `tfsdk:"base_url"`
	Username                types.String  `tfsdk:"username"`
	Password                types.String  `tfsdk:"password"`
	Passkey                 types.String  `tfsdk:"passkey"`
	EarlyReleaseLimit       types.Int64   `tfsdk:"early_release_limit"`
	SeedTime                types.Int64   `tfsdk:"seed_time"`
	Delay                   types.Int64   `tfsdk:"delay"`
	MinimumSeeders          types.Int64   `tfsdk:"minimum_seeders"`
	ID                      types.Int64   `tfsdk:"id"`
	SeedRatio               types.Float64 `tfsdk:"seed_ratio"`
	Priority                types.Int64   `tfsdk:"priority"`
	DiscographySeedTime     types.Int64   `tfsdk:"author_seed_time"`
	EnableInteractiveSearch types.Bool    `tfsdk:"enable_interactive_search"`
	EnableRss               types.Bool    `tfsdk:"enable_rss"`
	EnableAutomaticSearch   types.Bool    `tfsdk:"enable_automatic_search"`
	AllowZeroSize           types.Bool    `tfsdk:"allow_zero_size"`
	RankedOnly              types.Bool    `tfsdk:"ranked_only"`
}

func (r *IndexerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerResourceName
}

func (r *IndexerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->Generic Indexer resource. When possible use a specific resource instead.\nFor more information refer to [Indexer](https://wiki.servarr.com/readarr/settings#indexers) documentation.",
		Attributes: map[string]schema.Attribute{
			"enable_automatic_search": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic search flag.",
				Optional:            true,
				Computed:            true,
			},
			"enable_interactive_search": schema.BoolAttribute{
				MarkdownDescription: "Enable interactive search flag.",
				Optional:            true,
				Computed:            true,
			},
			"enable_rss": schema.BoolAttribute{
				MarkdownDescription: "Enable RSS flag.",
				Optional:            true,
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.",
				Optional:            true,
				Computed:            true,
			},
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "Indexer configuration template.",
				Required:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "Indexer implementation name.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Indexer name.",
				Required:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Protocol. Valid values are 'usenet' and 'torrent'.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("usenet", "torrent"),
				},
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Indexer ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"allow_zero_size": schema.BoolAttribute{
				MarkdownDescription: "Allow zero size files.",
				Optional:            true,
				Computed:            true,
			},
			"ranked_only": schema.BoolAttribute{
				MarkdownDescription: "Allow ranked only.",
				Optional:            true,
				Computed:            true,
			},
			"delay": schema.Int64Attribute{
				MarkdownDescription: "Delay before grabbing.",
				Optional:            true,
				Computed:            true,
			},
			"minimum_seeders": schema.Int64Attribute{
				MarkdownDescription: "Minimum seeders.",
				Optional:            true,
				Computed:            true,
			},
			"early_release_limit": schema.Int64Attribute{
				MarkdownDescription: "Early release limit.",
				Optional:            true,
				Computed:            true,
			},
			"seed_time": schema.Int64Attribute{
				MarkdownDescription: "Seed time.",
				Optional:            true,
				Computed:            true,
			},
			"author_seed_time": schema.Int64Attribute{
				MarkdownDescription: "Author seed time.",
				Optional:            true,
				Computed:            true,
			},
			"seed_ratio": schema.Float64Attribute{
				MarkdownDescription: "Seed ratio.",
				Optional:            true,
				Computed:            true,
			},
			"additional_parameters": schema.StringAttribute{
				MarkdownDescription: "Additional parameters.",
				Optional:            true,
				Computed:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
			},
			"api_user": schema.StringAttribute{
				MarkdownDescription: "API User.",
				Optional:            true,
				Computed:            true,
			},
			"api_path": schema.StringAttribute{
				MarkdownDescription: "API path.",
				Optional:            true,
				Computed:            true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
			},
			"captcha_token": schema.StringAttribute{
				MarkdownDescription: "Captcha token.",
				Optional:            true,
				Computed:            true,
			},
			"cookie": schema.StringAttribute{
				MarkdownDescription: "Cookie.",
				Optional:            true,
				Computed:            true,
			},
			"passkey": schema.StringAttribute{
				MarkdownDescription: "Passkey.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"categories": schema.SetAttribute{
				MarkdownDescription: "Series list.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *IndexerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *IndexerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var indexer *Indexer

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Indexer
	request := indexer.read(ctx)

	response, _, err := r.client.IndexerApi.CreateIndexer(ctx).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct.
	// this is needed because of many empty fields are unknown in both plan and read
	var state Indexer

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IndexerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var indexer *Indexer

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get Indexer current value
	response, _, err := r.client.IndexerApi.GetIndexerById(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct.
	// this is needed because of many empty fields are unknown in both plan and read
	var state Indexer

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IndexerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var indexer *Indexer

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Indexer
	request := indexer.read(ctx)

	response, _, err := r.client.IndexerApi.UpdateIndexer(ctx, strconv.Itoa(int(request.GetId()))).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, indexerResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct.
	// this is needed because of many empty fields are unknown in both plan and read
	var state Indexer

	state.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *IndexerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var indexer Indexer

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete Indexer current value
	_, err := r.client.IndexerApi.DeleteIndexer(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerResourceName+": "+strconv.Itoa(int(indexer.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerResourceName+": "+req.ID)
}

func (i *Indexer) write(ctx context.Context, indexer *readarr.IndexerResource) {
	i.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, indexer.Tags)
	i.EnableAutomaticSearch = types.BoolValue(indexer.GetEnableAutomaticSearch())
	i.EnableInteractiveSearch = types.BoolValue(indexer.GetEnableInteractiveSearch())
	i.EnableRss = types.BoolValue(indexer.GetEnableRss())
	i.Priority = types.Int64Value(int64(indexer.GetPriority()))
	i.ID = types.Int64Value(int64(indexer.GetId()))
	i.ConfigContract = types.StringValue(indexer.GetConfigContract())
	i.Implementation = types.StringValue(indexer.GetImplementation())
	i.Name = types.StringValue(indexer.GetName())
	i.Protocol = types.StringValue(string(indexer.GetProtocol()))
	i.Categories = types.SetValueMust(types.Int64Type, nil)
	helpers.WriteFields(ctx, i, indexer.GetFields(), indexerFields)
}

func (i *Indexer) read(ctx context.Context) *readarr.IndexerResource {
	tags := make([]*int32, len(i.Tags.Elements()))
	tfsdk.ValueAs(ctx, i.Tags, &tags)

	indexer := readarr.NewIndexerResource()
	indexer.SetEnableAutomaticSearch(i.EnableAutomaticSearch.ValueBool())
	indexer.SetEnableInteractiveSearch(i.EnableInteractiveSearch.ValueBool())
	indexer.SetEnableRss(i.EnableRss.ValueBool())
	indexer.SetPriority(int32(i.Priority.ValueInt64()))
	indexer.SetId(int32(i.ID.ValueInt64()))
	indexer.SetConfigContract(i.ConfigContract.ValueString())
	indexer.SetImplementation(i.Implementation.ValueString())
	indexer.SetName(i.Name.ValueString())
	indexer.SetProtocol(readarr.DownloadProtocol(i.Protocol.ValueString()))
	indexer.SetTags(tags)
	indexer.SetFields(helpers.ReadFields(ctx, i, indexerFields))

	return indexer
}
