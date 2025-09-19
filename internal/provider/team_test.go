package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccTeamResource(t *testing.T) {
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
`, accTestConfigData.Values.Organization, teamName),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("github_team.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
						statecheck.ExpectKnownValue("github_team.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue("github_team.test", tfjsonpath.New("name"), knownvalue.StringExact(teamName)),
						statecheck.ExpectKnownValue("github_team.test", tfjsonpath.New("notifications"), knownvalue.Bool(true)),
						statecheck.ExpectKnownValue("github_team.test", tfjsonpath.New("organization"), knownvalue.StringExact(accTestConfigData.Values.Organization)),
						statecheck.ExpectKnownValue("github_team.test", tfjsonpath.New("parent"), knownvalue.Null()),
						statecheck.ExpectKnownValue("github_team.test", tfjsonpath.New("privacy"), knownvalue.StringExact("closed")),
						statecheck.ExpectKnownValue("github_team.test", tfjsonpath.New("slug"), knownvalue.StringExact(teamName)),
					},
				},
			},
		})
	})

	t.Run("create_full", func(t *testing.T) {
		teamName := fmt.Sprintf("%s%s", accTestConfigData.ResourcePrefix, acctest.RandomWithPrefix("test"))
		description := "test description"
		privacy := "secret"
		notifications := false

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "github_team" "test" {
  organization = "%s"
  name         = "%s"
	description  = "%s"
	privacy      = "%s"
	notifications = %v
}
`, accTestConfigData.Values.Organization, teamName, description, privacy, notifications),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("github_team.test", tfjsonpath.New("description"), knownvalue.StringExact(description)),
						statecheck.ExpectKnownValue("github_team.test", tfjsonpath.New("id"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue("github_team.test", tfjsonpath.New("name"), knownvalue.StringExact(teamName)),
						statecheck.ExpectKnownValue("github_team.test", tfjsonpath.New("notifications"), knownvalue.Bool(notifications)),
						statecheck.ExpectKnownValue("github_team.test", tfjsonpath.New("organization"), knownvalue.StringExact(accTestConfigData.Values.Organization)),
						statecheck.ExpectKnownValue("github_team.test", tfjsonpath.New("parent"), knownvalue.Null()),
						statecheck.ExpectKnownValue("github_team.test", tfjsonpath.New("privacy"), knownvalue.StringExact(privacy)),
						statecheck.ExpectKnownValue("github_team.test", tfjsonpath.New("slug"), knownvalue.StringExact(teamName)),
					},
				},
			},
		})
	})

	t.Run("already_exists", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "github_team" "test" {
  organization = "%s"
  name         = "test-team"
}
`, accTestConfigData.Values.Organization),
					ExpectError: regexp.MustCompile("Error: Failed to create team"),
				},
			},
		})
	})
}
