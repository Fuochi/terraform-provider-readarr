package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	notificationGoodreadsBookshelvesResourceName   = "notification_goodreads_bookshelves"
	notificationGoodreadsBookshelvesImplementation = "GoodreadsBookshelf"
	notificationGoodreadsBookshelvesConfigContract = "GoodreadsBookshelfNotificationSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationGoodreadsBookshelvesResource{}
	_ resource.ResourceWithImportState = &NotificationGoodreadsBookshelvesResource{}
)

func NewNotificationGoodreadsBookshelvesResource() resource.Resource {
	return &NotificationGoodreadsBookshelvesResource{}
}

// NotificationGoodreadsBookshelvesResource defines the notification implementation.
type NotificationGoodreadsBookshelvesResource struct {
	client *readarr.APIClient
}

// NotificationGoodreadsBookshelves describes the notification data model.
type NotificationGoodreadsBookshelves struct {
	AddIds                     types.Set    `tfsdk:"add_ids"`
	RemoveIds                  types.Set    `tfsdk:"remove_ids"`
	Tags                       types.Set    `tfsdk:"tags"`
	AccessToken                types.String `tfsdk:"access_token"`
	AccessTokenSecret          types.String `tfsdk:"access_token_secret"`
	RequestTokenSecret         types.String `tfsdk:"request_token_secret"`
	Name                       types.String `tfsdk:"name"`
	Username                   types.String `tfsdk:"username"`
	UserID                     types.String `tfsdk:"user_id"`
	ID                         types.Int64  `tfsdk:"id"`
	IncludeHealthWarnings      types.Bool   `tfsdk:"include_health_warnings"`
	OnUpgrade                  types.Bool   `tfsdk:"on_upgrade"`
	OnReleaseImport            types.Bool   `tfsdk:"on_release_import"`
	OnAuthorDelete             types.Bool   `tfsdk:"on_author_delete"`
	OnBookDelete               types.Bool   `tfsdk:"on_book_delete"`
	OnBookFileDelete           types.Bool   `tfsdk:"on_book_file_delete"`
	OnBookFileDeleteForUpgrade types.Bool   `tfsdk:"on_book_file_delete_for_upgrade"`
}

func (n NotificationGoodreadsBookshelves) toNotification() *Notification {
	return &Notification{
		Tags:                       n.Tags,
		AddIds:                     n.AddIds,
		RemoveIds:                  n.RemoveIds,
		AccessToken:                n.AccessToken,
		AccessTokenSecret:          n.AccessTokenSecret,
		RequestTokenSecret:         n.RequestTokenSecret,
		Username:                   n.Username,
		UserID:                     n.UserID,
		Name:                       n.Name,
		ID:                         n.ID,
		OnReleaseImport:            n.OnReleaseImport,
		OnAuthorDelete:             n.OnAuthorDelete,
		IncludeHealthWarnings:      n.IncludeHealthWarnings,
		OnBookDelete:               n.OnBookDelete,
		OnBookFileDelete:           n.OnBookFileDelete,
		OnUpgrade:                  n.OnUpgrade,
		OnBookFileDeleteForUpgrade: n.OnBookFileDeleteForUpgrade,
		Implementation:             types.StringValue(notificationGoodreadsBookshelvesImplementation),
		ConfigContract:             types.StringValue(notificationGoodreadsBookshelvesConfigContract),
	}
}

func (n *NotificationGoodreadsBookshelves) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.AddIds = notification.AddIds
	n.RemoveIds = notification.RemoveIds
	n.AccessToken = notification.AccessToken
	n.AccessTokenSecret = notification.AccessTokenSecret
	n.RequestTokenSecret = notification.RequestTokenSecret
	n.Username = notification.Username
	n.UserID = notification.UserID
	n.Name = notification.Name
	n.ID = notification.ID
	n.OnBookFileDeleteForUpgrade = notification.OnBookFileDeleteForUpgrade
	n.OnBookFileDelete = notification.OnBookFileDelete
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnAuthorDelete = notification.OnAuthorDelete
	n.OnBookDelete = notification.OnBookDelete
	n.OnUpgrade = notification.OnUpgrade
	n.OnReleaseImport = notification.OnReleaseImport
}

func (r *NotificationGoodreadsBookshelvesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationGoodreadsBookshelvesResourceName
}

func (r *NotificationGoodreadsBookshelvesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification GoodreadsBookshelves resource.\nFor more information refer to [Notification](https://wiki.servarr.com/readarr/settings#connect) and [GoodreadsBookshelves](https://wiki.servarr.com/readarr/supported#goodreadsbookshelf).",
		Attributes: map[string]schema.Attribute{
			"on_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On upgrade flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_author_delete": schema.BoolAttribute{
				MarkdownDescription: "On author deleted flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_book_delete": schema.BoolAttribute{
				MarkdownDescription: "On book delete flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_book_file_delete": schema.BoolAttribute{
				MarkdownDescription: "On book file delete flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_book_file_delete_for_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On book file delete for upgrade flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_release_import": schema.BoolAttribute{
				MarkdownDescription: "On release import flag.",
				Optional:            true,
				Computed:            true,
			},
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Notification name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Notification ID.",
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
			"add_ids": schema.SetAttribute{
				MarkdownDescription: "Add IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"remove_ids": schema.SetAttribute{
				MarkdownDescription: "Remove IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *NotificationGoodreadsBookshelvesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationGoodreadsBookshelvesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationGoodreadsBookshelves

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationGoodreadsBookshelves
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationGoodreadsBookshelvesResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationGoodreadsBookshelvesResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationGoodreadsBookshelvesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationGoodreadsBookshelves

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationGoodreadsBookshelves current value
	response, _, err := r.client.NotificationAPI.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationGoodreadsBookshelvesResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationGoodreadsBookshelvesResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationGoodreadsBookshelvesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationGoodreadsBookshelves

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationGoodreadsBookshelves
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationGoodreadsBookshelvesResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationGoodreadsBookshelvesResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationGoodreadsBookshelvesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationGoodreadsBookshelves current value
	_, err := r.client.NotificationAPI.DeleteNotification(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationGoodreadsBookshelvesResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationGoodreadsBookshelvesResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationGoodreadsBookshelvesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationGoodreadsBookshelvesResourceName+": "+req.ID)
}

func (n *NotificationGoodreadsBookshelves) write(ctx context.Context, notification *readarr.NotificationResource, diags *diag.Diagnostics) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification, diags)
	n.fromNotification(genericNotification)
}

func (n *NotificationGoodreadsBookshelves) read(ctx context.Context, diags *diag.Diagnostics) *readarr.NotificationResource {
	return n.toNotification().read(ctx, diags)
}
