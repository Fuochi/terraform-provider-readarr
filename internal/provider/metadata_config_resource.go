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

const metadataConfigResourceName = "metadata_config"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MetadataConfigResource{}
	_ resource.ResourceWithImportState = &MetadataConfigResource{}
)

func NewMetadataConfigResource() resource.Resource {
	return &MetadataConfigResource{}
}

// MetadataConfigResource defines the metadata config implementation.
type MetadataConfigResource struct {
	client *readarr.APIClient
}

// MetadataConfig describes the metadata config data model.
type MetadataConfig struct {
	WriteBookTags  types.String `tfsdk:"write_book_tags"`
	WriteAudioTags types.String `tfsdk:"write_audio_tags"`
	ID             types.Int64  `tfsdk:"id"`
	ScrubAudioTags types.Bool   `tfsdk:"scrub_audio_tags"`
	UpdateCovers   types.Bool   `tfsdk:"update_covers"`
	EmbedMetadata  types.Bool   `tfsdk:"embed_metadata"`
}

func (r *MetadataConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataConfigResourceName
}

func (r *MetadataConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Metadata -->Metadata Config resource.\nFor more information refer to [Metadata](https://wiki.servarr.com/readarr/settings#options) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Metadata Config ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"write_audio_tags": schema.StringAttribute{
				MarkdownDescription: "Write audio tags.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("no", "newFiles", "allFiles", "sync"),
				},
			},
			"write_book_tags": schema.StringAttribute{
				MarkdownDescription: "Write book tags.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("newFiles", "allFiles", "sync"),
				},
			},
			"scrub_audio_tags": schema.BoolAttribute{
				MarkdownDescription: "Scrub audio tags.",
				Required:            true,
			},
			"update_covers": schema.BoolAttribute{
				MarkdownDescription: "Update covers.",
				Required:            true,
			},
			"embed_metadata": schema.BoolAttribute{
				MarkdownDescription: "Embed metadata in book files.",
				Required:            true,
			},
		},
	}
}

func (r *MetadataConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *MetadataConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var config *MetadataConfig

	resp.Diagnostics.Append(req.Plan.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	request := config.read()
	request.SetId(1)

	// Create new MetadataConfig
	response, _, err := r.client.MetadataProviderConfigApi.UpdateMetadataProviderConfig(ctx, strconv.Itoa(int(request.GetId()))).MetadataProviderConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, metadataConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+metadataConfigResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (r *MetadataConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var config *MetadataConfig

	resp.Diagnostics.Append(req.State.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get metadataConfig current value
	response, _, err := r.client.MetadataProviderConfigApi.GetMetadataProviderConfig(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataConfigResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (r *MetadataConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var config *MetadataConfig

	resp.Diagnostics.Append(req.Plan.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	request := config.read()

	// Update MetadataConfig
	response, _, err := r.client.MetadataProviderConfigApi.UpdateMetadataProviderConfig(ctx, strconv.Itoa(int(request.GetId()))).MetadataProviderConfigResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, metadataConfigResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+metadataConfigResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (r *MetadataConfigResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	// MetadataConfig cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled "+metadataConfigResourceName+": 1")
	resp.State.RemoveResource(ctx)
}

func (r *MetadataConfigResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Trace(ctx, "imported "+metadataConfigResourceName+": 1")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), 1)...)
}

func (c *MetadataConfig) write(metadataConfig *readarr.MetadataProviderConfigResource) {
	c.ID = types.Int64Value(int64(metadataConfig.GetId()))
	c.WriteAudioTags = types.StringValue(string(metadataConfig.GetWriteAudioTags()))
	c.WriteBookTags = types.StringValue(string(metadataConfig.GetWriteBookTags()))
	c.ScrubAudioTags = types.BoolValue(metadataConfig.GetScrubAudioTags())
	c.UpdateCovers = types.BoolValue(metadataConfig.GetUpdateCovers())
	c.EmbedMetadata = types.BoolValue(metadataConfig.GetEmbedMetadata())
}

func (c *MetadataConfig) read() *readarr.MetadataProviderConfigResource {
	config := readarr.NewMetadataProviderConfigResource()
	config.SetWriteAudioTags(readarr.WriteAudioTagsType(c.WriteAudioTags.ValueString()))
	config.SetWriteBookTags(readarr.WriteBookTagsType(c.WriteBookTags.ValueString()))
	config.SetScrubAudioTags(c.ScrubAudioTags.ValueBool())
	config.SetUpdateCovers(c.UpdateCovers.ValueBool())
	config.SetEmbedMetadata(c.EmbedMetadata.ValueBool())
	config.SetId(int32(c.ID.ValueInt64()))

	return config
}
