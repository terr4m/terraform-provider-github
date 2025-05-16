package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccTeamMembersDataSource(t *testing.T) {
	t.Run("no_members", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
data "github_team_members" "test" {
  organization = "%s"
  team         = "%s"
}
`, accTestConfigValues.Owner, accTestConfigValues.TeamSlug),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.github_team_members.test", tfjsonpath.New("members"), knownvalue.ListSizeExact(0)),
					},
				},
			},
		})
	})

	t.Run("members", func(t *testing.T) {
		teamName := fmt.Sprintf("%s%s", accTestConfigValues.ResourcePrefix, acctest.RandomWithPrefix("test"))
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "github_team" "test" {
  organization = "%s"
  name         = "%s"
}

resource "github_team_membership" "test" {
  organization = "%[1]s"
  team         = github_team.test.slug
  username     = "%[3]s"
}

data "github_team_members" "test" {
  organization = "%[1]s"
  team         = github_team.test.slug

  depends_on = [github_team_membership.test]
}
`, accTestConfigValues.Owner, teamName, accTestConfigValues.Username),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.github_team_members.test", tfjsonpath.New("members"), knownvalue.ListSizeExact(1)),
						statecheck.ExpectKnownValue("data.github_team_members.test", tfjsonpath.New("members").AtSliceIndex(0).AtMapKey("username"), knownvalue.StringExact(accTestConfigValues.Username)),
					},
				},
			},
		})
	})
}
