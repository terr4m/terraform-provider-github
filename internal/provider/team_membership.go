package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/google/go-github/v74/github"
)

var (
	_ resource.Resource              = &TeamMembershipResource{}
	_ resource.ResourceWithConfigure = &TeamMembershipResource{}
)

// NewTeamMembershipResource creates a new resource resource.
func NewTeamMembershipResource() resource.Resource {
	return &TeamMembershipResource{}
}

// TeamMembershipResource defines the resource implementation.
type TeamMembershipResource struct {
	providerData *GitHubProviderData
}

// TeamMembershipModel describes the data model.
type TeamMembershipModel struct {
	Organization types.String `tfsdk:"organization"`
	Role         types.String `tfsdk:"role"`
	State        types.String `tfsdk:"state"`
	Team         types.String `tfsdk:"team"`
	Username     types.String `tfsdk:"username"`
}

// Metadata returns the resource metadata.
func (d *TeamMembershipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_team_membership", req.ProviderTypeName)
}

// Schema returns the resource schema.
func (r *TeamMembershipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The _GitHub_ team membership resource (`github_team_membership`) allows you to manage membership for a _GitHub_ team.",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				MarkdownDescription: "Login of the organization the team belongs to.",
				Required:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The role of the membership. Can be `member` or `maintainer`.",
				Optional:            true,
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The state of the membership. Can be `active` or `pending`.",
				Computed:            true,
			},
			"team": schema.StringAttribute{
				MarkdownDescription: "Slug of the team.",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Login of the user to add to the team.",
				Required:            true,
			},
		},
	}
}

// Configure configures the resource.
func (r *TeamMembershipResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*GitHubProviderData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected resource provider data.", fmt.Sprintf("expected *GitHubProviderData, got: %T", req.ProviderData))
		return
	}

	r.providerData = providerData
}

// Create creates the resource.
func (r *TeamMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeamMembershipModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	client, err := r.providerData.ClientCreator.OrganizationClient(ctx, plan.Organization.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to create organization client", err.Error())
		return
	}

	_, response, _ := client.Teams.GetTeamMembershipBySlug(ctx, plan.Organization.ValueString(), plan.Team.ValueString(), plan.Username.ValueString())
	if response.StatusCode != 404 {
		resp.Diagnostics.AddError("Team membership already exists.", "can't add the same user to the same team multiple times")
		return
	}

	m, _, err := client.Teams.AddTeamMembershipBySlug(ctx, plan.Organization.ValueString(), plan.Team.ValueString(), plan.Username.ValueString(), &github.TeamAddTeamMembershipOptions{Role: plan.Role.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Failed to create team membership.", err.Error())
		return
	}

	state := TeamMembershipModel{
		Organization: plan.Organization,
		Role:         types.StringValue(m.GetRole()),
		State:        types.StringValue(m.GetState()),
		Team:         plan.Team,
		Username:     plan.Username,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read reads the resource state.
func (r *TeamMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeamMembershipModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	client, err := r.providerData.ClientCreator.OrganizationClient(ctx, state.Organization.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to create organization client", err.Error())
		return
	}

	m, _, err := client.Teams.GetTeamMembershipBySlug(ctx, state.Organization.ValueString(), state.Team.ValueString(), state.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get team membership.", err.Error())
		return
	}

	state.Role = types.StringValue(m.GetRole())
	state.State = types.StringValue(m.GetState())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource.
func (r *TeamMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeamMembershipModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	client, err := r.providerData.ClientCreator.OrganizationClient(ctx, plan.Organization.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to create organization client", err.Error())
		return
	}

	m, _, err := client.Teams.AddTeamMembershipBySlug(ctx, plan.Organization.ValueString(), plan.Team.ValueString(), plan.Username.ValueString(), &github.TeamAddTeamMembershipOptions{Role: plan.Role.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Failed to update team membership.", err.Error())
		return
	}

	state := TeamMembershipModel{
		Organization: plan.Organization,
		Role:         types.StringValue(m.GetRole()),
		State:        types.StringValue(m.GetState()),
		Team:         plan.Team,
		Username:     plan.Username,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete deletes the resource.
func (r *TeamMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeamMembershipModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	client, err := r.providerData.ClientCreator.OrganizationClient(ctx, state.Organization.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to create organization client", err.Error())
		return
	}

	_, err = client.Teams.RemoveTeamMembershipBySlug(ctx, state.Organization.ValueString(), state.Team.ValueString(), state.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete team membership.", err.Error())
		return
	}
}
