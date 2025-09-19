package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/google/go-github/v74/github"
)

var (
	_ datasource.DataSource              = &TeamMembersDataSource{}
	_ datasource.DataSourceWithConfigure = &TeamMembersDataSource{}
)

// NewTeamMembersDataSource creates a new team members data source.
func NewTeamMembersDataSource() datasource.DataSource {
	return &TeamMembersDataSource{}
}

// TeamMembersDataSource defines the data source implementation.
type TeamMembersDataSource struct {
	providerData *GitHubProviderData
}

// TeamMembersModel describes the data model.
type TeamMembersModel struct {
	Members      []TeamMemberModel `tfsdk:"members"`
	Organization types.String      `tfsdk:"organization"`
	Team         types.String      `tfsdk:"team"`
}

// TeamMemberModel describes the data model.
type TeamMemberModel struct {
	Role     types.String `tfsdk:"role"`
	Username types.String `tfsdk:"username"`
}

// Metadata returns the data source metadata.
func (d *TeamMembersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_team_members", req.ProviderTypeName)
}

// Schema returns the data source schema.
func (d *TeamMembersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The _GitHub_ team members data source (`github_team_members`) allows you to retrieve information about a _GitHub_ team's members.",
		Attributes: map[string]schema.Attribute{
			"members": schema.ListNestedAttribute{
				MarkdownDescription: "List of active team members.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"role": schema.StringAttribute{
							MarkdownDescription: "Role of the member. Can be `member` or `maintainer`.",
							Computed:            true,
						},
						"username": schema.StringAttribute{
							MarkdownDescription: "Username of the member.",
							Computed:            true,
						},
					},
				},
			},
			"organization": schema.StringAttribute{
				MarkdownDescription: "Login of the organization the team belongs to.",
				Required:            true,
			},
			"team": schema.StringAttribute{
				MarkdownDescription: "Slug of the team.",
				Required:            true,
			},
		},
	}
}

// Configure configures the data source.
func (d *TeamMembersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *TeamMembersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TeamMembersModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.providerData.ClientCreator.OrganizationClient(ctx, data.Organization.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to create organization client", err.Error())
		return
	}

	mt, _, err := client.Teams.ListTeamMembersBySlug(ctx, data.Organization.ValueString(), data.Team.ValueString(), &github.TeamListTeamMembersOptions{Role: "maintainer"})
	if err != nil {
		resp.Diagnostics.AddError("Failed to get team maintainers.", err.Error())
		return
	}

	mb, _, err := client.Teams.ListTeamMembersBySlug(ctx, data.Organization.ValueString(), data.Team.ValueString(), &github.TeamListTeamMembersOptions{Role: "member"})
	if err != nil {
		resp.Diagnostics.AddError("Failed to get team members.", err.Error())
		return
	}

	members := make([]TeamMemberModel, 0, len(mt)+len(mb))

	for _, user := range mt {
		members = append(data.Members, TeamMemberModel{
			Role:     types.StringValue("maintainer"),
			Username: types.StringValue(user.GetLogin()),
		})
	}

	for _, user := range mb {
		members = append(data.Members, TeamMemberModel{
			Role:     types.StringValue("member"),
			Username: types.StringValue(user.GetLogin()),
		})
	}

	data.Members = members

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
