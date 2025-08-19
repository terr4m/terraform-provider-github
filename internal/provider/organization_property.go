package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/google/go-github/v74/github"
)

var (
	_ resource.Resource                = &OrganizationPropertyResource{}
	_ resource.ResourceWithConfigure   = &OrganizationPropertyResource{}
	_ resource.ResourceWithImportState = &OrganizationPropertyResource{}
)

// NewOrganizationPropertyResource creates a new OrganizationPropertyResource.
func NewOrganizationPropertyResource() resource.Resource {
	return &OrganizationPropertyResource{}
}

// OrganizationPropertyResource defines the resource implementation.
type OrganizationPropertyResource struct {
	providerData *GitHubProviderData
}

// PropertyModel describes the data model.
type PropertyModel struct {
	AllowedValues types.List   `tfsdk:"allowed_values"`
	DefaultValue  types.String `tfsdk:"default_value"`
	Description   types.String `tfsdk:"description"`
	EditableBy    types.String `tfsdk:"editable_by"`
	Name          types.String `tfsdk:"name"`
	Required      types.Bool   `tfsdk:"required"`
	SourceType    types.String `tfsdk:"source_type"`
	ValueType     types.String `tfsdk:"value_type"`
}

// OrganizationPropertyModel describes the data model.
type OrganizationPropertyModel struct {
	Organization types.String `tfsdk:"organization"`
	PropertyModel
}

// Metadata returns the resource metadata.
func (d *OrganizationPropertyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_organization_property", req.ProviderTypeName)
}

// Schema returns the resource schema.
func (r *OrganizationPropertyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The _GitHub_ organization property resource (`github_organization_property`) allows you to manage custom properties for a _GitHub_ organization.",
		Attributes: map[string]schema.Attribute{
			"allowed_values": schema.ListAttribute{
				MarkdownDescription: "An ordered list of the allowed values of the property; the property can have up to 200 allowed values.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Validators: []validator.List{
					listvalidator.SizeAtMost(200),
				},
			},
			"default_value": schema.StringAttribute{
				MarkdownDescription: "Default value of the property.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Short description of the property.",
				Optional:            true,
			},
			"editable_by": schema.StringAttribute{
				MarkdownDescription: "Who can edit the values of the property.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("org_actors"),
				Validators: []validator.String{
					stringvalidator.OneOf("org_actors", "org_and_repo_actors"),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the property.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"organization": schema.StringAttribute{
				MarkdownDescription: "Name of the organization.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"required": schema.BoolAttribute{
				MarkdownDescription: "Whether the property is required.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"source_type": schema.StringAttribute{
				MarkdownDescription: "The source type of the property.",
				Computed:            true,
			},
			"value_type": schema.StringAttribute{
				MarkdownDescription: "The type of the value for the property.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("string", "single_select", "multi_select", "true_false"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Configure configures the resource.
func (r *OrganizationPropertyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *OrganizationPropertyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model OrganizationPropertyModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	organization := model.Organization.ValueString()

	client, err := r.providerData.ClientCreator.OrganizationClient(ctx, organization)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create organization client", err.Error())
		return
	}

	property, diags := fromPropertyModel(ctx, model.PropertyModel)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	p, _, err := client.Organizations.CreateOrUpdateCustomProperty(ctx, organization, model.Name.ValueString(), &property)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create organization property.", err.Error())
		return
	}

	m, diags := toOrganizationPropertyModel(ctx, organization, p)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}

// Read reads the resource state.
func (r *OrganizationPropertyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model OrganizationPropertyModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	organization := model.Organization.ValueString()

	client, err := r.providerData.ClientCreator.OrganizationClient(ctx, organization)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create organization client", err.Error())
		return
	}

	p, _, err := client.Organizations.GetCustomProperty(ctx, organization, model.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get organization property.", err.Error())
		return
	}

	model, diags := toOrganizationPropertyModel(ctx, organization, p)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// Update updates the resource.
func (r *OrganizationPropertyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model OrganizationPropertyModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	organization := model.Organization.ValueString()

	client, err := r.providerData.ClientCreator.OrganizationClient(ctx, organization)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create organization client", err.Error())
		return
	}

	property, diags := fromPropertyModel(ctx, model.PropertyModel)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	p, _, err := client.Organizations.CreateOrUpdateCustomProperty(ctx, organization, model.Name.ValueString(), &property)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update organization property.", err.Error())
		return
	}

	m, diags := toOrganizationPropertyModel(ctx, organization, p)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}

// Delete deletes the resource.
func (r *OrganizationPropertyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model OrganizationPropertyModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	organization := model.Organization.ValueString()

	client, err := r.providerData.ClientCreator.OrganizationClient(ctx, organization)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create organization client", err.Error())
		return
	}

	_, err = client.Organizations.RemoveCustomProperty(ctx, organization, model.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete organization property.", err.Error())
		return
	}
}

// ImportState imports the resource state.
func (r *OrganizationPropertyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID.", "import id must be in the format \"organization:property_name\"")
		return
	}

	organization := parts[0]
	propertyName := parts[1]

	if len(organization) == 0 {
		resp.Diagnostics.AddError("Invalid import ID.", "organization must be non-empty")
		return
	}

	if len(propertyName) == 0 || propertyName == "" {
		resp.Diagnostics.AddError("Invalid import ID.", "property_name must be non-empty")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization"), organization)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), propertyName)...)
}

func toOrganizationPropertyModel(ctx context.Context, org string, p *github.CustomProperty) (OrganizationPropertyModel, diag.Diagnostics) {
	pm, diags := toPropertyModel(ctx, p)
	if diags.HasError() {
		return OrganizationPropertyModel{}, diags
	}

	m := OrganizationPropertyModel{
		Organization:  types.StringValue(org),
		PropertyModel: pm,
	}

	return m, diag.Diagnostics{}
}

func toPropertyModel(ctx context.Context, p *github.CustomProperty) (PropertyModel, diag.Diagnostics) {
	if p == nil {
		diags := diag.Diagnostics{}
		diags.AddError("Failed to convert to property model.", "property is nil")
		return PropertyModel{}, diags
	}

	allowedValues, diags := types.ListValueFrom(ctx, types.StringType, p.AllowedValues)
	if diags.HasError() {
		return PropertyModel{}, diags
	}

	m := PropertyModel{
		AllowedValues: allowedValues,
		DefaultValue:  types.StringPointerValue(p.DefaultValue),
		Description:   types.StringPointerValue(p.Description),
		Name:          types.StringValue(p.GetPropertyName()),
		Required:      types.BoolPointerValue(p.Required),
		SourceType:    types.StringPointerValue(p.SourceType),
		EditableBy:    types.StringPointerValue(p.ValuesEditableBy),
		ValueType:     types.StringValue(p.ValueType),
	}

	return m, diag.Diagnostics{}
}

func fromPropertyModel(ctx context.Context, m PropertyModel) (github.CustomProperty, diag.Diagnostics) {
	p := github.CustomProperty{
		DefaultValue: m.DefaultValue.ValueStringPointer(),
		Description:  m.Description.ValueStringPointer(),
		PropertyName: github.Ptr(m.Name.ValueString()),
		Required:     m.Required.ValueBoolPointer(),
		// SourceType:       m.SourceType.ValueStringPointer(),
		ValuesEditableBy: m.EditableBy.ValueStringPointer(),
		ValueType:        m.ValueType.ValueString(),
	}

	if !m.AllowedValues.IsNull() && !m.AllowedValues.IsUnknown() {
		allowedValues := make([]string, 0, len(m.AllowedValues.Elements()))
		if diags := m.AllowedValues.ElementsAs(ctx, &allowedValues, false); diags.HasError() {
			return github.CustomProperty{}, diags
		}
		p.AllowedValues = allowedValues
	}

	return p, diag.Diagnostics{}
}
