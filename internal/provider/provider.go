package provider

import (
	"context"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/providervalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/terr4m/terraform-provider-github/internal/ghutil"
)

// Ensure GitHubProvider satisfies various provider interfaces.
var (
	_ provider.Provider                       = &GitHubProvider{}
	_ provider.ProviderWithConfigValidators   = &GitHubProvider{}
	_ provider.ProviderWithFunctions          = &GitHubProvider{}
	_ provider.ProviderWithEphemeralResources = &GitHubProvider{}
)

// New returns a new provider implementation.
func New(version, commit string) func() provider.Provider {
	return func() provider.Provider {
		return &GitHubProvider{
			version: version,
			commit:  commit,
		}
	}
}

// GitHubProviderData is the data available to the resource and data sources.
type GitHubProviderData struct {
	provider        *GitHubProvider
	Model           *GitHubProviderModel
	ClientCreator   ghutil.ClientCreator
	DefaultTimeouts *Timeouts
}

// Timeouts represents a set of timeouts.
type Timeouts struct {
	Create time.Duration
	Read   time.Duration
	Update time.Duration
	Delete time.Duration
}

// GitHubProviderModel describes the provider data model.
type GitHubProviderModel struct {
	AppAuth       *AppAuthModel  `tfsdk:"app_auth"`
	CacheRequests types.Bool     `tfsdk:"cache_requests"`
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
	Token         types.String   `tfsdk:"token"`
}

// AppAuth describes the application authentication configuration.
type AppAuthModel struct {
	ID             types.Int64  `tfsdk:"id"`
	PrivateKey     types.String `tfsdk:"private_key"`
	PrivateKeyFile types.String `tfsdk:"private_key_file"`
}

// GitHubProvider defines the provider implementation.
type GitHubProvider struct {
	version string
	commit  string
}

// Metadata returns the provider metadata.
func (p *GitHubProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "github"
	resp.Version = p.version
}

// Schema returns the provider schema.
func (p *GitHubProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The GitHub provider provides a way to manage _GitHub_ resources available via the REST API using _Terraform_.",
		Attributes: map[string]schema.Attribute{
			"app_auth": schema.SingleNestedAttribute{
				MarkdownDescription: "GitHub application authentication configuration; this is mutually exclusive with `token`. If `private_key` or `private_key_file` are not provided, the provider will attempt to use the `GITHUB_APP_PRIVATE_KEY` and then `GITHUB_APP_PRIVATE_KEY_FILE` environment variables.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						MarkdownDescription: "The GitHub application ID.",
						Required:            true,
					},
					"private_key": schema.StringAttribute{
						MarkdownDescription: "The private key for the GitHub application; this is mutually exclusive with `private_key_file`.",
						Optional:            true,
					},
					"private_key_file": schema.StringAttribute{
						MarkdownDescription: "The file containing the private key for the GitHub application; this is mutually exclusive with `private_key`.",
						Optional:            true,
					},
				},
			},
			"cache_requests": schema.BoolAttribute{
				MarkdownDescription: "If `true`, the provider will cache requests to the GitHub API. This can help reduce the number of requests made to the API, but may result in stale data being returned. Defaults to `false`.",
				Optional:            true,
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create:            true,
				CreateDescription: "Timeout for resource creation; defaults to `10m`. This should be a string that can be [parsed as a duration] (https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).",
				Read:              true,
				ReadDescription:   "Timeout for resource or data source reads; defaults to `10m`. This should be a string that can be [parsed as a duration] (https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).",
				Update:            true,
				UpdateDescription: "Timeout for resource update; defaults to `10m`. This should be a string that can be [parsed as a duration] (https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).",
				Delete:            true,
				DeleteDescription: "Timeout for resource deletion; defaults to `10m`. This should be a string that can be [parsed as a duration] (https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).",
			}),
			"token": schema.StringAttribute{
				MarkdownDescription: "A GitHub token to use for authentication; this is mutually exclusive with `app_auth`. If `app_auth` isn;t configured and this isn't set the provider will look for the `GITHUB_TOKEN` environment variable.",
				Optional:            true,
			},
		},
	}
}

// ConfigValidators returns the provider config validators.
func (p *GitHubProvider) ConfigValidators(ctx context.Context) []provider.ConfigValidator {
	return []provider.ConfigValidator{
		providervalidator.Conflicting(
			path.MatchRoot("app_auth"),
			path.MatchRoot("token"),
		),
		providervalidator.Conflicting(
			path.MatchRoot("app_auth").AtName("private_key"),
			path.MatchRoot("app_auth").AtName("private_key_file"),
		),
	}
}

// Configure configures the provider.
func (p *GitHubProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	if req.ClientCapabilities.DeferralAllowed && !req.Config.Raw.IsFullyKnown() {
		resp.Deferred = &provider.Deferred{
			Reason: provider.DeferredReasonProviderConfigUnknown,
		}
	}

	model := &GitHubProviderModel{}
	if resp.Diagnostics.Append(req.Config.Get(ctx, model)...); resp.Diagnostics.HasError() {
		return
	}

	var clientCreator ghutil.ClientCreator
	cacheRequests := model.CacheRequests.ValueBool()
	if model.AppAuth != nil {
		var privateKey []byte
		appID := model.AppAuth.ID.ValueInt64()

		if !model.AppAuth.PrivateKey.IsNull() {
			privateKey = []byte(model.AppAuth.PrivateKey.ValueString())
		} else if !model.AppAuth.PrivateKeyFile.IsNull() {
			k, err := os.ReadFile(model.AppAuth.PrivateKeyFile.String())
			if err != nil {
				resp.Diagnostics.AddError("Failed to read private key file", err.Error())
				return
			}
			privateKey = k
		} else if v := os.Getenv("GITHUB_APP_PRIVATE_KEY"); len(v) != 0 {
			privateKey = []byte(v)
		} else if v := os.Getenv("GITHUB_APP_PRIVATE_KEY_FILE"); len(v) != 0 {
			k, err := os.ReadFile(v)
			if err != nil {
				resp.Diagnostics.AddError("Failed to read private key file", err.Error())
				return
			}
			privateKey = k
		} else {
			resp.Diagnostics.AddError("Private key not provided", "no private key was provided for the app auth")
			return
		}

		cc, err := ghutil.NewAppClientCreator(appID, privateKey, 10, cacheRequests)
		if err != nil {
			resp.Diagnostics.AddError("Failed to create GitHub client creator", err.Error())
			return
		}
		clientCreator = cc
	} else {
		var token *string

		if !model.Token.IsNull() {
			token = model.Token.ValueStringPointer()
		} else if v := os.Getenv("GITHUB_TOKEN"); len(v) != 0 {
			token = &v
		}

		cc, err := ghutil.NewClientCreator(token, cacheRequests)
		if err != nil {
			resp.Diagnostics.AddError("Failed to create GitHub client creator", err.Error())
			return
		}
		clientCreator = cc
	}

	createTimeout, diags := model.Timeouts.Create(ctx, 10*time.Minute)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}
	readTimeout, diags := model.Timeouts.Read(ctx, 10*time.Minute)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}
	updateTimeout, diags := model.Timeouts.Update(ctx, 10*time.Minute)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}
	deleteTimeout, diags := model.Timeouts.Delete(ctx, 10*time.Minute)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	providerData := &GitHubProviderData{
		provider:      p,
		Model:         model,
		ClientCreator: clientCreator,
		DefaultTimeouts: &Timeouts{
			Create: createTimeout,
			Read:   readTimeout,
			Update: updateTimeout,
			Delete: deleteTimeout,
		},
	}

	resp.DataSourceData = providerData
	resp.ResourceData = providerData
}

// Resources returns the provider resources.
func (p *GitHubProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewOrganizationPropertyResource,
		NewTeamMembershipResource,
		NewTeamResource,
	}
}

// EphemeralResources returns the provider ephemeral resources.
func (p *GitHubProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

// DataSources returns the provider data sources.
func (p *GitHubProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOrganizationDataSource,
		NewOrganizationPropertiesDataSource,
		NewTeamDataSource,
		NewTeamMembersDataSource,
		NewUserDataSource,
	}
}

// Functions returns the provider functions.
func (p *GitHubProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}
