package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/google/go-github/v74/github"
)

var (
	_ resource.Resource              = &TeamResource{}
	_ resource.ResourceWithConfigure = &TeamResource{}
)

// NewTeamResource creates a new resource resource.
func NewTeamResource() resource.Resource {
	return &TeamResource{}
}

// TeamResource defines the resource implementation.
type TeamResource struct {
	providerData *GitHubProviderData
}

// TeamModel describes the data model.
type TeamModel struct {
	Description   types.String `tfsdk:"description"`
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Notifications types.Bool   `tfsdk:"notifications"`
	Organization  types.String `tfsdk:"organization"`
	Parent        *TeamModel   `tfsdk:"parent"`
	Privacy       types.String `tfsdk:"privacy"`
	Slug          types.String `tfsdk:"slug"`
}

const (
	TeamNotificationsEnabled  = "notifications_enabled"
	TeamNotificationsDisabled = "notifications_disabled"
)

// Metadata returns the resource metadata.
func (d *TeamResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_team", req.ProviderTypeName)
}

// Schema returns the resource schema.
func (r *TeamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The _GitHub_ team resource (`github_team`) allows you to manage teams for a _GitHub_ organization.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the team.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Unique identifier of the team.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the team.",
				Required:            true,
			},
			"notifications": schema.BoolAttribute{
				MarkdownDescription: "If team members receive notifications when the team is `@mentioned`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"organization": schema.StringAttribute{
				MarkdownDescription: "Name of the organization the team belongs to.",
				Required:            true,
			},
			"parent": schema.SingleNestedAttribute{
				MarkdownDescription: "Parent team of the team.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						MarkdownDescription: "Unique identifier of the parent team.",
						Required:            true,
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
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("closed"),
				Validators: []validator.String{
					stringvalidator.OneOf("closed", "secret"),
				},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "Slug of the team name.",
				Computed:            true,
			},
		},
	}
}

// Configure configures the resource.
func (r *TeamResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *TeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeamModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	client, err := r.providerData.ClientCreator.OrganizationClient(ctx, plan.Organization.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to create organization client", err.Error())
		return
	}

	n := github.NewTeam{
		Description: github.Ptr(plan.Description.ValueString()),
		Name:        plan.Name.ValueString(),
		Privacy:     github.Ptr(plan.Privacy.ValueString()),
	}

	if plan.Notifications.ValueBool() {
		n.NotificationSetting = github.Ptr(TeamNotificationsEnabled)
	} else {
		n.NotificationSetting = github.Ptr(TeamNotificationsDisabled)
	}

	if plan.Parent != nil {
		n.ParentTeamID = github.Ptr(plan.Parent.ID.ValueInt64())
	}

	t, _, err := client.Teams.CreateTeam(ctx, plan.Organization.ValueString(), n)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create team.", err.Error())
		return
	}

	if t.GetMembersCount() > 0 {
		m, _, err := client.Teams.ListTeamMembersBySlug(ctx, plan.Organization.ValueString(), t.GetSlug(), &github.TeamListTeamMembersOptions{})
		if err != nil {
			resp.Diagnostics.AddError("Failed to get team members.", err.Error())
			return
		}

		for _, member := range m {
			_, err := client.Teams.RemoveTeamMembershipBySlug(ctx, plan.Organization.ValueString(), t.GetSlug(), member.GetLogin())
			if err != nil {
				resp.Diagnostics.AddError("Failed to remove team member.", err.Error())
				return
			}
		}
	}

	state := TeamModel{
		Description:   types.StringValue(t.GetDescription()),
		ID:            types.Int64Value(t.GetID()),
		Name:          types.StringValue(t.GetName()),
		Notifications: types.BoolValue(t.GetNotificationSetting() == TeamNotificationsEnabled),
		Organization:  types.StringValue(t.GetOrganization().GetLogin()),
		Privacy:       types.StringValue(t.GetPrivacy()),
		Slug:          types.StringValue(t.GetSlug()),
	}

	if parent := t.GetParent(); parent != nil {
		state.Parent = &TeamModel{
			ID:   types.Int64Value(parent.GetID()),
			Name: types.StringValue(parent.GetName()),
			Slug: types.StringValue(parent.GetSlug()),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read reads the resource state.
func (r *TeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeamModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	client, err := r.providerData.ClientCreator.OrganizationClient(ctx, state.Organization.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to create organization client", err.Error())
		return
	}

	t, _, err := client.Teams.GetTeamBySlug(ctx, state.Organization.ValueString(), state.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get team.", err.Error())
		return
	}

	state.Description = types.StringValue(t.GetDescription())
	state.ID = types.Int64Value(t.GetID())
	state.Name = types.StringValue(t.GetName())
	state.Notifications = types.BoolValue(t.GetNotificationSetting() == TeamNotificationsEnabled)
	state.Organization = types.StringValue(t.GetOrganization().GetLogin())
	state.Privacy = types.StringValue(t.GetPrivacy())
	state.Slug = types.StringValue(t.GetSlug())

	if parent := t.GetParent(); parent != nil {
		state.Parent = &TeamModel{
			ID:   types.Int64Value(parent.GetID()),
			Name: types.StringValue(parent.GetName()),
			Slug: types.StringValue(parent.GetSlug()),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource.
func (r *TeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeamModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	client, err := r.providerData.ClientCreator.OrganizationClient(ctx, plan.Organization.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to create organization client", err.Error())
		return
	}

	n := github.NewTeam{
		Description: github.Ptr(plan.Description.ValueString()),
		Name:        plan.Name.ValueString(),
		Privacy:     github.Ptr(plan.Privacy.ValueString()),
	}

	if plan.Notifications.ValueBool() {
		n.NotificationSetting = github.Ptr(TeamNotificationsEnabled)
	} else {
		n.NotificationSetting = github.Ptr(TeamNotificationsDisabled)
	}

	if plan.Parent != nil {
		n.ParentTeamID = github.Ptr(plan.Parent.ID.ValueInt64())
	}

	t, _, err := client.Teams.EditTeamBySlug(ctx, plan.Organization.ValueString(), plan.Slug.ValueString(), n, plan.Parent == nil)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update team.", err.Error())
		return
	}

	state := TeamModel{
		Description:   types.StringValue(t.GetDescription()),
		ID:            types.Int64Value(t.GetID()),
		Name:          types.StringValue(t.GetName()),
		Notifications: types.BoolValue(t.GetNotificationSetting() == TeamNotificationsEnabled),
		Organization:  types.StringValue(t.GetOrganization().GetLogin()),
		Privacy:       types.StringValue(t.GetPrivacy()),
		Slug:          types.StringValue(t.GetSlug()),
	}

	if parent := t.GetParent(); parent != nil {
		state.Parent = &TeamModel{
			ID:   types.Int64Value(parent.GetID()),
			Name: types.StringValue(parent.GetName()),
			Slug: types.StringValue(parent.GetSlug()),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete deletes the resource.
func (r *TeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeamModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	client, err := r.providerData.ClientCreator.OrganizationClient(ctx, state.Organization.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to create organization client", err.Error())
		return
	}

	_, err = client.Teams.DeleteTeamBySlug(ctx, state.Organization.ValueString(), state.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete team.", err.Error())
		return
	}
}
