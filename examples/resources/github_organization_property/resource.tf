resource "github_organization_property" "example" {
  organization  = "example-org"
  name          = "example-property"
  description   = "An example organization property"
  value_type    = "string"
  default_value = "example-value"
}
