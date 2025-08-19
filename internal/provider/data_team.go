package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &TeamDataSource{}
	_ datasource.DataSourceWithConfigure = &TeamDataSource{}
)

// NewTeamDataSource creates a new team data source.
func NewTeamDataSource() datasource.DataSource {
	return &TeamDataSource{}
}

// TeamDataSource defines the data source implementation.
type TeamDataSource struct {
	providerData *GitHubProviderData
}

// Metadata returns the data source metadata.
func (d *TeamDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_team", req.ProviderTypeName)
}

// Schema returns the data source schema.
func (d *TeamDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The _GitHub_ team data source (`github_team`) allows you to retrieve information about a _GitHub_ team.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the team",
				Computed:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Unique identifier of the team",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the team.",
				Computed:            true,
			},
			"notifications": schema.BoolAttribute{
				MarkdownDescription: "If team members receive notifications when the team is `@mentioned`.",
				Computed:            true,
			},
			"organization": schema.StringAttribute{
				MarkdownDescription: "Name of the organization the team belongs to.",
				Required:            true,
			},
			"parent": schema.SingleNestedAttribute{
				MarkdownDescription: "Parent team of the team.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						MarkdownDescription: "Unique identifier of the parent team.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Name of the parent team.",
						Computed:            true,
					},
					"slug": schema.StringAttribute{
						MarkdownDescription: "Slug of the parent team name.",
						Computed:            true,
					},
				},
			},
			"privacy": schema.StringAttribute{
				MarkdownDescription: "The level of privacy this team should have. This can be one of `closed` or `secret`.",
				Computed:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "Slug of the team name.",
				Required:            true,
			},
		},
	}
}

// Configure configures the data source.
func (d *TeamDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*GitHubProviderData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected data source provider data.", fmt.Sprintf("expected *provider.GitHubProviderData, got: %T", req.ProviderData))
		return
	}

	d.providerData = providerData
}

// Read reads the data source.
func (d *TeamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TeamModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.providerData.ClientCreator.OrganizationClient(ctx, data.Organization.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to create organization client", err.Error())
		return
	}

	t, _, err := client.Teams.GetTeamBySlug(ctx, data.Organization.ValueString(), data.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get team.", err.Error())
		return
	}

	data.Description = types.StringValue(t.GetDescription())
	data.ID = types.Int64Value(t.GetID())
	data.Name = types.StringValue(t.GetName())
	data.Notifications = types.BoolValue(t.GetNotificationSetting() == TeamNotificationsEnabled)
	data.Organization = types.StringValue(t.GetOrganization().GetLogin())
	data.Privacy = types.StringValue(t.GetPrivacy())
	data.Slug = types.StringValue(t.GetSlug())

	if parent := t.GetParent(); parent != nil {
		data.Parent = &TeamModel{
			ID:   types.Int64Value(parent.GetID()),
			Name: types.StringValue(parent.GetName()),
			Slug: types.StringValue(parent.GetSlug()),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
