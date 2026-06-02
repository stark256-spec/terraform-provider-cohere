package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.ProviderWithFunctions = &CohereProvider{}

type CohereProvider struct{ version string }

type CohereProviderModel struct {
	APIKey  types.String `tfsdk:"api_key"`
	BaseURL types.String `tfsdk:"base_url"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider { return &CohereProvider{version: version} }
}

func (p *CohereProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cohere"
	resp.Version = p.version
}

func (p *CohereProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Cohere organization resources — teams and API keys.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key. Can be set via `COHERE_API_KEY` env var.",
				Optional:            true,
				Sensitive:           true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Override the API base URL.",
				Optional:            true,
			},
		},
	}
}

func (p *CohereProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg CohereProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiKey := os.Getenv("COHERE_API_KEY")
	if !cfg.APIKey.IsNull() {
		apiKey = cfg.APIKey.ValueString()
	}
	if apiKey == "" {
		resp.Diagnostics.AddError("Missing API key", "Set api_key or COHERE_API_KEY env var.")
		return
	}
	baseURL := ""
	if !cfg.BaseURL.IsNull() {
		baseURL = cfg.BaseURL.ValueString()
	}
	client := newClient(apiKey, baseURL)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *CohereProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{NewWorkspaceResource, NewAPIKeyResource}
}

func (p *CohereProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *CohereProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{}
}
