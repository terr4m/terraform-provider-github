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
	t.Run("team_exists", func(t *testing.T) {
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
`, accTestConfigValues.Owner, accTestConfigValues.TeamSlug),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.github_team.test", tfjsonpath.New("name"), knownvalue.StringExact(accTestConfigValues.TeamSlug)),
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
`, accTestConfigValues.Owner),
					ExpectError: regexp.MustCompile("Error: Failed to get team"),
				},
			},
		})
	})
}
