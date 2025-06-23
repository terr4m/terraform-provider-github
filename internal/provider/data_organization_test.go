package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccOrganizationDataSource(t *testing.T) {
	t.Run("organization_exists", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
data "github_organization" "test" {
  login = "github"
}
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.github_organization.test", tfjsonpath.New("name"), knownvalue.StringExact("GitHub")),
					},
				},
			},
		})
	})

	t.Run("organization_does_not_exist", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
data "github_organization" "test" {
  login = "should-not-exist"
}
`,
					ExpectError: regexp.MustCompile("Error: Failed to get organization"),
				},
			},
		})
	})
}
