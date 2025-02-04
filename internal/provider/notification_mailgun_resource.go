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
	notificationMailgunResourceName   = "notification_mailgun"
	notificationMailgunImplementation = "Mailgun"
	notificationMailgunConfigContract = "MailgunSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationMailgunResource{}
	_ resource.ResourceWithImportState = &NotificationMailgunResource{}
)

func NewNotificationMailgunResource() resource.Resource {
	return &NotificationMailgunResource{}
}

// NotificationMailgunResource defines the notification implementation.
type NotificationMailgunResource struct {
	client *readarr.APIClient
}

// NotificationMailgun describes the notification data model.
type NotificationMailgun struct {
	Tags                       types.Set    `tfsdk:"tags"`
	Recipients                 types.Set    `tfsdk:"recipients"`
	From                       types.String `tfsdk:"from"`
	SenderDomain               types.String `tfsdk:"sender_domain"`
	APIKey                     types.String `tfsdk:"api_key"`
	Name                       types.String `tfsdk:"name"`
	ID                         types.Int64  `tfsdk:"id"`
	UseEuEndpoint              types.Bool   `tfsdk:"use_eu_endpoint"`
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

func (n NotificationMailgun) toNotification() *Notification {
	return &Notification{
		Tags:                       n.Tags,
		From:                       n.From,
		Recipients:                 n.Recipients,
		SenderDomain:               n.SenderDomain,
		APIKey:                     n.APIKey,
		Name:                       n.Name,
		ID:                         n.ID,
		UseEuEndpoint:              n.UseEuEndpoint,
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
		Implementation:             types.StringValue(notificationMailgunImplementation),
		ConfigContract:             types.StringValue(notificationMailgunConfigContract),
	}
}

func (n *NotificationMailgun) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.From = notification.From
	n.Recipients = notification.Recipients
	n.SenderDomain = notification.SenderDomain
	n.APIKey = notification.APIKey
	n.Name = notification.Name
	n.ID = notification.ID
	n.UseEuEndpoint = notification.UseEuEndpoint
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

func (r *NotificationMailgunResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationMailgunResourceName
}

func (r *NotificationMailgunResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Mailgun resource.\nFor more information refer to [Notification](https://wiki.servarr.com/readarr/settings#connect) and [Mailgun](https://wiki.servarr.com/readarr/supported#mailgun).",
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
			"use_eu_endpoint": schema.BoolAttribute{
				MarkdownDescription: "Use EU endpoint flag.",
				Optional:            true,
				Computed:            true,
			},
			"sender_domain": schema.StringAttribute{
				MarkdownDescription: "Sender domain.",
				Optional:            true,
				Computed:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API Key.",
				Required:            true,
				Sensitive:           true,
			},
			"from": schema.StringAttribute{
				MarkdownDescription: "From.",
				Required:            true,
			},
			"recipients": schema.SetAttribute{
				MarkdownDescription: "Recipients.",
				Required:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *NotificationMailgunResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationMailgunResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationMailgun

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationMailgun
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationMailgunResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationMailgunResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationMailgunResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationMailgun

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationMailgun current value
	response, _, err := r.client.NotificationAPI.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationMailgunResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationMailgunResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationMailgunResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationMailgun

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationMailgun
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationMailgunResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationMailgunResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationMailgunResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationMailgun current value
	_, err := r.client.NotificationAPI.DeleteNotification(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationMailgunResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationMailgunResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationMailgunResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationMailgunResourceName+": "+req.ID)
}

func (n *NotificationMailgun) write(ctx context.Context, notification *readarr.NotificationResource, diags *diag.Diagnostics) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification, diags)
	n.fromNotification(genericNotification)
}

func (n *NotificationMailgun) read(ctx context.Context, diags *diag.Diagnostics) *readarr.NotificationResource {
	return n.toNotification().read(ctx, diags)
}
