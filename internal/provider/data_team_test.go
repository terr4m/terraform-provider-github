package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccTeamDataSource(t *testing.T) {
	if accTestConfigData.AuthType == accAuthTypeUnauthenticated || !accTestConfigData.Features.Organization {
		t.Skip("Skipping test because the organization testing feature isn't enabled")
	}

	t.Run("team_exists", func(t *testing.T) {
		if len(accTestConfigData.Values.TeamSlug) == 0 {
			t.Skip("Skipping test because the team slug is not set")
		}

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
data "github_team" "test" {
  organization = "%s"
  slug         = "%s"
}
`, accTestConfigData.Values.Organization, accTestConfigData.Values.TeamSlug),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.github_team.test", tfjsonpath.New("name"), knownvalue.StringExact(accTestConfigData.Values.TeamSlug)),
					},
				},
			},
		})
	})

	t.Run("team_does_not_exist", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
data "github_team" "test" {
  organization = "%s"
  slug         = "should-not-exist"
}
`, accTestConfigData.Values.Organization),
					ExpectError: regexp.MustCompile("Error: Failed to get team"),
				},
			},
		})
	})
}
