# Copilot Instructions

This project is a Terraform Provider for GitHub built using the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework). The provider is intended to map to the GitHub REST API using [google/go-github](https://github.com/google/go-github).

## Dependencies

- Use [google/go-github](https://github.com/google/go-github) to interact with the GitHub REST API.

## Coding Standards

- Follow idiomatic Go coding standards.
- Prefer table tests over sub-tests.
- Use [google/go-cmp](https://github.com/google/go-cmp) for test comparisons.
