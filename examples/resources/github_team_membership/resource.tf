data "github_team" "example" {
  organization = "example-org"
  slug         = "example-team"
}

data "github_user" "example" {
  username = "example-user"
}

resource "github_team_membership" "example" {
  organization = "example-org"
  team_id      = data.github_team.example.id
  username     = data.github_user.example.login
  role         = "member"
}
