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
	downloadClientUtorrentResourceName   = "download_client_utorrent"
	downloadClientUtorrentImplementation = "UTorrent"
	downloadClientUtorrentConfigContract = "UTorrentSettings"
	downloadClientUtorrentProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientUtorrentResource{}
	_ resource.ResourceWithImportState = &DownloadClientUtorrentResource{}
)

func NewDownloadClientUtorrentResource() resource.Resource {
	return &DownloadClientUtorrentResource{}
}

// DownloadClientUtorrentResource defines the download client implementation.
type DownloadClientUtorrentResource struct {
	client *readarr.APIClient
}

// DownloadClientUtorrent describes the download client data model.
type DownloadClientUtorrent struct {
	Tags                  types.Set    `tfsdk:"tags"`
	MusicImportedCategory types.String `tfsdk:"book_imported_category"`
	Name                  types.String `tfsdk:"name"`
	Host                  types.String `tfsdk:"host"`
	URLBase               types.String `tfsdk:"url_base"`
	Username              types.String `tfsdk:"username"`
	Password              types.String `tfsdk:"password"`
	MusicCategory         types.String `tfsdk:"book_category"`
	RecentTVPriority      types.Int64  `tfsdk:"recent_book_priority"`
	Priority              types.Int64  `tfsdk:"priority"`
	Port                  types.Int64  `tfsdk:"port"`
	ID                    types.Int64  `tfsdk:"id"`
	OlderTVPriority       types.Int64  `tfsdk:"older_book_priority"`
	IntialState           types.Int64  `tfsdk:"intial_state"`
	UseSsl                types.Bool   `tfsdk:"use_ssl"`
	Enable                types.Bool   `tfsdk:"enable"`
}

func (d DownloadClientUtorrent) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:                  d.Tags,
		Name:                  d.Name,
		Host:                  d.Host,
		URLBase:               d.URLBase,
		Username:              d.Username,
		Password:              d.Password,
		MusicCategory:         d.MusicCategory,
		RecentTVPriority:      d.RecentTVPriority,
		OlderTVPriority:       d.OlderTVPriority,
		Priority:              d.Priority,
		Port:                  d.Port,
		ID:                    d.ID,
		MusicImportedCategory: d.MusicImportedCategory,
		IntialState:           d.IntialState,
		UseSsl:                d.UseSsl,
		Enable:                d.Enable,
		Implementation:        types.StringValue(downloadClientUtorrentImplementation),
		ConfigContract:        types.StringValue(downloadClientUtorrentConfigContract),
		Protocol:              types.StringValue(downloadClientUtorrentProtocol),
	}
}

func (d *DownloadClientUtorrent) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.Name = client.Name
	d.Host = client.Host
	d.URLBase = client.URLBase
	d.Username = client.Username
	d.Password = client.Password
	d.MusicCategory = client.MusicCategory
	d.RecentTVPriority = client.RecentTVPriority
	d.OlderTVPriority = client.OlderTVPriority
	d.Priority = client.Priority
	d.Port = client.Port
	d.ID = client.ID
	d.MusicImportedCategory = client.MusicImportedCategory
	d.IntialState = client.IntialState
	d.UseSsl = client.UseSsl
	d.Enable = client.Enable
}

func (r *DownloadClientUtorrentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientUtorrentResourceName
}

func (r *DownloadClientUtorrentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Download Client uTorrent resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/readarr/settings#download-clients) and [uTorrent](https://wiki.servarr.com/readarr/supported#utorrent).",
		Attributes: map[string]schema.Attribute{
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Enable flag.",
				Optional:            true,
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Download Client name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Download Client ID.",
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
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
			},
			"recent_book_priority": schema.Int64Attribute{
				MarkdownDescription: "Recent Music priority. `0` Last, `1` First.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"older_book_priority": schema.Int64Attribute{
				MarkdownDescription: "Older Music priority. `0` Last, `1` First.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"intial_state": schema.Int64Attribute{
				MarkdownDescription: "Initial state, with Stop support. `0` Start, `1` ForceStart, `2` Pause, `3` Stop.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1, 2, 3),
				},
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "host.",
				Optional:            true,
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
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
			"book_category": schema.StringAttribute{
				MarkdownDescription: "Book category.",
				Optional:            true,
				Computed:            true,
			},
			"book_imported_category": schema.StringAttribute{
				MarkdownDescription: "Book imported category.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *DownloadClientUtorrentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *DownloadClientUtorrentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientUtorrent

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientUtorrent
	request := client.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.DownloadClientApi.CreateDownloadClient(ctx).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, downloadClientUtorrentResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientUtorrentResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientUtorrentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientUtorrent

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientUtorrent current value
	response, _, err := r.client.DownloadClientApi.GetDownloadClientById(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientUtorrentResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientUtorrentResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientUtorrentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientUtorrent

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientUtorrent
	request := client.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.DownloadClientApi.UpdateDownloadClient(ctx, strconv.Itoa(int(request.GetId()))).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, downloadClientUtorrentResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientUtorrentResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientUtorrentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientUtorrent current value
	_, err := r.client.DownloadClientApi.DeleteDownloadClient(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, downloadClientUtorrentResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientUtorrentResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientUtorrentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+downloadClientUtorrentResourceName+": "+req.ID)
}

func (d *DownloadClientUtorrent) write(ctx context.Context, downloadClient *readarr.DownloadClientResource, diags *diag.Diagnostics) {
	genericDownloadClient := d.toDownloadClient()
	genericDownloadClient.write(ctx, downloadClient, diags)
	d.fromDownloadClient(genericDownloadClient)
}

func (d *DownloadClientUtorrent) read(ctx context.Context, diags *diag.Diagnostics) *readarr.DownloadClientResource {
	return d.toDownloadClient().read(ctx, diags)
}
