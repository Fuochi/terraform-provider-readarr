package provider

import (
	"context"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const systemStatusDataSourceName = "system_status"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SystemStatusDataSource{}

func NewSystemStatusDataSource() datasource.DataSource {
	return &SystemStatusDataSource{}
}

// SystemStatusDataSource defines the system status implementation.
type SystemStatusDataSource struct {
	client *readarr.APIClient
}

// SystemStatus describes the system status data model.
type SystemStatus struct {
	PackageAuthor                 types.String `tfsdk:"package_author"`
	RuntimeName                   types.String `tfsdk:"runtime_name"`
	AppData                       types.String `tfsdk:"app_data"`
	DatabaseVersion               types.String `tfsdk:"database_version"`
	OsName                        types.String `tfsdk:"os_name"`
	AppName                       types.String `tfsdk:"app_name"`
	OsVersion                     types.String `tfsdk:"os_version"`
	StartTime                     types.String `tfsdk:"start_time"`
	BuildTime                     types.String `tfsdk:"build_time"`
	PackageUpdateMechanism        types.String `tfsdk:"package_update_mechanism"`
	PackageUpdateMechanismMessage types.String `tfsdk:"package_update_mechanism_message"`
	PackageVersion                types.String `tfsdk:"package_version"`
	RuntimeVersion                types.String `tfsdk:"runtime_version"`
	Version                       types.String `tfsdk:"version"`
	StartupPath                   types.String `tfsdk:"startup_path"`
	InstanceName                  types.String `tfsdk:"instance_name"`
	DatabaseType                  types.String `tfsdk:"database_type"`
	URLBase                       types.String `tfsdk:"url_base"`
	Mode                          types.String `tfsdk:"mode"`
	Branch                        types.String `tfsdk:"branch"`
	Authentication                types.String `tfsdk:"authentication"`
	MigrationVersion              types.Int64  `tfsdk:"migration_version"`
	ID                            types.Int64  `tfsdk:"id"`
	IsLinux                       types.Bool   `tfsdk:"is_linux"`
	IsProduction                  types.Bool   `tfsdk:"is_production"`
	IsDebug                       types.Bool   `tfsdk:"is_debug"`
	IsDocker                      types.Bool   `tfsdk:"is_docker"`
	IsWindows                     types.Bool   `tfsdk:"is_windows"`
	IsOsx                         types.Bool   `tfsdk:"is_osx"`
	IsNetCore                     types.Bool   `tfsdk:"is_net_core"`
	IsUserInteractive             types.Bool   `tfsdk:"is_user_interactive"`
	IsAdmin                       types.Bool   `tfsdk:"is_admin"`
}

func (d *SystemStatusDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + systemStatusDataSourceName
}

func (d *SystemStatusDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:System -->System Status resource. User must have rights to read `config.xml`.\nFor more information refer to [System Status](https://wiki.servarr.com/readarr/system#status) documentation.",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.Int64Attribute{
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
			},
			"is_debug": schema.BoolAttribute{
				MarkdownDescription: "Is debug flag.",
				Computed:            true,
			},
			"is_production": schema.BoolAttribute{
				MarkdownDescription: "Is production flag.",
				Computed:            true,
			},
			"is_admin": schema.BoolAttribute{
				MarkdownDescription: "Is admin flag.",
				Computed:            true,
			},
			"is_user_interactive": schema.BoolAttribute{
				MarkdownDescription: "Is user interactive flag.",
				Computed:            true,
			},
			"is_net_core": schema.BoolAttribute{
				MarkdownDescription: "Is net core flag.",
				Computed:            true,
			},
			"is_linux": schema.BoolAttribute{
				MarkdownDescription: "Is linux flag.",
				Computed:            true,
			},
			"is_osx": schema.BoolAttribute{
				MarkdownDescription: "Is osx flag.",
				Computed:            true,
			},
			"is_windows": schema.BoolAttribute{
				MarkdownDescription: "Is windows flag.",
				Computed:            true,
			},
			"is_docker": schema.BoolAttribute{
				MarkdownDescription: "Is docker flag.",
				Computed:            true,
			},
			"migration_version": schema.Int64Attribute{
				MarkdownDescription: "Migration version.",
				Computed:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Version.",
				Computed:            true,
			},
			"startup_path": schema.StringAttribute{
				MarkdownDescription: "Startup path.",
				Computed:            true,
			},
			"app_data": schema.StringAttribute{
				MarkdownDescription: "App data folder.",
				Computed:            true,
			},
			"os_name": schema.StringAttribute{
				MarkdownDescription: "OS name.",
				Computed:            true,
			},
			"os_version": schema.StringAttribute{
				MarkdownDescription: "OS version.",
				Computed:            true,
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "Mode.",
				Computed:            true,
			},
			"branch": schema.StringAttribute{
				MarkdownDescription: "Branch.",
				Computed:            true,
			},
			"authentication": schema.StringAttribute{
				MarkdownDescription: "Authentication.",
				Computed:            true,
			},
			"app_name": schema.StringAttribute{
				MarkdownDescription: "App name.",
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Computed:            true,
			},
			"runtime_version": schema.StringAttribute{
				MarkdownDescription: "Runtime version.",
				Computed:            true,
			},
			"runtime_name": schema.StringAttribute{
				MarkdownDescription: "Runtime name.",
				Computed:            true,
			},
			"package_version": schema.StringAttribute{
				MarkdownDescription: "Package version.",
				Computed:            true,
			},
			"package_author": schema.StringAttribute{
				MarkdownDescription: "Package author.",
				Computed:            true,
			},
			"package_update_mechanism": schema.StringAttribute{
				MarkdownDescription: "Package update mechanism.",
				Computed:            true,
			},
			"package_update_mechanism_message": schema.StringAttribute{
				MarkdownDescription: "Package update mechanism message.",
				Computed:            true,
			},
			"build_time": schema.StringAttribute{
				MarkdownDescription: "Build time.",
				Computed:            true,
			},
			"start_time": schema.StringAttribute{
				MarkdownDescription: "Start time.",
				Computed:            true,
			},
			"database_type": schema.StringAttribute{
				MarkdownDescription: "Database type.",
				Computed:            true,
			},
			"database_version": schema.StringAttribute{
				MarkdownDescription: "Database version.",
				Computed:            true,
			},
			"instance_name": schema.StringAttribute{
				MarkdownDescription: "Instance name.",
				Computed:            true,
			},
		},
	}
}

func (d *SystemStatusDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *SystemStatusDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get system status current value
	response, _, err := d.client.SystemAPI.GetSystemStatus(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, systemStatusDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+systemStatusDataSourceName)

	status := SystemStatus{}
	status.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, status)...)
}

func (s *SystemStatus) write(status *readarr.SystemResource) {
	s.IsDebug = types.BoolValue(status.GetIsDebug())
	s.IsProduction = types.BoolValue(status.GetIsProduction())
	s.IsAdmin = types.BoolValue(status.GetIsProduction())
	s.IsUserInteractive = types.BoolValue(status.GetIsUserInteractive())
	s.IsNetCore = types.BoolValue(status.GetIsNetCore())
	s.IsLinux = types.BoolValue(status.GetIsLinux())
	s.IsOsx = types.BoolValue(status.GetIsOsx())
	s.IsWindows = types.BoolValue(status.GetIsWindows())
	s.IsDocker = types.BoolValue(status.GetIsDocker())
	s.ID = types.Int64Value(int64(1))
	s.MigrationVersion = types.Int64Value(int64(status.GetMigrationVersion()))
	s.Version = types.StringValue(status.GetVersion())
	s.StartupPath = types.StringValue(status.GetStartupPath())
	s.AppData = types.StringValue(status.GetAppData())
	s.OsName = types.StringValue(status.GetOsName())
	s.OsVersion = types.StringValue(status.GetOsVersion())
	s.Mode = types.StringValue(string(status.GetMode()))
	s.Branch = types.StringValue(status.GetBranch())
	s.Authentication = types.StringValue(string(status.GetAuthentication()))
	s.URLBase = types.StringValue(status.GetUrlBase())
	s.RuntimeVersion = types.StringValue(status.GetRuntimeVersion())
	s.RuntimeName = types.StringValue(status.GetRuntimeName())
	s.PackageVersion = types.StringValue(status.GetPackageVersion())
	s.PackageAuthor = types.StringValue(status.GetPackageAuthor())
	s.PackageUpdateMechanism = types.StringValue(string(status.GetPackageUpdateMechanism()))
	s.PackageUpdateMechanismMessage = types.StringValue(status.GetPackageUpdateMechanismMessage())
	s.BuildTime = types.StringValue(status.GetBuildTime().String())
	s.StartTime = types.StringValue(status.GetStartTime().String())
	s.AppName = types.StringValue(status.GetAppName())
	s.DatabaseType = types.StringValue(string(status.GetDatabaseType()))
	s.DatabaseVersion = types.StringValue(status.GetDatabaseVersion())
	s.InstanceName = types.StringValue(status.GetInstanceName())
}
