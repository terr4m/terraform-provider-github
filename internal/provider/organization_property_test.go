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

func TestAccOrganizationPropertyResource(t *testing.T) {
	if accTestConfigData.AuthType == accAuthTypeUnauthenticated || !accTestConfigData.Features.Organization {
		t.Skip("Skipping test because the organization testing feature isn't enabled")
	}

	t.Run("create_string_property", func(t *testing.T) {
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
`, propertyName, accTestConfigData.Values.Organization),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("allowed_values"), knownvalue.Null()),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("default_value"), knownvalue.Null()),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("description"), knownvalue.Null()),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("editable_by"), knownvalue.StringExact("org_actors")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("name"), knownvalue.StringExact(propertyName)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("organization"), knownvalue.StringExact(accTestConfigData.Values.Organization)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("required"), knownvalue.Bool(false)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("source_type"), knownvalue.StringExact("organization")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("value_type"), knownvalue.StringExact("string")),
					},
				},
			},
		})
	})

	t.Run("create_string_property_full", func(t *testing.T) {
		propertyName := fmt.Sprintf("%s%s", accTestConfigData.ResourcePrefix, acctest.RandomWithPrefix("test"))

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "github_organization_property" "test" {
  default_value = "default-value"
  description   = "Test description."
  editable_by   = "org_and_repo_actors"
  name          = "%s"
  organization  = "%s"
  required      = true
  value_type    = "string"
}
`, propertyName, accTestConfigData.Values.Organization),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("allowed_values"), knownvalue.Null()),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("default_value"), knownvalue.StringExact("default-value")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("description"), knownvalue.StringExact("Test description.")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("editable_by"), knownvalue.StringExact("org_and_repo_actors")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("name"), knownvalue.StringExact(propertyName)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("organization"), knownvalue.StringExact(accTestConfigData.Values.Organization)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("required"), knownvalue.Bool(true)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("source_type"), knownvalue.StringExact("organization")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("value_type"), knownvalue.StringExact("string")),
					},
				},
			},
		})
	})

	t.Run("create_single_select_property", func(t *testing.T) {
		propertyName := fmt.Sprintf("%s%s", accTestConfigData.ResourcePrefix, acctest.RandomWithPrefix("test"))

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "github_organization_property" "test" {
  allowed_values = ["option1", "option2", "option3"]
  name           = "%s"
  organization   = "%s"
  value_type     = "single_select"
}
`, propertyName, accTestConfigData.Values.Organization),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("allowed_values"), knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("option1"),
							knownvalue.StringExact("option2"),
							knownvalue.StringExact("option3"),
						})),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("default_value"), knownvalue.Null()),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("description"), knownvalue.Null()),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("editable_by"), knownvalue.StringExact("org_actors")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("name"), knownvalue.StringExact(propertyName)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("organization"), knownvalue.StringExact(accTestConfigData.Values.Organization)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("required"), knownvalue.Bool(false)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("source_type"), knownvalue.StringExact("organization")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("value_type"), knownvalue.StringExact("single_select")),
					},
				},
			},
		})
	})

	t.Run("create_single_select_property_full", func(t *testing.T) {
		propertyName := fmt.Sprintf("%s%s", accTestConfigData.ResourcePrefix, acctest.RandomWithPrefix("test"))

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "github_organization_property" "test" {
  allowed_values = ["option1", "option2", "option3"]
  default_value  = "option1"
  description    = "Test description."
  editable_by    = "org_and_repo_actors"
  name           = "%s"
  organization   = "%s"
  required       = true
  value_type     = "single_select"
}
`, propertyName, accTestConfigData.Values.Organization),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("allowed_values"), knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("option1"),
							knownvalue.StringExact("option2"),
							knownvalue.StringExact("option3"),
						})),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("default_value"), knownvalue.StringExact("option1")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("description"), knownvalue.StringExact("Test description.")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("editable_by"), knownvalue.StringExact("org_and_repo_actors")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("name"), knownvalue.StringExact(propertyName)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("organization"), knownvalue.StringExact(accTestConfigData.Values.Organization)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("required"), knownvalue.Bool(true)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("source_type"), knownvalue.StringExact("organization")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("value_type"), knownvalue.StringExact("single_select")),
					},
				},
			},
		})
	})

	t.Run("create_multi_select_property", func(t *testing.T) {
		propertyName := fmt.Sprintf("%s%s", accTestConfigData.ResourcePrefix, acctest.RandomWithPrefix("test"))

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "github_organization_property" "test" {
  allowed_values = ["tag1", "tag2", "tag3", "tag4"]
  name           = "%s"
  organization   = "%s"
  value_type     = "multi_select"
}
`, propertyName, accTestConfigData.Values.Organization),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("allowed_values"), knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("tag1"),
							knownvalue.StringExact("tag2"),
							knownvalue.StringExact("tag3"),
							knownvalue.StringExact("tag4"),
						})),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("default_value"), knownvalue.Null()),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("description"), knownvalue.Null()),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("editable_by"), knownvalue.StringExact("org_actors")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("name"), knownvalue.StringExact(propertyName)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("organization"), knownvalue.StringExact(accTestConfigData.Values.Organization)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("required"), knownvalue.Bool(false)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("source_type"), knownvalue.StringExact("organization")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("value_type"), knownvalue.StringExact("multi_select")),
					},
				},
			},
		})
	})

	t.Run("create_multi_select_property_full", func(t *testing.T) {
		propertyName := fmt.Sprintf("%s%s", accTestConfigData.ResourcePrefix, acctest.RandomWithPrefix("test"))

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "github_organization_property" "test" {
  allowed_values = ["tag1", "tag2", "tag3", "tag4"]
  # default_value  = "tag1"
  description    = "Test description."
  editable_by    = "org_and_repo_actors"
  name           = "%s"
  organization   = "%s"
  # required       = true
  value_type     = "multi_select"
}
`, propertyName, accTestConfigData.Values.Organization),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("allowed_values"), knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("tag1"),
							knownvalue.StringExact("tag2"),
							knownvalue.StringExact("tag3"),
							knownvalue.StringExact("tag4"),
						})),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("default_value"), knownvalue.Null()),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("description"), knownvalue.StringExact("Test description.")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("editable_by"), knownvalue.StringExact("org_and_repo_actors")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("name"), knownvalue.StringExact(propertyName)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("organization"), knownvalue.StringExact(accTestConfigData.Values.Organization)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("required"), knownvalue.Bool(false)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("source_type"), knownvalue.StringExact("organization")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("value_type"), knownvalue.StringExact("multi_select")),
					},
				},
			},
		})
	})

	t.Run("create_true_false_property", func(t *testing.T) {
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
  value_type   = "true_false"
}
`, propertyName, accTestConfigData.Values.Organization),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("allowed_values"), knownvalue.Null()),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("default_value"), knownvalue.Null()),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("description"), knownvalue.Null()),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("editable_by"), knownvalue.StringExact("org_actors")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("name"), knownvalue.StringExact(propertyName)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("organization"), knownvalue.StringExact(accTestConfigData.Values.Organization)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("required"), knownvalue.Bool(false)),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("source_type"), knownvalue.StringExact("organization")),
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("value_type"), knownvalue.StringExact("true_false")),
					},
				},
			},
		})
	})

	t.Run("update_allowed_values", func(t *testing.T) {
		propertyName := fmt.Sprintf("%s%s", accTestConfigData.ResourcePrefix, acctest.RandomWithPrefix("test"))

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "github_organization_property" "test" {
  allowed_values = ["option1", "option2"]
  name           = "%s"
  organization   = "%s"
  value_type     = "single_select"
}
`, propertyName, accTestConfigData.Values.Organization),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("allowed_values"), knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("option1"),
							knownvalue.StringExact("option2"),
						})),
					},
				},
				{
					Config: fmt.Sprintf(`
resource "github_organization_property" "test" {
  allowed_values = ["option1", "option2", "option3", "option4"]
  name           = "%s"
  organization   = "%s"
  value_type     = "single_select"
}
`, propertyName, accTestConfigData.Values.Organization),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("allowed_values"), knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("option1"),
							knownvalue.StringExact("option2"),
							knownvalue.StringExact("option3"),
							knownvalue.StringExact("option4"),
						})),
					},
				},
			},
		})
	})

	t.Run("update_default_value", func(t *testing.T) {
		propertyName := fmt.Sprintf("%s%s", accTestConfigData.ResourcePrefix, acctest.RandomWithPrefix("test"))
		originalDefault := "foo"
		updatedDefault := "bar"

		config := `
resource "github_organization_property" "test" {
  default_value = "%s"
  name          = "%s"
  organization  = "%s"
	required      = true
  value_type    = "string"
}
`

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(config, originalDefault, propertyName, accTestConfigData.Values.Organization),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("default_value"), knownvalue.StringExact(originalDefault)),
					},
				},
				{
					Config: fmt.Sprintf(config, updatedDefault, propertyName, accTestConfigData.Values.Organization),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("default_value"), knownvalue.StringExact(updatedDefault)),
					},
				},
			},
		})
	})

	t.Run("update_description", func(t *testing.T) {
		propertyName := fmt.Sprintf("%s%s", accTestConfigData.ResourcePrefix, acctest.RandomWithPrefix("test"))
		originalDescription := "Original description"
		updatedDescription := "Updated description"

		config := `
resource "github_organization_property" "test" {
  description  = "%s"
  name         = "%s"
  organization = "%s"
  value_type   = "string"
}
`

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(config, originalDescription, propertyName, accTestConfigData.Values.Organization),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("description"), knownvalue.StringExact(originalDescription)),
					},
				},
				{
					Config: fmt.Sprintf(config, updatedDescription, propertyName, accTestConfigData.Values.Organization),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("github_organization_property.test", tfjsonpath.New("description"), knownvalue.StringExact(updatedDescription)),
					},
				},
			},
		})
	})

	t.Run("invalid_required_no_default", func(t *testing.T) {
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
  required     = true
  value_type   = "true_false"
}
`, propertyName, accTestConfigData.Values.Organization),
					ExpectError: regexp.MustCompile(`Default value must be present`),
				},
			},
		})
	})

	t.Run("invalid_value_type", func(t *testing.T) {
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
  value_type   = "invalid_type"
}
`, propertyName, accTestConfigData.Values.Organization),
					ExpectError: regexp.MustCompile(`Attribute value_type value must be one of`),
				},
			},
		})
	})

	t.Run("invalid_editable_by", func(t *testing.T) {
		propertyName := fmt.Sprintf("%s%s", accTestConfigData.ResourcePrefix, acctest.RandomWithPrefix("test"))

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "github_organization_property" "test" {
  editable_by  = "invalid_option"
  name         = "%s"
  organization = "%s"
  value_type   = "string"
}
`, propertyName, accTestConfigData.Values.Organization),
					ExpectError: regexp.MustCompile(`Attribute editable_by value must be one of`),
				},
			},
		})
	})
}
