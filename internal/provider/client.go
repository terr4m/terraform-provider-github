package provider

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	ghcht "github.com/bored-engineer/github-conditional-http-transport"
	bboltstorage "github.com/bored-engineer/github-conditional-http-transport/bbolt"
	"github.com/bradleyfalzon/ghinstallation/v2"
	ratelimit "github.com/gofri/go-github-ratelimit/v2/github_ratelimit"
	"github.com/google/go-github/v72/github"
)

// getGitHubClient creates a new GitHub client.
func getGitHubClient(model *GitHubProviderModel) (*github.Client, error) {
	tr := http.DefaultTransport

	if model != nil && model.AppAuth != nil {
		atr, err := getAppTransport(tr, model.AppAuth)
		if err != nil {
			return nil, fmt.Errorf("failed to create app transport: %w", err)
		}
		tr = atr
	}

	tr = ratelimit.New(tr)

	if model != nil && model.CacheRequests.ValueBool() {
		ctr, err := getCacheTransport(tr)
		if err != nil {
			return nil, fmt.Errorf("failed to create cache transport: %w", err)
		}
		tr = ctr
	}

	client := github.NewClient(&http.Client{Transport: tr})

	if model != nil && model.AppAuth == nil {
		if !model.Token.IsNull() {
			return client.WithAuthToken(model.Token.String()), nil
		}

		if v := os.Getenv("GITHUB_TOKEN"); len(v) != 0 {
			return client.WithAuthToken(v), nil
		}
	}

	return client, nil
}

// getAppTransport creates a new http.RoundTripper for the given app auth.
func getAppTransport(tr http.RoundTripper, appAuth *AppAuthModel) (http.RoundTripper, error) {
	var key []byte

	appID := appAuth.ID.ValueInt64()
	installationID := appAuth.InstallationID.ValueInt64()

	if !appAuth.PrivateKey.IsNull() {
		key = []byte(appAuth.PrivateKey.String())
	} else if !appAuth.PrivateKeyFile.IsNull() {
		k, err := os.ReadFile(appAuth.PrivateKeyFile.String())
		if err != nil {
			return nil, fmt.Errorf("failed to read private key file: %w", err)
		}
		key = k
	} else if v := os.Getenv("GITHUB_APP_PRIVATE_KEY"); len(v) != 0 {
		key = []byte(v)
	} else if v := os.Getenv("GITHUB_APP_PRIVATE_KEY_FILE"); len(v) != 0 {
		k, err := os.ReadFile(v)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key file: %w", err)
		}
		key = k
	} else {
		return nil, fmt.Errorf("no private key provided")
	}

	atr, err := ghinstallation.NewAppsTransport(tr, appID, key)
	if err != nil {
		return nil, fmt.Errorf("failed to create apps transport: %w", err)
	}

	return ghinstallation.NewFromAppsTransport(atr, installationID), nil
}

// getCacheTransport creates a new http.RoundTripper with a cache.
func getCacheTransport(tr http.RoundTripper) (http.RoundTripper, error) {
	dir, err := os.MkdirTemp("", "terraform-provider-github")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	stor, err := bboltstorage.Open(filepath.Join(dir, "cache.db"), 0644, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open cache: %w", err)
	}

	return ghcht.NewTransport(stor, tr), nil
}
