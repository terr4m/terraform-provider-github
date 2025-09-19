data "github_team" "example" {
  organization = "example-org"
  slug         = "example-team"
}

data "github_team_members" "example" {
  organization = "example-org"
  team_id      = data.github_team.example.id
}
