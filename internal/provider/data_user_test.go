package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccUserDataSource(t *testing.T) {
	t.Run("user_exists_by_login", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
data "github_user" "test" {
  login = "github"
}
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.github_user.test", tfjsonpath.New("name"), knownvalue.StringExact("GitHub")),
					},
				},
			},
		})
	})

	t.Run("user_exists_by_id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
data "github_user" "test" {
  id = 9919
}
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.github_user.test", tfjsonpath.New("name"), knownvalue.StringExact("GitHub")),
					},
				},
			},
		})
	})

	t.Run("user_does_not_exist", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
data "github_user" "test" {
  login = "should-not-exist"
}
`,
					ExpectError: regexp.MustCompile("Error: Failed to get user"),
				},
			},
		})
	})
}
