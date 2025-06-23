# Terraform Provider GitHub

![Terraform](https://img.shields.io/badge/Terraform-terr4m/github-purple?logo=terraform&link=https%3A%2F%2Fregistry.terraform.io%2Fproviders%2Fterr4m%github)
![GitHub Release (latest SemVer)](https://img.shields.io/github/v/release/terr4m/terraform-provider-github?logo=github&label=Release&sort=semver)
![Release](https://github.com/terr4m/terraform-provider-github/actions/workflows/release.yaml/badge.svg)
![Validate](https://github.com/terr4m/terraform-provider-github/actions/workflows/validate.yaml/badge.svg?branch=main)

This _Terraform_ provider allows you to interact with _GitHub_ using the REST API from within _Terraform_ and can be found on the _Terraform_ registry at [terr4m/github](https://registry.terraform.io/providers/terr4m/github/latest).

## Usage

For full documentation on how to use this provider, please see the [Terraform Registry](https://registry.terraform.io/providers/terr4m/github/latest/docs).

### Status

This provider is currently in active development as a replacement for the _"official"_ GitHub TF provider ([integrations/github](https://registry.terraform.io/providers/integrations/github/latest)) which appears to be no longer supported or taking contributions. Although this provider is not yet feature complete, it is being actively developed and maintained unlike `integrations/github`. In addition to the maintained status, this provider can make use of the latest TF patterns as it uses the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) instead of the older [Terraform SDK v2](https://developer.hashicorp.com/terraform/plugin/sdkv2).

## Architecture

This provider is build using the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) and the [google/go-github](https://github.com/google/go-github) SDK. The provider resources and data sources are not intended to be a direct mapping of the _GitHub_ REST API but instead are intended to be a more TF friendly abstraction of the API.

## Contributing

Contributions are welcome in all forms including documentation and issues as well as code; please see [CONTRIBUTING.md](./.github/CONTRIBUTING.md).
