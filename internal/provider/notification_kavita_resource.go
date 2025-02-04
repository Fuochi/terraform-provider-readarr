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
	notificationKavitaResourceName   = "notification_kavita"
	notificationKavitaImplementation = "Kavita"
	notificationKavitaConfigContract = "KavitaSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationKavitaResource{}
	_ resource.ResourceWithImportState = &NotificationKavitaResource{}
)

func NewNotificationKavitaResource() resource.Resource {
	return &NotificationKavitaResource{}
}

// NotificationKavitaResource defines the notification implementation.
type NotificationKavitaResource struct {
	client *readarr.APIClient
}

// NotificationKavita describes the notification data model.
type NotificationKavita struct {
	Tags                       types.Set    `tfsdk:"tags"`
	APIKey                     types.String `tfsdk:"api_key"`
	Host                       types.String `tfsdk:"host"`
	Name                       types.String `tfsdk:"name"`
	Port                       types.Int64  `tfsdk:"port"`
	ID                         types.Int64  `tfsdk:"id"`
	UseSSL                     types.Bool   `tfsdk:"use_ssl"`
	Notify                     types.Bool   `tfsdk:"notify"`
	OnUpgrade                  types.Bool   `tfsdk:"on_upgrade"`
	OnReleaseImport            types.Bool   `tfsdk:"on_release_import"`
	OnBookDelete               types.Bool   `tfsdk:"on_book_delete"`
	OnBookFileDelete           types.Bool   `tfsdk:"on_book_file_delete"`
	OnBookFileDeleteForUpgrade types.Bool   `tfsdk:"on_book_file_delete_for_upgrade"`
	OnBookRetag                types.Bool   `tfsdk:"on_book_retag"`
}

func (n NotificationKavita) toNotification() *Notification {
	return &Notification{
		Tags:                       n.Tags,
		APIKey:                     n.APIKey,
		Host:                       n.Host,
		Name:                       n.Name,
		Port:                       n.Port,
		ID:                         n.ID,
		UseSSL:                     n.UseSSL,
		OnReleaseImport:            n.OnReleaseImport,
		Notify:                     n.Notify,
		OnBookDelete:               n.OnBookDelete,
		OnBookFileDelete:           n.OnBookFileDelete,
		OnUpgrade:                  n.OnUpgrade,
		OnBookFileDeleteForUpgrade: n.OnBookFileDeleteForUpgrade,
		OnBookRetag:                n.OnBookRetag,
		Implementation:             types.StringValue(notificationKavitaImplementation),
		ConfigContract:             types.StringValue(notificationKavitaConfigContract),
	}
}

func (n *NotificationKavita) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.APIKey = notification.APIKey
	n.Host = notification.Host
	n.Name = notification.Name
	n.Port = notification.Port
	n.ID = notification.ID
	n.UseSSL = notification.UseSSL
	n.OnBookFileDeleteForUpgrade = notification.OnBookFileDeleteForUpgrade
	n.OnBookFileDelete = notification.OnBookFileDelete
	n.Notify = notification.Notify
	n.OnBookDelete = notification.OnBookDelete
	n.OnUpgrade = notification.OnUpgrade
	n.OnReleaseImport = notification.OnReleaseImport
	n.OnBookRetag = notification.OnBookRetag
}

func (r *NotificationKavitaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationKavitaResourceName
}

func (r *NotificationKavitaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Kavita resource.\nFor more information refer to [Notification](https://wiki.servarr.com/readarr/settings#connect) and [Kavita](https://wiki.servarr.com/readarr/supported#kavita).",
		Attributes: map[string]schema.Attribute{
			"on_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On upgrade flag.",
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
			"use_ssl": schema.BoolAttribute{
				MarkdownDescription: "Use SSL flag.",
				Optional:            true,
				Computed:            true,
			},
			"notify": schema.BoolAttribute{
				MarkdownDescription: "Update library.",
				Required:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Required:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Required:            true,
				Sensitive:           true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Host.",
				Required:            true,
			},
		},
	}
}

func (r *NotificationKavitaResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationKavitaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationKavita

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationKavita
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationKavitaResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationKavitaResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationKavitaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationKavita

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationKavita current value
	response, _, err := r.client.NotificationAPI.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationKavitaResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationKavitaResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationKavitaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationKavita

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationKavita
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationKavitaResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationKavitaResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationKavitaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationKavita current value
	_, err := r.client.NotificationAPI.DeleteNotification(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationKavitaResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationKavitaResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationKavitaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationKavitaResourceName+": "+req.ID)
}

func (n *NotificationKavita) write(ctx context.Context, notification *readarr.NotificationResource, diags *diag.Diagnostics) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification, diags)
	n.fromNotification(genericNotification)
}

func (n *NotificationKavita) read(ctx context.Context, diags *diag.Diagnostics) *readarr.NotificationResource {
	return n.toNotification().read(ctx, diags)
}
