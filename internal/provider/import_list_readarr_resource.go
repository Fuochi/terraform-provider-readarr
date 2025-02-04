package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	importListReadarrResourceName   = "import_list_readarr"
	importListReadarrImplementation = "ReadarrImport"
	importListReadarrConfigContract = "ReadarrSettings"
	importListReadarrType           = "program"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListReadarrResource{}
	_ resource.ResourceWithImportState = &ImportListReadarrResource{}
)

func NewImportListReadarrResource() resource.Resource {
	return &ImportListReadarrResource{}
}

// ImportListReadarrResource defines the import list implementation.
type ImportListReadarrResource struct {
	client *readarr.APIClient
}

// ImportListReadarr describes the import list data model.
type ImportListReadarr struct {
	ProfileIds            types.Set    `tfsdk:"profile_ids"`
	TagIds                types.Set    `tfsdk:"tag_ids"`
	Tags                  types.Set    `tfsdk:"tags"`
	Name                  types.String `tfsdk:"name"`
	MonitorNewItems       types.String `tfsdk:"monitor_new_items"`
	ShouldMonitor         types.String `tfsdk:"should_monitor"`
	RootFolderPath        types.String `tfsdk:"root_folder_path"`
	BaseURL               types.String `tfsdk:"base_url"`
	APIKey                types.String `tfsdk:"api_key"`
	QualityProfileID      types.Int64  `tfsdk:"quality_profile_id"`
	MetadataProfileID     types.Int64  `tfsdk:"metadata_profile_id"`
	ListOrder             types.Int64  `tfsdk:"list_order"`
	ID                    types.Int64  `tfsdk:"id"`
	EnableAutomaticAdd    types.Bool   `tfsdk:"enable_automatic_add"`
	ShouldMonitorExisting types.Bool   `tfsdk:"should_monitor_existing"`
	ShouldSearch          types.Bool   `tfsdk:"should_search"`
}

func (i ImportListReadarr) toImportList() *ImportList {
	return &ImportList{
		ProfileIds:            i.ProfileIds,
		TagIds:                i.TagIds,
		Tags:                  i.Tags,
		Name:                  i.Name,
		MonitorNewItems:       i.MonitorNewItems,
		ShouldMonitor:         i.ShouldMonitor,
		RootFolderPath:        i.RootFolderPath,
		BaseURL:               i.BaseURL,
		APIKey:                i.APIKey,
		QualityProfileID:      i.QualityProfileID,
		MetadataProfileID:     i.MetadataProfileID,
		ListOrder:             i.ListOrder,
		ID:                    i.ID,
		EnableAutomaticAdd:    i.EnableAutomaticAdd,
		ShouldMonitorExisting: i.ShouldMonitorExisting,
		ShouldSearch:          i.ShouldSearch,
		Implementation:        types.StringValue(importListReadarrImplementation),
		ConfigContract:        types.StringValue(importListReadarrConfigContract),
		ListType:              types.StringValue(importListReadarrType),
	}
}

func (i *ImportListReadarr) fromImportList(importList *ImportList) {
	i.ProfileIds = importList.ProfileIds
	i.TagIds = importList.TagIds
	i.Tags = importList.Tags
	i.Name = importList.Name
	i.MonitorNewItems = importList.MonitorNewItems
	i.ShouldMonitor = importList.ShouldMonitor
	i.RootFolderPath = importList.RootFolderPath
	i.BaseURL = importList.BaseURL
	i.APIKey = importList.APIKey
	i.QualityProfileID = importList.QualityProfileID
	i.MetadataProfileID = importList.MetadataProfileID
	i.ListOrder = importList.ListOrder
	i.ID = importList.ID
	i.EnableAutomaticAdd = importList.EnableAutomaticAdd
	i.ShouldMonitorExisting = importList.ShouldMonitorExisting
	i.ShouldSearch = importList.ShouldSearch
}

func (r *ImportListReadarrResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListReadarrResourceName
}

func (r *ImportListReadarrResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Import List Readarr resource.\nFor more information refer to [Import List](https://wiki.servarr.com/readarr/settings#import-lists) and [Readarr](https://wiki.servarr.com/readarr/supported#readarrimport).",
		Attributes: map[string]schema.Attribute{
			"enable_automatic_add": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic add flag.",
				Optional:            true,
				Computed:            true,
			},
			"should_monitor_existing": schema.BoolAttribute{
				MarkdownDescription: "Should monitor existing flag.",
				Optional:            true,
				Computed:            true,
			},
			"should_search": schema.BoolAttribute{
				MarkdownDescription: "Should search flag.",
				Optional:            true,
				Computed:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality profile ID.",
				Optional:            true,
				Computed:            true,
			},
			"metadata_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Metadata profile ID.",
				Optional:            true,
				Computed:            true,
			},
			"list_order": schema.Int64Attribute{
				MarkdownDescription: "List order.",
				Optional:            true,
				Computed:            true,
			},
			"root_folder_path": schema.StringAttribute{
				MarkdownDescription: "Root folder path.",
				Optional:            true,
				Computed:            true,
			},
			"should_monitor": schema.StringAttribute{
				MarkdownDescription: "Should monitor.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("none", "specificBook", "entireAuthor"),
				},
			},
			"monitor_new_items": schema.StringAttribute{
				MarkdownDescription: "Monitor new items.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("none", "all", "new"),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Import List name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Import List ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Required:            true,
				Sensitive:           true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Required:            true,
			},
			"profile_ids": schema.SetAttribute{
				MarkdownDescription: "Profile IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"tag_ids": schema.SetAttribute{
				MarkdownDescription: "Tag IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *ImportListReadarrResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ImportListReadarrResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListReadarr

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListReadarr
	request := importList.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ImportListAPI.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListReadarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListReadarrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListReadarrResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportListReadarr

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportListReadarr current value
	response, _, err := r.client.ImportListAPI.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListReadarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListReadarrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListReadarrResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportListReadarr

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListReadarr
	request := importList.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.ImportListAPI.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListReadarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListReadarrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListReadarrResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListReadarr current value
	_, err := r.client.ImportListAPI.DeleteImportList(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, importListReadarrResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListReadarrResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListReadarrResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListReadarrResourceName+": "+req.ID)
}

func (i *ImportListReadarr) write(ctx context.Context, importList *readarr.ImportListResource, diags *diag.Diagnostics) {
	genericImportList := i.toImportList()
	genericImportList.write(ctx, importList, diags)
	i.fromImportList(genericImportList)
}

func (i *ImportListReadarr) read(ctx context.Context, diags *diag.Diagnostics) *readarr.ImportListResource {
	return i.toImportList().read(ctx, diags)
}
