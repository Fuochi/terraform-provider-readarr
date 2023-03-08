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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	importListLazyLibrarianResourceName   = "import_list_lazy_librarian"
	importListLazyLibrarianImplementation = "LazyLibrarianImport"
	importListLazyLibrarianConfigContract = "LazyLibrarianImportSettings"
	importListLazyLibrarianType           = "other"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListLazyLibrarianResource{}
	_ resource.ResourceWithImportState = &ImportListLazyLibrarianResource{}
)

func NewImportListLazyLibrarianResource() resource.Resource {
	return &ImportListLazyLibrarianResource{}
}

// ImportListLazyLibrarianResource defines the import list implementation.
type ImportListLazyLibrarianResource struct {
	client *readarr.APIClient
}

// ImportListLazyLibrarian describes the import list data model.
type ImportListLazyLibrarian struct {
	Tags types.Set    `tfsdk:"tags"`
	Name types.String `tfsdk:"name"`
	// MonitorNewItems       types.String `tfsdk:"monitor_new_items"`
	BaseURL            types.String `tfsdk:"base_url"`
	APIKey             types.String `tfsdk:"api_key"`
	ShouldMonitor      types.String `tfsdk:"should_monitor"`
	RootFolderPath     types.String `tfsdk:"root_folder_path"`
	QualityProfileID   types.Int64  `tfsdk:"quality_profile_id"`
	MetadataProfileID  types.Int64  `tfsdk:"metadata_profile_id"`
	ListOrder          types.Int64  `tfsdk:"list_order"`
	ID                 types.Int64  `tfsdk:"id"`
	EnableAutomaticAdd types.Bool   `tfsdk:"enable_automatic_add"`
	// ShouldMonitorExisting types.Bool   `tfsdk:"should_monitor_existing"`
	ShouldSearch types.Bool `tfsdk:"should_search"`
}

func (i ImportListLazyLibrarian) toImportList() *ImportList {
	return &ImportList{
		Tags: i.Tags,
		Name: i.Name,
		// MonitorNewItems:       i.MonitorNewItems,
		ShouldMonitor:      i.ShouldMonitor,
		RootFolderPath:     i.RootFolderPath,
		BaseURL:            i.BaseURL,
		APIKey:             i.APIKey,
		QualityProfileID:   i.QualityProfileID,
		MetadataProfileID:  i.MetadataProfileID,
		ListOrder:          i.ListOrder,
		ID:                 i.ID,
		EnableAutomaticAdd: i.EnableAutomaticAdd,
		// ShouldMonitorExisting: i.ShouldMonitorExisting,
		ShouldSearch:   i.ShouldSearch,
		Implementation: types.StringValue(importListLazyLibrarianImplementation),
		ConfigContract: types.StringValue(importListLazyLibrarianConfigContract),
		ListType:       types.StringValue(importListLazyLibrarianType),
	}
}

func (i *ImportListLazyLibrarian) fromImportList(importList *ImportList) {
	i.Tags = importList.Tags
	i.Name = importList.Name
	// i.MonitorNewItems = importList.MonitorNewItems
	i.ShouldMonitor = importList.ShouldMonitor
	i.RootFolderPath = importList.RootFolderPath
	i.BaseURL = importList.BaseURL
	i.APIKey = importList.APIKey
	i.QualityProfileID = importList.QualityProfileID
	i.MetadataProfileID = importList.MetadataProfileID
	i.ListOrder = importList.ListOrder
	i.ID = importList.ID
	i.EnableAutomaticAdd = importList.EnableAutomaticAdd
	// i.ShouldMonitorExisting = importList.ShouldMonitorExisting
	i.ShouldSearch = importList.ShouldSearch
}

func (r *ImportListLazyLibrarianResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListLazyLibrarianResourceName
}

func (r *ImportListLazyLibrarianResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Import List Lazy Librarian resource.\nFor more information refer to [Import List](https://wiki.servarr.com/readarr/settings#import-lists) and [Lazy Librarian](https://wiki.servarr.com/readarr/supported#lazylibrarianimport).",
		Attributes: map[string]schema.Attribute{
			"enable_automatic_add": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic add flag.",
				Optional:            true,
				Computed:            true,
			},
			// "should_monitor_existing": schema.BoolAttribute{
			// 	MarkdownDescription: "Should monitor existing flag.",
			// 	Optional:            true,
			// 	Computed:            true,
			// },
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
			// "monitor_new_items": schema.StringAttribute{
			// 	MarkdownDescription: "Monitor new items.",
			// 	Optional:            true,
			// 	Computed:            true,
			// 	Validators: []validator.String{
			// 		stringvalidator.OneOf("none", "all", "new"),
			// 	},
			// },
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
		},
	}
}

func (r *ImportListLazyLibrarianResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ImportListLazyLibrarianResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListLazyLibrarian

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListLazyLibrarian
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListLazyLibrarianResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListLazyLibrarianResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListLazyLibrarianResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportListLazyLibrarian

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportListLazyLibrarian current value
	response, _, err := r.client.ImportListApi.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListLazyLibrarianResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListLazyLibrarianResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListLazyLibrarianResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportListLazyLibrarian

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListLazyLibrarian
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListLazyLibrarianResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListLazyLibrarianResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListLazyLibrarianResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var importList *ImportListLazyLibrarian

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListLazyLibrarian current value
	_, err := r.client.ImportListApi.DeleteImportList(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListLazyLibrarianResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListLazyLibrarianResourceName+": "+strconv.Itoa(int(importList.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListLazyLibrarianResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListLazyLibrarianResourceName+": "+req.ID)
}

func (i *ImportListLazyLibrarian) write(ctx context.Context, importList *readarr.ImportListResource) {
	genericImportList := i.toImportList()
	genericImportList.write(ctx, importList)
	i.fromImportList(genericImportList)
}

func (i *ImportListLazyLibrarian) read(ctx context.Context) *readarr.ImportListResource {
	return i.toImportList().read(ctx)
}
