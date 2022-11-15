package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/readarr"
)

const systemStatusDataSourceName = "system_status"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SystemStatusDataSource{}

func NewSystemStatusDataSource() datasource.DataSource {
	return &SystemStatusDataSource{}
}

// SystemStatusDataSource defines the system status implementation.
type SystemStatusDataSource struct {
	client *readarr.Readarr
}

// SystemStatus describes the system status data model.
type SystemStatus struct {
	IsDebug           types.Bool `tfsdk:"is_debug"`
	IsProduction      types.Bool `tfsdk:"is_production"`
	IsAdmin           types.Bool `tfsdk:"is_admin"`
	IsUserInteractive types.Bool `tfsdk:"is_user_interactive"`
	IsNetCore         types.Bool `tfsdk:"is_net_core"`
	IsMono            types.Bool `tfsdk:"is_mono"`
	IsLinux           types.Bool `tfsdk:"is_linux"`
	IsOsx             types.Bool `tfsdk:"is_osx"`
	IsWindows         types.Bool `tfsdk:"is_windows"`
	IsDocker          types.Bool `tfsdk:"is_docker"`
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	ID                     types.Int64  `tfsdk:"id"`
	MigrationVersion       types.Int64  `tfsdk:"migration_version"`
	Version                types.String `tfsdk:"version"`
	StartupPath            types.String `tfsdk:"startup_path"`
	AppData                types.String `tfsdk:"app_data"`
	OsName                 types.String `tfsdk:"os_name"`
	OsVersion              types.String `tfsdk:"os_version"`
	Mode                   types.String `tfsdk:"mode"`
	Branch                 types.String `tfsdk:"branch"`
	Authentication         types.String `tfsdk:"authentication"`
	URLBase                types.String `tfsdk:"url_base"`
	RuntimeVersion         types.String `tfsdk:"runtime_version"`
	RuntimeName            types.String `tfsdk:"runtime_name"`
	PackageVersion         types.String `tfsdk:"package_version"`
	PackageAuthor          types.String `tfsdk:"package_author"`
	PackageUpdateMechanism types.String `tfsdk:"package_update_mechanism"`
	BuildTime              types.String `tfsdk:"build_time"`
	StartTime              types.String `tfsdk:"start_time"`
	AppName                types.String `tfsdk:"app_name"`
	DatabaseType           types.String `tfsdk:"database_type"`
	DatabaseVersion        types.String `tfsdk:"database_version"`
	InstanceName           types.String `tfsdk:"instance_name"`
}

func (d *SystemStatusDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + systemStatusDataSourceName
}

func (d *SystemStatusDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Status -->System Status resource. User must have rights to read `config.xml`.\nFor more information refer to [System Status](https://wiki.servarr.com/readarr/system#status) documentation.",
		Attributes: map[string]tfsdk.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": {
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"is_debug": {
				MarkdownDescription: "Is debug flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_production": {
				MarkdownDescription: "Is production flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_admin": {
				MarkdownDescription: "Is admin flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_user_interactive": {
				MarkdownDescription: "Is user interactive flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_net_core": {
				MarkdownDescription: "Is net core flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_mono": {
				MarkdownDescription: "Is mono flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_linux": {
				MarkdownDescription: "Is linux flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_osx": {
				MarkdownDescription: "Is osx flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_windows": {
				MarkdownDescription: "Is windows flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"is_docker": {
				MarkdownDescription: "Is docker flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"migration_version": {
				MarkdownDescription: "Migration version.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"version": {
				MarkdownDescription: "Version.",
				Computed:            true,
				Type:                types.StringType,
			},
			"startup_path": {
				MarkdownDescription: "Startup path.",
				Computed:            true,
				Type:                types.StringType,
			},
			"app_data": {
				MarkdownDescription: "App data folder.",
				Computed:            true,
				Type:                types.StringType,
			},
			"os_name": {
				MarkdownDescription: "OS name.",
				Computed:            true,
				Type:                types.StringType,
			},
			"os_version": {
				MarkdownDescription: "OS version.",
				Computed:            true,
				Type:                types.StringType,
			},
			"mode": {
				MarkdownDescription: "Mode.",
				Computed:            true,
				Type:                types.StringType,
			},
			"branch": {
				MarkdownDescription: "Branch.",
				Computed:            true,
				Type:                types.StringType,
			},
			"authentication": {
				MarkdownDescription: "Authentication.",
				Computed:            true,
				Type:                types.StringType,
			},
			"app_name": {
				MarkdownDescription: "App name.",
				Computed:            true,
				Type:                types.StringType,
			},
			"url_base": {
				MarkdownDescription: "Base URL.",
				Computed:            true,
				Type:                types.StringType,
			},
			"runtime_version": {
				MarkdownDescription: "Runtime version.",
				Computed:            true,
				Type:                types.StringType,
			},
			"runtime_name": {
				MarkdownDescription: "Runtime name.",
				Computed:            true,
				Type:                types.StringType,
			},
			"package_version": {
				MarkdownDescription: "Package version.",
				Computed:            true,
				Type:                types.StringType,
			},
			"package_author": {
				MarkdownDescription: "Package author.",
				Computed:            true,
				Type:                types.StringType,
			},
			"package_update_mechanism": {
				MarkdownDescription: "Package update mechanism.",
				Computed:            true,
				Type:                types.StringType,
			},
			"build_time": {
				MarkdownDescription: "Build time.",
				Computed:            true,
				Type:                types.StringType,
			},
			"start_time": {
				MarkdownDescription: "Start time.",
				Computed:            true,
				Type:                types.StringType,
			},
			"database_type": {
				MarkdownDescription: "Database type.",
				Computed:            true,
				Type:                types.StringType,
			},
			"database_version": {
				MarkdownDescription: "Database version.",
				Computed:            true,
				Type:                types.StringType,
			},
			"instance_name": {
				MarkdownDescription: "Instance name.",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (d *SystemStatusDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*readarr.Readarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *readarr.Readarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *SystemStatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get naming current value
	response, err := d.client.GetSystemStatusContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", systemStatusDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+systemStatusDataSourceName)

	result := writeSystemStatus(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func writeSystemStatus(status *readarr.SystemStatus) *SystemStatus {
	return &SystemStatus{
		IsDebug:                types.BoolValue(status.IsDebug),
		IsProduction:           types.BoolValue(status.IsProduction),
		IsAdmin:                types.BoolValue(status.IsProduction),
		IsUserInteractive:      types.BoolValue(status.IsUserInteractive),
		IsNetCore:              types.BoolValue(status.IsNetCore),
		IsMono:                 types.BoolValue(status.IsMono),
		IsLinux:                types.BoolValue(status.IsLinux),
		IsOsx:                  types.BoolValue(status.IsOsx),
		IsWindows:              types.BoolValue(status.IsWindows),
		IsDocker:               types.BoolValue(status.IsDocker),
		ID:                     types.Int64Value(int64(1)),
		MigrationVersion:       types.Int64Value(status.MigrationVersion),
		Version:                types.StringValue(status.Version),
		StartupPath:            types.StringValue(status.StartupPath),
		AppData:                types.StringValue(status.AppData),
		OsName:                 types.StringValue(status.OsName),
		OsVersion:              types.StringValue(status.OsVersion),
		Mode:                   types.StringValue(status.Mode),
		Branch:                 types.StringValue(status.Branch),
		Authentication:         types.StringValue(status.Authentication),
		URLBase:                types.StringValue(status.URLBase),
		RuntimeVersion:         types.StringValue(status.RuntimeVersion),
		RuntimeName:            types.StringValue(status.RuntimeName),
		PackageVersion:         types.StringValue(status.PackageVersion),
		PackageAuthor:          types.StringValue(status.PackageAuthor),
		PackageUpdateMechanism: types.StringValue(status.PackageUpdateMechanism),
		BuildTime:              types.StringValue(status.BuildTime.String()),
		StartTime:              types.StringValue(status.StartTime.String()),
		AppName:                types.StringValue(status.AppName),
		DatabaseType:           types.StringValue(status.DatabaseType),
		DatabaseVersion:        types.StringValue(status.DatabaseVersion),
		InstanceName:           types.StringValue(status.InstanceName),
	}
}
