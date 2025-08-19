package provider

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

type accAuthType string

const (
	accAuthTypeUnauthenticated     accAuthType = "UNAUTHENTICATED"
	accAuthTypePersonalAccessToken accAuthType = "PERSONAL_ACCESS_TOKEN"
	accAuthTypeGitHubApp           accAuthType = "GITHUB_APP"
)

type accTestConfig struct {
	ResourcePrefix string
	AuthType       accAuthType
	Features       accTestFeatures
	Values         accTestValues
}

type accTestFeatures struct {
	Organization     bool
	Enterprise       bool
	AdvancedSecurity bool
}

type accTestValues struct {
	Username     string
	Organization string
	TeamSlug     string
}

var accTestConfigData accTestConfig

func TestMain(m *testing.M) {
	accTestConfigData = accTestConfig{
		ResourcePrefix: fmt.Sprintf("test-acc-%s-", acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)),
		AuthType:       accAuthType(strings.ToUpper(os.Getenv("ACC_GITHUB_AUTH_TYPE"))),
		Features: accTestFeatures{
			Organization:     os.Getenv("ACC_GITHUB_FEATURE_ORGANIZATION") == "true",
			Enterprise:       os.Getenv("ACC_GITHUB_FEATURE_ENTERPRISE") == "true",
			AdvancedSecurity: os.Getenv("ACC_GITHUB_FEATURE_ADVANCED_SECURITY") == "true",
		},
		Values: accTestValues{
			Username:     os.Getenv("ACC_GITHUB_VALUE_USERNAME"),
			Organization: os.Getenv("ACC_GITHUB_VALUE_ORGANIZATION"),
			TeamSlug:     os.Getenv("ACC_GITHUB_VALUE_TEAM"),
		},
	}

	// if accTestConfigValues.Authenticated {
	// 	resource.TestMain(m)

	// 	client, err := getGitHubClient(nil)
	// 	if err != nil {
	// 		fmt.Printf("error creating GitHub client: %v", err)
	// 		os.Exit(1)
	// 	}

	// 	resource.AddTestSweepers("repos", &resource.Sweeper{
	// 		Name: "repos",
	// 		F: func(prefix string) error {
	// 			repos, _, err := client.Repositories.ListByOrg(context.Background(), accTestConfigValues.Owner, &github.RepositoryListByOrgOptions{})
	// 			if err != nil {
	// 				return err
	// 			}

	// 			for _, r := range repos {
	// 				if name := r.GetName(); strings.HasPrefix(name, prefix) {
	// 					if _, err := client.Repositories.Delete(context.Background(), accTestConfigValues.Owner, name); err != nil {
	// 						return err
	// 					}
	// 				}
	// 			}

	// 			return nil
	// 		},
	// 	})

	// 	if accTestConfigValues.OwnerType != "USER" {
	// 		resource.AddTestSweepers("teams", &resource.Sweeper{
	// 			Name: "teams",
	// 			F: func(prefix string) error {
	// 				teams, _, err := client.Teams.ListTeams(context.Background(), accTestConfigValues.Owner, &github.ListOptions{})
	// 				if err != nil {
	// 					return err
	// 				}

	// 				for _, t := range teams {
	// 					if slug := t.GetSlug(); strings.HasPrefix(slug, prefix) {
	// 						if _, err := client.Teams.DeleteTeamBySlug(context.Background(), accTestConfigValues.Owner, slug); err != nil {
	// 							return err
	// 						}
	// 					}
	// 				}

	// 				return nil
	// 			},
	// 		})
	// 	}
	// }

	m.Run()
}
