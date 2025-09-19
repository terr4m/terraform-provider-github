package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/google/go-github/v74/github"
)

var (
	_ datasource.DataSource                     = &UserDataSource{}
	_ datasource.DataSourceWithConfigValidators = &UserDataSource{}
	_ datasource.DataSourceWithConfigure        = &UserDataSource{}
)

// NewUserDataSource creates a new user data source.
func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

// UserDataSource defines the data source implementation.
type UserDataSource struct {
	providerData *GitHubProviderData
}

// UserModel describes the data source data model.
type UserModel struct {
	Bio      types.String `tfsdk:"bio"`
	Company  types.String `tfsdk:"company"`
	Email    types.String `tfsdk:"email"`
	ID       types.Int64  `tfsdk:"id"`
	Location types.String `tfsdk:"location"`
	Login    types.String `tfsdk:"login"`
	Name     types.String `tfsdk:"name"`
}

// Metadata returns the data source metadata.
func (d *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_user", req.ProviderTypeName)
}

// Schema returns the data source schema.
func (d *UserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The _GitHub_ user data source (`github_user`) allows you to retrieve information about a user on _GitHub_.",
		Attributes: map[string]schema.Attribute{
			"bio": schema.StringAttribute{
				MarkdownDescription: "Bio of the user.",
				Computed:            true,
			},
			"company": schema.StringAttribute{
				MarkdownDescription: "Company of the user.",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Email of the user.",
				Computed:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "ID of the user.",
				Optional:            true,
				Computed:            true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "Location of the user.",
				Computed:            true,
			},
			"login": schema.StringAttribute{
				MarkdownDescription: "Login of the user.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the user.",
				Computed:            true,
			},
		},
	}
}

// ConfigValidators returns the data source config validators.
func (d *UserDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("login"),
		),
	}
}

// Configure configures the data source.
func (d *UserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.providerData.ClientCreator.DefaultClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create organization client", err.Error())
		return
	}

	var user *github.User
	if !data.ID.IsNull() {
		u, _, err := client.Users.GetByID(ctx, data.ID.ValueInt64())
		if err != nil {
			resp.Diagnostics.AddError("Failed to get user.", err.Error())
			return
		}
		user = u
	} else if !data.Login.IsNull() {
		u, _, err := client.Users.Get(ctx, data.Login.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Failed to get user.", err.Error())
			return
		}
		user = u
	}

	data.Bio = types.StringValue(user.GetBio())
	data.Company = types.StringValue(user.GetCompany())
	data.Email = types.StringValue(user.GetEmail())
	data.ID = types.Int64Value(user.GetID())
	data.Location = types.StringValue(user.GetLocation())
	data.Login = types.StringValue(user.GetLogin())
	data.Name = types.StringValue(user.GetName())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
