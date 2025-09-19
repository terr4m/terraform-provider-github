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

func TestAccOrganizationPropertiesDataSource(t *testing.T) {
	if accTestConfigData.AuthType == accAuthTypeUnauthenticated || !accTestConfigData.Features.Organization {
		t.Skip("Skipping test because the organization testing feature isn't enabled")
	}

	t.Run("properties", func(t *testing.T) {
		propertyName := fmt.Sprintf("%s%s", accTestConfigData.ResourcePrefix, acctest.RandomWithPrefix("test"))

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "github_organization_property" "test" {
  name         = "%s"
  organization = "%s"
  value_type   = "string"
}

data "github_organization_properties" "test" {
  organization = "%[2]s"

  depends_on = [
    github_organization_property.test
  ]
}
`, propertyName, accTestConfigData.Values.Organization),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.github_organization_properties.test", tfjsonpath.New("organization"), knownvalue.StringExact(accTestConfigData.Values.Organization)),
						statecheck.ExpectKnownValue("data.github_organization_properties.test", tfjsonpath.New("properties"), knownvalue.ListPartial(map[int]knownvalue.Check{0: knownvalue.NotNull()})),
					},
				},
			},
		})
	})
}
