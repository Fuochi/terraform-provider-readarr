package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const authorResourceName = "author"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &AuthorResource{}
	_ resource.ResourceWithImportState = &AuthorResource{}
)

func NewAuthorResource() resource.Resource {
	return &AuthorResource{}
}

// AuthorResource defines the author implementation.
type AuthorResource struct {
	client *readarr.APIClient
}

// Author describes the author data model.
type Author struct {
	Genres           types.Set    `tfsdk:"genres"`
	Tags             types.Set    `tfsdk:"tags"`
	AuthorName       types.String `tfsdk:"author_name"`
	ForeignAuthorID  types.String `tfsdk:"foreign_author_id"`
	Status           types.String `tfsdk:"status"`
	Path             types.String `tfsdk:"path"`
	Overview         types.String `tfsdk:"overview"`
	ID               types.Int64  `tfsdk:"id"`
	QualityProfileID types.Int64  `tfsdk:"quality_profile_id"`
	Monitored        types.Bool   `tfsdk:"monitored"`

	// TODO: future Implementation
	// Links          types.Set    `tfsdk:"links"`
	// SortName       types.String `tfsdk:"sortName"`
	// Ended          types.Bool   `tfsdk:"ended"`
	// RootFolderPath types.String `tfsdk:"root_folder_path"`
	// FolderName     types.String `tfsdk:"folderName"`
	// CleanName      types.String `tfsdk:"cleanName"`
	// Added          types.String `tfsdk:"added"`
	// Ratings        types.Object `tfsdk:"ratings"`
}

func (a Author) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(
		map[string]attr.Type{
			"genres":             types.SetType{}.WithElementType(types.StringType),
			"tags":               types.SetType{}.WithElementType(types.Int64Type),
			"author_name":        types.StringType,
			"foreign_author_id":  types.StringType,
			"status":             types.StringType,
			"path":               types.StringType,
			"overview":           types.StringType,
			"id":                 types.Int64Type,
			"quality_profile_id": types.Int64Type,
			"monitored":          types.BoolType,
		})
}

func (r *AuthorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + authorResourceName
}

func (r *AuthorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Authors -->Author resource.\nFor more information refer to [Authors](https://wiki.servarr.com/readarr/library#authors) documentation.",
		Attributes: map[string]schema.Attribute{
			"monitored": schema.BoolAttribute{
				MarkdownDescription: "Monitored flag.",
				Required:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality profile ID.",
				Required:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Author ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"author_name": schema.StringAttribute{
				MarkdownDescription: "Author name.",
				Required:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Full author path.",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Author status.",
				Computed:            true,
			},
			"overview": schema.StringAttribute{
				MarkdownDescription: "Overview.",
				Computed:            true,
			},
			"foreign_author_id": schema.StringAttribute{
				MarkdownDescription: "Foreign author ID.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"genres": schema.SetAttribute{
				MarkdownDescription: "List genres.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *AuthorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *AuthorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var author *Author

	resp.Diagnostics.Append(req.Plan.Get(ctx, &author)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Author
	request := author.read(ctx, &resp.Diagnostics)
	// TODO: can parametrize AddAuthorOptions
	options := readarr.NewAddAuthorOptions()
	options.SetMonitor(readarr.MONITORTYPES_ALL)
	options.SetSearchForMissingBooks(true)

	response, _, err := r.client.AuthorAPI.CreateAuthor(ctx).AuthorResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, authorResourceName, err))

		return
	}

	tflog.Trace(ctx, "created author: "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	author.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &author)...)
}

func (r *AuthorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var author *Author

	resp.Diagnostics.Append(req.State.Get(ctx, &author)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get author current value
	response, _, err := r.client.AuthorAPI.GetAuthorById(ctx, int32(author.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, authorResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+authorResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	author.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &author)...)
}

func (r *AuthorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var author *Author

	resp.Diagnostics.Append(req.Plan.Get(ctx, &author)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Author
	request := author.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.AuthorAPI.UpdateAuthor(ctx, fmt.Sprint(request.GetId())).AuthorResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, authorResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+authorResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	author.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &author)...)
}

func (r *AuthorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete author current value
	_, err := r.client.AuthorAPI.DeleteAuthor(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, authorResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+authorResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *AuthorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+authorResourceName+": "+req.ID)
}

func (a *Author) write(ctx context.Context, author *readarr.AuthorResource, diags *diag.Diagnostics) {
	var tempDiag diag.Diagnostics

	a.Monitored = types.BoolValue(author.GetMonitored())
	a.ID = types.Int64Value(int64(author.GetId()))
	a.AuthorName = types.StringValue(author.GetAuthorName())
	a.Path = types.StringValue(author.GetPath())
	a.QualityProfileID = types.Int64Value(int64(author.GetQualityProfileId()))
	a.ForeignAuthorID = types.StringValue(author.GetForeignAuthorId())
	a.Status = types.StringValue(string(author.GetStatus()))
	a.Overview = types.StringValue(author.GetOverview())
	a.Genres = types.SetValueMust(types.StringType, nil)
	a.Tags, tempDiag = types.SetValueFrom(ctx, types.Int64Type, author.GetTags())
	diags.Append(tempDiag...)
	a.Genres, tempDiag = types.SetValueFrom(ctx, types.StringType, author.GetGenres())
	diags.Append(tempDiag...)
}

func (a *Author) read(ctx context.Context, diags *diag.Diagnostics) *readarr.AuthorResource {
	author := readarr.NewAuthorResource()
	author.SetMonitored(a.Monitored.ValueBool())
	author.SetAuthorName(a.AuthorName.ValueString())
	author.SetPath(a.Path.ValueString())
	author.SetQualityProfileId(int32(a.QualityProfileID.ValueInt64()))
	author.SetForeignAuthorId(a.ForeignAuthorID.ValueString())
	author.SetId(int32(a.ID.ValueInt64()))
	diags.Append(a.Tags.ElementsAs(ctx, &author.Tags, true)...)
	// Fix unused but required profile
	author.SetMetadataProfileId(1)

	return author
}
