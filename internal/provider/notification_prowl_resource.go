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
	notificationProwlResourceName   = "notification_prowl"
	notificationProwlImplementation = "Prowl"
	notificationProwlConfigContract = "ProwlSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationProwlResource{}
	_ resource.ResourceWithImportState = &NotificationProwlResource{}
)

func NewNotificationProwlResource() resource.Resource {
	return &NotificationProwlResource{}
}

// NotificationProwlResource defines the notification implementation.
type NotificationProwlResource struct {
	client *readarr.APIClient
}

// NotificationProwl describes the notification data model.
type NotificationProwl struct {
	Tags                       types.Set    `tfsdk:"tags"`
	APIKey                     types.String `tfsdk:"api_key"`
	Name                       types.String `tfsdk:"name"`
	Priority                   types.Int64  `tfsdk:"priority"`
	ID                         types.Int64  `tfsdk:"id"`
	OnGrab                     types.Bool   `tfsdk:"on_grab"`
	IncludeHealthWarnings      types.Bool   `tfsdk:"include_health_warnings"`
	OnHealthIssue              types.Bool   `tfsdk:"on_health_issue"`
	OnApplicationUpdate        types.Bool   `tfsdk:"on_application_update"`
	OnUpgrade                  types.Bool   `tfsdk:"on_upgrade"`
	OnReleaseImport            types.Bool   `tfsdk:"on_release_import"`
	OnAuthorDelete             types.Bool   `tfsdk:"on_author_delete"`
	OnBookDelete               types.Bool   `tfsdk:"on_book_delete"`
	OnBookFileDelete           types.Bool   `tfsdk:"on_book_file_delete"`
	OnBookFileDeleteForUpgrade types.Bool   `tfsdk:"on_book_file_delete_for_upgrade"`
}

func (n NotificationProwl) toNotification() *Notification {
	return &Notification{
		Tags:                       n.Tags,
		APIKey:                     n.APIKey,
		Name:                       n.Name,
		Priority:                   n.Priority,
		ID:                         n.ID,
		OnGrab:                     n.OnGrab,
		OnReleaseImport:            n.OnReleaseImport,
		OnAuthorDelete:             n.OnAuthorDelete,
		IncludeHealthWarnings:      n.IncludeHealthWarnings,
		OnHealthIssue:              n.OnHealthIssue,
		OnApplicationUpdate:        n.OnApplicationUpdate,
		OnBookDelete:               n.OnBookDelete,
		OnBookFileDelete:           n.OnBookFileDelete,
		OnUpgrade:                  n.OnUpgrade,
		OnBookFileDeleteForUpgrade: n.OnBookFileDeleteForUpgrade,
		Implementation:             types.StringValue(notificationProwlImplementation),
		ConfigContract:             types.StringValue(notificationProwlConfigContract),
	}
}

func (n *NotificationProwl) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.APIKey = notification.APIKey
	n.Name = notification.Name
	n.Priority = notification.Priority
	n.ID = notification.ID
	n.OnGrab = notification.OnGrab
	n.OnBookFileDeleteForUpgrade = notification.OnBookFileDeleteForUpgrade
	n.OnBookFileDelete = notification.OnBookFileDelete
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnApplicationUpdate = notification.OnApplicationUpdate
	n.OnAuthorDelete = notification.OnAuthorDelete
	n.OnBookDelete = notification.OnBookDelete
	n.OnUpgrade = notification.OnUpgrade
	n.OnReleaseImport = notification.OnReleaseImport
}

func (r *NotificationProwlResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationProwlResourceName
}

func (r *NotificationProwlResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Prowl resource.\nFor more information refer to [Notification](https://wiki.servarr.com/readarr/settings#connect) and [Prowl](https://wiki.servarr.com/readarr/supported#prowl).",
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
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.`-2` Very Low, `-1` Low, `0` Normal, `1` High, `2` Emergency.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(-2, -1, 0, 1, 2),
				},
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API Key.",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *NotificationProwlResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationProwlResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationProwl

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationProwl
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationProwlResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationProwlResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationProwlResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationProwl

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationProwl current value
	response, _, err := r.client.NotificationAPI.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationProwlResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationProwlResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationProwlResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationProwl

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationProwl
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationProwlResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationProwlResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationProwlResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationProwl current value
	_, err := r.client.NotificationAPI.DeleteNotification(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationProwlResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationProwlResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationProwlResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationProwlResourceName+": "+req.ID)
}

func (n *NotificationProwl) write(ctx context.Context, notification *readarr.NotificationResource, diags *diag.Diagnostics) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification, diags)
	n.fromNotification(genericNotification)
}

func (n *NotificationProwl) read(ctx context.Context, diags *diag.Diagnostics) *readarr.NotificationResource {
	return n.toNotification().read(ctx, diags)
}
