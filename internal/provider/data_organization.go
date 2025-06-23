package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &OrganizationDataSource{}

// NewOrganizationDataSource creates a new organization data source.
func NewOrganizationDataSource() datasource.DataSource {
	return &OrganizationDataSource{}
}

// OrganizationDataSource defines the data source implementation.
type OrganizationDataSource struct {
	providerData *GitHubProviderData
}

// OrganizationModel describes the data source data model.
type OrganizationModel struct {
	Blog        types.String `tfsdk:"blog"`
	Company     types.String `tfsdk:"company"`
	Description types.String `tfsdk:"description"`
	Email       types.String `tfsdk:"email"`
	ID          types.Int64  `tfsdk:"id"`
	Location    types.String `tfsdk:"location"`
	Login       types.String `tfsdk:"login"`
	Name        types.String `tfsdk:"name"`
	Verified    types.Bool   `tfsdk:"verified"`
}

// Metadata returns the data source metadata.
func (d *OrganizationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_organization", req.ProviderTypeName)
}

// Schema returns the data source schema.
func (d *OrganizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The _GitHub_ organization data source (`github_organization`) allows you to retrieve information about a _GitHub_ organization.",
		Attributes: map[string]schema.Attribute{
			"blog": schema.StringAttribute{
				MarkdownDescription: "URL of the organization's website.",
				Computed:            true,
			},
			"company": schema.StringAttribute{
				MarkdownDescription: "Company name of the organization.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the organization.",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Email of the organization.",
				Computed:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "ID of the organization.",
				Computed:            true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "Location of the organization.",
				Computed:            true,
			},
			"login": schema.StringAttribute{
				MarkdownDescription: "Login of the organization.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the organization.",
				Computed:            true,
			},
			"verified": schema.BoolAttribute{
				MarkdownDescription: "Whether the organization is verified.",
				Computed:            true,
			},
		},
	}
}

// Configure configures the data source.
func (d *OrganizationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *OrganizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	o, _, err := d.providerData.Client.Organizations.Get(ctx, data.Login.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get organization.", err.Error())
		return
	}

	data.Blog = types.StringValue(o.GetBlog())
	data.Company = types.StringValue(o.GetCompany())
	data.Description = types.StringValue(o.GetDescription())
	data.Email = types.StringValue(o.GetEmail())
	data.ID = types.Int64Value(o.GetID())
	data.Location = types.StringValue(o.GetLocation())
	data.Login = types.StringValue(o.GetLogin())
	data.Name = types.StringValue(o.GetName())
	data.Verified = types.BoolValue(o.GetIsVerified())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
