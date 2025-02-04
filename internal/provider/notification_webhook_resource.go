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
	notificationWebhookResourceName   = "notification_webhook"
	notificationWebhookImplementation = "Webhook"
	notificationWebhookConfigContract = "WebhookSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationWebhookResource{}
	_ resource.ResourceWithImportState = &NotificationWebhookResource{}
)

func NewNotificationWebhookResource() resource.Resource {
	return &NotificationWebhookResource{}
}

// NotificationWebhookResource defines the notification implementation.
type NotificationWebhookResource struct {
	client *readarr.APIClient
}

// NotificationWebhook describes the notification data model.
type NotificationWebhook struct {
	Tags                       types.Set    `tfsdk:"tags"`
	URL                        types.String `tfsdk:"url"`
	Name                       types.String `tfsdk:"name"`
	Username                   types.String `tfsdk:"username"`
	Password                   types.String `tfsdk:"password"`
	ID                         types.Int64  `tfsdk:"id"`
	Method                     types.Int64  `tfsdk:"method"`
	OnGrab                     types.Bool   `tfsdk:"on_grab"`
	IncludeHealthWarnings      types.Bool   `tfsdk:"include_health_warnings"`
	OnHealthIssue              types.Bool   `tfsdk:"on_health_issue"`
	OnApplicationUpdate        types.Bool   `tfsdk:"on_application_update"`
	OnRename                   types.Bool   `tfsdk:"on_rename"`
	OnUpgrade                  types.Bool   `tfsdk:"on_upgrade"`
	OnReleaseImport            types.Bool   `tfsdk:"on_release_import"`
	OnAuthorDelete             types.Bool   `tfsdk:"on_author_delete"`
	OnBookDelete               types.Bool   `tfsdk:"on_book_delete"`
	OnBookFileDelete           types.Bool   `tfsdk:"on_book_file_delete"`
	OnBookFileDeleteForUpgrade types.Bool   `tfsdk:"on_book_file_delete_for_upgrade"`
	OnBookRetag                types.Bool   `tfsdk:"on_book_retag"`
}

func (n NotificationWebhook) toNotification() *Notification {
	return &Notification{
		Tags:                       n.Tags,
		URL:                        n.URL,
		Method:                     n.Method,
		Username:                   n.Username,
		Password:                   n.Password,
		Name:                       n.Name,
		ID:                         n.ID,
		OnGrab:                     n.OnGrab,
		OnReleaseImport:            n.OnReleaseImport,
		OnAuthorDelete:             n.OnAuthorDelete,
		IncludeHealthWarnings:      n.IncludeHealthWarnings,
		OnHealthIssue:              n.OnHealthIssue,
		OnApplicationUpdate:        n.OnApplicationUpdate,
		OnBookDelete:               n.OnBookDelete,
		OnBookFileDelete:           n.OnBookFileDelete,
		OnRename:                   n.OnRename,
		OnUpgrade:                  n.OnUpgrade,
		OnBookFileDeleteForUpgrade: n.OnBookFileDeleteForUpgrade,
		OnBookRetag:                n.OnBookRetag,
		Implementation:             types.StringValue(notificationWebhookImplementation),
		ConfigContract:             types.StringValue(notificationWebhookConfigContract),
	}
}

func (n *NotificationWebhook) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.URL = notification.URL
	n.Method = notification.Method
	n.Username = notification.Username
	n.Password = notification.Password
	n.Name = notification.Name
	n.ID = notification.ID
	n.OnGrab = notification.OnGrab
	n.OnBookFileDeleteForUpgrade = notification.OnBookFileDeleteForUpgrade
	n.OnBookFileDelete = notification.OnBookFileDelete
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnApplicationUpdate = notification.OnApplicationUpdate
	n.OnAuthorDelete = notification.OnAuthorDelete
	n.OnBookDelete = notification.OnBookDelete
	n.OnRename = notification.OnRename
	n.OnUpgrade = notification.OnUpgrade
	n.OnBookRetag = notification.OnBookRetag
	n.OnReleaseImport = notification.OnReleaseImport
}

func (r *NotificationWebhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationWebhookResourceName
}

func (r *NotificationWebhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Webhook resource.\nFor more information refer to [Notification](https://wiki.servarr.com/readarr/settings#connect) and [Webhook](https://wiki.servarr.com/readarr/supported#webhook).",
		Attributes: map[string]schema.Attribute{
			"on_grab": schema.BoolAttribute{
				MarkdownDescription: "On grab flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On upgrade flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_rename": schema.BoolAttribute{
				MarkdownDescription: "On rename flag.",
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
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_application_update": schema.BoolAttribute{
				MarkdownDescription: "On application update flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_book_retag": schema.BoolAttribute{
				MarkdownDescription: "On book retag flag.",
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
			"url": schema.StringAttribute{
				MarkdownDescription: "URL.",
				Required:            true,
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
			"method": schema.Int64Attribute{
				MarkdownDescription: "Method. `1` POST, `2` PUT.",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(1, 2),
				},
			},
		},
	}
}

func (r *NotificationWebhookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationWebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationWebhook

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationWebhook
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationWebhookResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationWebhookResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationWebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationWebhook

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationWebhook current value
	response, _, err := r.client.NotificationAPI.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationWebhookResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationWebhookResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationWebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationWebhook

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationWebhook
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationWebhookResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationWebhookResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationWebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationWebhook current value
	_, err := r.client.NotificationAPI.DeleteNotification(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationWebhookResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationWebhookResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationWebhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationWebhookResourceName+": "+req.ID)
}

func (n *NotificationWebhook) write(ctx context.Context, notification *readarr.NotificationResource, diags *diag.Diagnostics) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification, diags)
	n.fromNotification(genericNotification)
}

func (n *NotificationWebhook) read(ctx context.Context, diags *diag.Diagnostics) *readarr.NotificationResource {
	return n.toNotification().read(ctx, diags)
}
