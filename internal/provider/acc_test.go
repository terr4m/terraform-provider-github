package provider

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/google/go-github/v72/github"
)

type accTestConfig struct {
	Authenticated    bool
	Owner            string
	OwnerType        string
	AdvancedSecurity bool
	Username         string
	TeamSlug         string
	ResourcePrefix   string
}

var accTestConfigValues accTestConfig

func TestMain(m *testing.M) {
	accTestConfigValues = accTestConfig{
		Authenticated:    os.Getenv("ACC_GITHUB_AUTHENTICATED") == "true",
		Owner:            os.Getenv("ACC_GITHUB_OWNER"),
		OwnerType:        strings.ToUpper(os.Getenv("ACC_GITHUB_OWNER_TYPE")),
		AdvancedSecurity: os.Getenv("ACC_GITHUB_ADVANCED_SECURITY") == "true",
		Username:         "stevehipwelltesting",
		TeamSlug:         "test-team",
		ResourcePrefix:   fmt.Sprintf("test-acc-%s-", acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)),
	}

	if accTestConfigValues.Authenticated {
		if len(accTestConfigValues.Owner) == 0 {
			fmt.Println("no owner configured")
			os.Exit(1)
		}

		if len(accTestConfigValues.Owner) == 0 {
			fmt.Println("no owner configured")
			os.Exit(1)
		}

		if accTestConfigValues.OwnerType != "USER" && accTestConfigValues.OwnerType != "ORGANIZATION" && accTestConfigValues.OwnerType != "ENTERPRISE" {
			fmt.Println("invalid owner type")
			os.Exit(1)
		}
	}

	if accTestConfigValues.Authenticated {
		resource.TestMain(m)

		client, err := getGitHubClient(nil)
		if err != nil {
			fmt.Printf("error creating GitHub client: %v", err)
			os.Exit(1)
		}

		resource.AddTestSweepers("repos", &resource.Sweeper{
			Name: "repos",
			F: func(prefix string) error {
				repos, _, err := client.Repositories.ListByOrg(context.Background(), accTestConfigValues.Owner, &github.RepositoryListByOrgOptions{})
				if err != nil {
					return err
				}

				for _, r := range repos {
					if name := r.GetName(); strings.HasPrefix(name, prefix) {
						if _, err := client.Repositories.Delete(context.Background(), accTestConfigValues.Owner, name); err != nil {
							return err
						}
					}
				}

				return nil
			},
		})

		if accTestConfigValues.OwnerType != "USER" {
			resource.AddTestSweepers("teams", &resource.Sweeper{
				Name: "teams",
				F: func(prefix string) error {
					teams, _, err := client.Teams.ListTeams(context.Background(), accTestConfigValues.Owner, &github.ListOptions{})
					if err != nil {
						return err
					}

					for _, t := range teams {
						if slug := t.GetSlug(); strings.HasPrefix(slug, prefix) {
							if _, err := client.Teams.DeleteTeamBySlug(context.Background(), accTestConfigValues.Owner, slug); err != nil {
								return err
							}
						}
					}

					return nil
				},
			})
		}
	}

	m.Run()
}
