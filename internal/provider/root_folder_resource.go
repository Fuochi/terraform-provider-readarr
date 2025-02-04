package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const rootFolderResourceName = "root_folder"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &RootFolderResource{}
	_ resource.ResourceWithImportState = &RootFolderResource{}
)

func NewRootFolderResource() resource.Resource {
	return &RootFolderResource{}
}

// RootFolderResource defines the root folder implementation.
type RootFolderResource struct {
	client *readarr.APIClient
}

// RootFolder describes the root folder data model.
type RootFolder struct {
	DefaultTags                 types.Set    `tfsdk:"default_tags"`
	Path                        types.String `tfsdk:"path"`
	Name                        types.String `tfsdk:"name"`
	DefaultMonitorOption        types.String `tfsdk:"default_monitor_option"`
	DefaultNewItemMonitorOption types.String `tfsdk:"default_monitor_new_item_option"`
	Host                        types.String `tfsdk:"host"`
	Username                    types.String `tfsdk:"username"`
	Password                    types.String `tfsdk:"password"`
	Library                     types.String `tfsdk:"library"`
	OutputProfile               types.String `tfsdk:"output_profile"`
	Port                        types.Int64  `tfsdk:"port"`
	DefaultMetadataProfileID    types.Int64  `tfsdk:"default_metadata_profile_id"`
	DefaultQualityProfileID     types.Int64  `tfsdk:"default_quality_profile_id"`
	ID                          types.Int64  `tfsdk:"id"`
	Accessible                  types.Bool   `tfsdk:"accessible"`
	IsCalibreLibrary            types.Bool   `tfsdk:"is_calibre_library"`
	UseSSL                      types.Bool   `tfsdk:"use_ssl"`
}

func (r RootFolder) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"default_tags":                    types.SetType{}.WithElementType(types.Int64Type),
			"path":                            types.StringType,
			"name":                            types.StringType,
			"default_monitor_option":          types.StringType,
			"default_monitor_new_item_option": types.StringType,
			"host":                            types.StringType,
			"username":                        types.StringType,
			"password":                        types.StringType,
			"library":                         types.StringType,
			"output_profile":                  types.StringType,
			"id":                              types.Int64Type,
			"port":                            types.Int64Type,
			"default_metadata_profile_id":     types.Int64Type,
			"default_quality_profile_id":      types.Int64Type,
			"accessible":                      types.BoolType,
			"is_calibre_library":              types.BoolType,
			"use_ssl":                         types.BoolType,
		})
}

func (r *RootFolderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + rootFolderResourceName
}

func (r *RootFolderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Media Management -->Root Folder resource.\nFor more information refer to [Root Folders](https://wiki.servarr.com/readarr/settings#root-folders) documentation.",
		Attributes: map[string]schema.Attribute{
			// TODO: add validator
			"path": schema.StringAttribute{
				MarkdownDescription: "Root Folder absolute path.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Root Folder name.",
				Required:            true,
			},
			"default_monitor_option": schema.StringAttribute{
				MarkdownDescription: "Default monitor option.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("all", "future", "missing", "existing", "latest", "first", "none", "unknown"),
				},
			},
			"default_monitor_new_item_option": schema.StringAttribute{
				MarkdownDescription: "Default monitor new item option.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("all", "future", "missing", "existing", "latest", "first", "none", "unknown"),
				},
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Calibre host.",
				Optional:            true,
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Calibre username.",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Calibre password.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"library": schema.StringAttribute{
				MarkdownDescription: "Calibre library.",
				Optional:            true,
				Computed:            true,
			},
			"output_profile": schema.StringAttribute{
				MarkdownDescription: "Calibre output profile.",
				Optional:            true,
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Calibre Port.",
				Optional:            true,
				Computed:            true,
			},
			"default_metadata_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Default metadata profile ID.",
				Required:            true,
			},
			"default_quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Default metadata profile ID.",
				Required:            true,
			},
			"accessible": schema.BoolAttribute{
				MarkdownDescription: "Access flag.",
				Computed:            true,
			},
			"is_calibre_library": schema.BoolAttribute{
				MarkdownDescription: "Is calibre library flag.",
				Required:            true,
			},
			"use_ssl": schema.BoolAttribute{
				MarkdownDescription: "Use SSL for calibre connection.",
				Optional:            true,
				Computed:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Root Folder ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"default_tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *RootFolderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *RootFolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var folder *RootFolder

	resp.Diagnostics.Append(req.Plan.Get(ctx, &folder)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new RootFolder
	request := folder.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.RootFolderAPI.CreateRootFolder(ctx).RootFolderResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, rootFolderResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+rootFolderResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	folder.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &folder)...)
}

func (r *RootFolderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var folder *RootFolder

	resp.Diagnostics.Append(req.State.Get(ctx, &folder)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get rootFolder current value
	response, _, err := r.client.RootFolderAPI.GetRootFolderById(ctx, int32(folder.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, rootFolderResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+rootFolderResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	folder.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &folder)...)
}

func (r *RootFolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var folder *RootFolder

	resp.Diagnostics.Append(req.Plan.Get(ctx, &folder)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update RootFolder
	request := folder.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.RootFolderAPI.UpdateRootFolder(ctx, strconv.Itoa(int(request.GetId()))).RootFolderResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, rootFolderResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+rootFolderResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	folder.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &folder)...)
}

func (r *RootFolderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete rootFolder current value
	_, err := r.client.RootFolderAPI.DeleteRootFolder(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, rootFolderResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+rootFolderResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *RootFolderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+rootFolderResourceName+": "+req.ID)
}

func (r *RootFolder) write(ctx context.Context, rootFolder *readarr.RootFolderResource, diags *diag.Diagnostics) {
	var tempDiag diag.Diagnostics

	r.Accessible = types.BoolValue(rootFolder.GetAccessible())
	r.Name = types.StringValue(rootFolder.GetName())
	r.ID = types.Int64Value(int64(rootFolder.GetId()))
	r.Path = types.StringValue(rootFolder.GetPath())
	r.DefaultMonitorOption = types.StringValue(string(rootFolder.GetDefaultMonitorOption()))
	r.DefaultNewItemMonitorOption = types.StringValue(string(rootFolder.GetDefaultNewItemMonitorOption()))
	r.Host = types.StringValue(rootFolder.GetHost())
	r.Username = types.StringValue(rootFolder.GetUsername())
	r.Password = types.StringValue(rootFolder.GetPassword())
	r.Library = types.StringValue(rootFolder.GetLibrary())
	r.OutputProfile = types.StringValue(rootFolder.GetOutputProfile())
	r.DefaultMetadataProfileID = types.Int64Value(int64(rootFolder.GetDefaultMetadataProfileId()))
	r.DefaultQualityProfileID = types.Int64Value(int64(rootFolder.GetDefaultQualityProfileId()))
	r.Port = types.Int64Value(int64(rootFolder.GetPort()))
	r.IsCalibreLibrary = types.BoolValue(rootFolder.GetIsCalibreLibrary())
	r.UseSSL = types.BoolValue(rootFolder.GetUseSsl())
	r.DefaultTags, tempDiag = types.SetValueFrom(ctx, types.Int64Type, rootFolder.GetDefaultTags())
	diags.Append(tempDiag...)
}

func (r *RootFolder) read(ctx context.Context, diags *diag.Diagnostics) *readarr.RootFolderResource {
	folder := readarr.NewRootFolderResource()
	folder.SetId(int32(r.ID.ValueInt64()))
	folder.SetName(r.Name.ValueString())
	folder.SetPath(r.Path.ValueString())
	folder.SetDefaultMonitorOption(readarr.MonitorTypes(r.DefaultMonitorOption.ValueString()))
	folder.SetDefaultNewItemMonitorOption(readarr.NewItemMonitorTypes(r.DefaultNewItemMonitorOption.ValueString()))
	folder.SetHost(r.Host.ValueString())
	folder.SetUsername(r.Username.ValueString())
	folder.SetPassword(r.Password.ValueString())
	folder.SetLibrary(r.Library.ValueString())
	folder.SetOutputProfile(r.OutputProfile.ValueString())
	folder.SetDefaultMetadataProfileId(int32(r.DefaultMetadataProfileID.ValueInt64()))
	folder.SetDefaultQualityProfileId(int32(r.DefaultQualityProfileID.ValueInt64()))
	folder.SetPort(int32(r.Port.ValueInt64()))
	folder.SetIsCalibreLibrary(r.IsCalibreLibrary.ValueBool())
	folder.SetUseSsl(r.UseSSL.ValueBool())
	diags.Append(r.DefaultTags.ElementsAs(ctx, &folder.DefaultTags, true)...)

	return folder
}
