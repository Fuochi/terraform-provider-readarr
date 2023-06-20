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
	importListGoodreadsOwnedBooksResourceName   = "import_list_goodreads_owned_books"
	importListGoodreadsOwnedBooksImplementation = "GoodreadsOwnedBooks"
	importListGoodreadsOwnedBooksConfigContract = "GoodreadsOwnedBooksImportListSettings"
	importListGoodreadsOwnedBooksType           = "goodreads"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListGoodreadsOwnedBooksResource{}
	_ resource.ResourceWithImportState = &ImportListGoodreadsOwnedBooksResource{}
)

func NewImportListGoodreadsOwnedBooksResource() resource.Resource {
	return &ImportListGoodreadsOwnedBooksResource{}
}

// ImportListGoodreadsOwnedBooksResource defines the import list implementation.
type ImportListGoodreadsOwnedBooksResource struct {
	client *readarr.APIClient
}

// ImportListGoodreadsOwnedBooks describes the import list data model.
type ImportListGoodreadsOwnedBooks struct {
	Tags types.Set    `tfsdk:"tags"`
	Name types.String `tfsdk:"name"`
	// MonitorNewItems       types.String `tfsdk:"monitor_new_items"`
	ShouldMonitor      types.String `tfsdk:"should_monitor"`
	RootFolderPath     types.String `tfsdk:"root_folder_path"`
	Username           types.String `tfsdk:"username"`
	UserID             types.String `tfsdk:"user_id"`
	AccessToken        types.String `tfsdk:"access_token"`
	AccessTokenSecret  types.String `tfsdk:"access_token_secret"`
	RequestTokenSecret types.String `tfsdk:"request_token_secret"`
	QualityProfileID   types.Int64  `tfsdk:"quality_profile_id"`
	MetadataProfileID  types.Int64  `tfsdk:"metadata_profile_id"`
	ListOrder          types.Int64  `tfsdk:"list_order"`
	ID                 types.Int64  `tfsdk:"id"`
	EnableAutomaticAdd types.Bool   `tfsdk:"enable_automatic_add"`
	// ShouldMonitorExisting types.Bool   `tfsdk:"should_monitor_existing"`
	ShouldSearch types.Bool `tfsdk:"should_search"`
}

func (i ImportListGoodreadsOwnedBooks) toImportList() *ImportList {
	return &ImportList{
		Tags: i.Tags,
		Name: i.Name,
		// MonitorNewItems:       i.MonitorNewItems,
		ShouldMonitor:      i.ShouldMonitor,
		RootFolderPath:     i.RootFolderPath,
		Username:           i.Username,
		UserID:             i.UserID,
		AccessToken:        i.AccessToken,
		AccessTokenSecret:  i.AccessTokenSecret,
		RequestTokenSecret: i.RequestTokenSecret,
		QualityProfileID:   i.QualityProfileID,
		MetadataProfileID:  i.MetadataProfileID,
		ListOrder:          i.ListOrder,
		ID:                 i.ID,
		EnableAutomaticAdd: i.EnableAutomaticAdd,
		// ShouldMonitorExisting: i.ShouldMonitorExisting,
		ShouldSearch:   i.ShouldSearch,
		Implementation: types.StringValue(importListGoodreadsOwnedBooksImplementation),
		ConfigContract: types.StringValue(importListGoodreadsOwnedBooksConfigContract),
		ListType:       types.StringValue(importListGoodreadsOwnedBooksType),
	}
}

func (i *ImportListGoodreadsOwnedBooks) fromImportList(importList *ImportList) {
	i.Tags = importList.Tags
	i.Name = importList.Name
	// i.MonitorNewItems = importList.MonitorNewItems
	i.ShouldMonitor = importList.ShouldMonitor
	i.RootFolderPath = importList.RootFolderPath
	i.Username = importList.Username
	i.UserID = importList.UserID
	i.AccessToken = importList.AccessToken
	i.AccessTokenSecret = importList.AccessTokenSecret
	i.RequestTokenSecret = importList.RequestTokenSecret
	i.QualityProfileID = importList.QualityProfileID
	i.MetadataProfileID = importList.MetadataProfileID
	i.ListOrder = importList.ListOrder
	i.ID = importList.ID
	i.EnableAutomaticAdd = importList.EnableAutomaticAdd
	// i.ShouldMonitorExisting = importList.ShouldMonitorExisting
	i.ShouldSearch = importList.ShouldSearch
}

func (r *ImportListGoodreadsOwnedBooksResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListGoodreadsOwnedBooksResourceName
}

func (r *ImportListGoodreadsOwnedBooksResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Import List Goodreads OwnedBooks resource.\nFor more information refer to [Import List](https://wiki.servarr.com/readarr/settings#import-lists) and [Goodreads OwnedBooks](https://wiki.servarr.com/readarr/supported#goodreadsownedbooks).",
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
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token.",
				Required:            true,
				Sensitive:           true,
			},
			"access_token_secret": schema.StringAttribute{
				MarkdownDescription: "Access token secret.",
				Required:            true,
				Sensitive:           true,
			},
			"request_token_secret": schema.StringAttribute{
				MarkdownDescription: "Request token secret.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Required:            true,
			},
			"user_id": schema.StringAttribute{
				MarkdownDescription: "User ID.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *ImportListGoodreadsOwnedBooksResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ImportListGoodreadsOwnedBooksResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListGoodreadsOwnedBooks

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListGoodreadsOwnedBooks
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListGoodreadsOwnedBooksResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListGoodreadsOwnedBooksResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListGoodreadsOwnedBooksResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportListGoodreadsOwnedBooks

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportListGoodreadsOwnedBooks current value
	response, _, err := r.client.ImportListApi.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListGoodreadsOwnedBooksResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListGoodreadsOwnedBooksResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListGoodreadsOwnedBooksResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportListGoodreadsOwnedBooks

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListGoodreadsOwnedBooks
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListGoodreadsOwnedBooksResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListGoodreadsOwnedBooksResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListGoodreadsOwnedBooksResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var importList *ImportListGoodreadsOwnedBooks

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListGoodreadsOwnedBooks current value
	_, err := r.client.ImportListApi.DeleteImportList(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, importListGoodreadsOwnedBooksResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListGoodreadsOwnedBooksResourceName+": "+strconv.Itoa(int(importList.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListGoodreadsOwnedBooksResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListGoodreadsOwnedBooksResourceName+": "+req.ID)
}

func (i *ImportListGoodreadsOwnedBooks) write(ctx context.Context, importList *readarr.ImportListResource) {
	genericImportList := i.toImportList()
	genericImportList.write(ctx, importList)
	i.fromImportList(genericImportList)
}

func (i *ImportListGoodreadsOwnedBooks) read(ctx context.Context) *readarr.ImportListResource {
	return i.toImportList().read(ctx)
}
