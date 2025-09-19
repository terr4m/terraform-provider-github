resource "github_team" "example" {
  organization = "example-org"
  name         = "example-team"
  description  = "An example team"
  privacy      = "closed"
}
