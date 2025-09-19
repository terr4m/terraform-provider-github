package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &OrganizationPropertiesDataSource{}
	_ datasource.DataSourceWithConfigure = &OrganizationPropertiesDataSource{}
)

// NewOrganizationPropertiesDataSource creates a new organization data source.
func NewOrganizationPropertiesDataSource() datasource.DataSource {
	return &OrganizationPropertiesDataSource{}
}

// OrganizationPropertiesDataSource defines the data source implementation.
type OrganizationPropertiesDataSource struct {
	providerData *GitHubProviderData
}

// OrganizationPropertiesModel describes the data source data model.
type OrganizationPropertiesModel struct {
	Properties   []PropertyModel `tfsdk:"properties"`
	Organization types.String    `tfsdk:"organization"`
}

// Metadata returns the data source metadata.
func (d *OrganizationPropertiesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_organization_properties", req.ProviderTypeName)
}

// Schema returns the data source schema.
func (d *OrganizationPropertiesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The _GitHub_ organization properties data source (`github_organization_properties`) allows you to retrieve information about a _GitHub_ organization's properties.",
		Attributes: map[string]schema.Attribute{
			"properties": schema.ListNestedAttribute{
				MarkdownDescription: "List of organization properties.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"allowed_values": schema.ListAttribute{
							MarkdownDescription: "List of allowed values for the property.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"default_value": schema.StringAttribute{
							MarkdownDescription: "Default value of the property.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the property.",
							Computed:            true,
						},
						"editable_by": schema.StringAttribute{
							MarkdownDescription: "Who can edit the property values.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the property.",
							Computed:            true,
						},
						"required": schema.BoolAttribute{
							MarkdownDescription: "Whether the property is required.",
							Computed:            true,
						},
						"source_type": schema.StringAttribute{
							MarkdownDescription: "Source type of the property.",
							Computed:            true,
						},
						"value_type": schema.StringAttribute{
							MarkdownDescription: "Value type of the property.",
							Computed:            true,
						},
					},
				},
			},
			"organization": schema.StringAttribute{
				MarkdownDescription: "Login of the organization the team belongs to.",
				Required:            true,
			},
		},
	}
}

// Configure configures the data source.
func (d *OrganizationPropertiesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *OrganizationPropertiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationPropertiesModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	organization := data.Organization.ValueString()

	client, err := d.providerData.ClientCreator.OrganizationClient(ctx, organization)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create organization client", err.Error())
		return
	}

	cp, _, err := client.Organizations.GetAllCustomProperties(ctx, organization)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get organization properties.", err.Error())
		return
	}

	props := make([]PropertyModel, 0, len(cp))
	for _, p := range cp {
		prop, diags := toPropertyModel(ctx, p)
		if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
			return
		}

		props = append(props, prop)
	}
	data.Properties = props

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
