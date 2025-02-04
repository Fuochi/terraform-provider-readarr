package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	notificationGoodreadsOwnedBooksResourceName   = "notification_goodreads_owned_books"
	notificationGoodreadsOwnedBooksImplementation = "GoodreadsOwnedBooks"
	notificationGoodreadsOwnedBooksConfigContract = "GoodreadsOwnedBooksNotificationSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationGoodreadsOwnedBooksResource{}
	_ resource.ResourceWithImportState = &NotificationGoodreadsOwnedBooksResource{}
)

func NewNotificationGoodreadsOwnedBooksResource() resource.Resource {
	return &NotificationGoodreadsOwnedBooksResource{}
}

// NotificationGoodreadsOwnedBooksResource defines the notification implementation.
type NotificationGoodreadsOwnedBooksResource struct {
	client *readarr.APIClient
}

// NotificationGoodreadsOwnedBooks describes the notification data model.
type NotificationGoodreadsOwnedBooks struct {
	Tags               types.Set    `tfsdk:"tags"`
	AccessToken        types.String `tfsdk:"access_token"`
	AccessTokenSecret  types.String `tfsdk:"access_token_secret"`
	RequestTokenSecret types.String `tfsdk:"request_token_secret"`
	Name               types.String `tfsdk:"name"`
	Username           types.String `tfsdk:"username"`
	UserID             types.String `tfsdk:"user_id"`
	Description        types.String `tfsdk:"description"`
	Location           types.String `tfsdk:"location"`
	ID                 types.Int64  `tfsdk:"id"`
	Condition          types.Int64  `tfsdk:"condition"`
	OnUpgrade          types.Bool   `tfsdk:"on_upgrade"`
	OnReleaseImport    types.Bool   `tfsdk:"on_release_import"`
}

func (n NotificationGoodreadsOwnedBooks) toNotification() *Notification {
	return &Notification{
		Tags:               n.Tags,
		AccessToken:        n.AccessToken,
		AccessTokenSecret:  n.AccessTokenSecret,
		RequestTokenSecret: n.RequestTokenSecret,
		Username:           n.Username,
		UserID:             n.UserID,
		Description:        n.Description,
		Location:           n.Location,
		Name:               n.Name,
		ID:                 n.ID,
		Condition:          n.Condition,
		OnReleaseImport:    n.OnReleaseImport,
		OnUpgrade:          n.OnUpgrade,
		Implementation:     types.StringValue(notificationGoodreadsOwnedBooksImplementation),
		ConfigContract:     types.StringValue(notificationGoodreadsOwnedBooksConfigContract),
	}
}

func (n *NotificationGoodreadsOwnedBooks) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.AccessToken = notification.AccessToken
	n.AccessTokenSecret = notification.AccessTokenSecret
	n.RequestTokenSecret = notification.RequestTokenSecret
	n.Username = notification.Username
	n.UserID = notification.UserID
	n.Description = notification.Description
	n.Location = notification.Location
	n.Name = notification.Name
	n.ID = notification.ID
	n.Condition = notification.Condition
	n.OnUpgrade = notification.OnUpgrade
	n.OnReleaseImport = notification.OnReleaseImport
}

func (r *NotificationGoodreadsOwnedBooksResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationGoodreadsOwnedBooksResourceName
}

func (r *NotificationGoodreadsOwnedBooksResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification GoodreadsOwnedBooks resource.\nFor more information refer to [Notification](https://wiki.servarr.com/readarr/settings#connect) and [GoodreadsOwnedBooks](https://wiki.servarr.com/readarr/supported#goodreadsownedbooks).",
		Attributes: map[string]schema.Attribute{
			"on_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On upgrade flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_release_import": schema.BoolAttribute{
				MarkdownDescription: "On release import flag.",
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
			"condition": schema.Int64Attribute{
				MarkdownDescription: "Condition. `10` BrandNew, `20` LikeNew, `30` VeryGood, `40` Good, `50` Acceptable, `60` Poor.",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(10, 20, 30, 40, 50, 60),
				},
			},
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
			"description": schema.StringAttribute{
				MarkdownDescription: "Condition description.",
				Optional:            true,
				Computed:            true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "Purchase location.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *NotificationGoodreadsOwnedBooksResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationGoodreadsOwnedBooksResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationGoodreadsOwnedBooks

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationGoodreadsOwnedBooks
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationGoodreadsOwnedBooksResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationGoodreadsOwnedBooksResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationGoodreadsOwnedBooksResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationGoodreadsOwnedBooks

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationGoodreadsOwnedBooks current value
	response, _, err := r.client.NotificationAPI.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationGoodreadsOwnedBooksResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationGoodreadsOwnedBooksResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationGoodreadsOwnedBooksResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationGoodreadsOwnedBooks

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationGoodreadsOwnedBooks
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationGoodreadsOwnedBooksResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationGoodreadsOwnedBooksResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationGoodreadsOwnedBooksResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationGoodreadsOwnedBooks current value
	_, err := r.client.NotificationAPI.DeleteNotification(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationGoodreadsOwnedBooksResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationGoodreadsOwnedBooksResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationGoodreadsOwnedBooksResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationGoodreadsOwnedBooksResourceName+": "+req.ID)
}

func (n *NotificationGoodreadsOwnedBooks) write(ctx context.Context, notification *readarr.NotificationResource, diags *diag.Diagnostics) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification, diags)
	n.fromNotification(genericNotification)
}

func (n *NotificationGoodreadsOwnedBooks) read(ctx context.Context, diags *diag.Diagnostics) *readarr.NotificationResource {
	return n.toNotification().read(ctx, diags)
}
