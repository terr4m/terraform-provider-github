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

func TestAccTeamMembershipResource(t *testing.T) {
	if accTestConfigData.AuthType == accAuthTypeUnauthenticated || !accTestConfigData.Features.Organization {
		t.Skip("Skipping test because the organization testing feature isn't enabled")
	}

	t.Run("create_default", func(t *testing.T) {
		teamName := fmt.Sprintf("%s%s", accTestConfigData.ResourcePrefix, acctest.RandomWithPrefix("test"))
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
`, accTestConfigData.Values.Organization, teamName, accTestConfigData.Values.Username),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("github_team_membership.test", tfjsonpath.New("organization"), knownvalue.StringExact(accTestConfigData.Values.Organization)),
						statecheck.ExpectKnownValue("github_team_membership.test", tfjsonpath.New("role"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue("github_team_membership.test", tfjsonpath.New("state"), knownvalue.StringExact("active")),
						statecheck.ExpectKnownValue("github_team_membership.test", tfjsonpath.New("team"), knownvalue.StringExact(teamName)),
						statecheck.ExpectKnownValue("github_team_membership.test", tfjsonpath.New("username"), knownvalue.StringExact(accTestConfigData.Values.Username)),
					},
				},
			},
		})
	})
}
