package ghutil

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	ghcht "github.com/bored-engineer/github-conditional-http-transport"
	bboltstorage "github.com/bored-engineer/github-conditional-http-transport/bbolt"
	"github.com/bradleyfalzon/ghinstallation/v2"
	ratelimit "github.com/gofri/go-github-ratelimit/v2/github_ratelimit"
	"github.com/google/go-github/v74/github"
)

// NewGitHubClient creates a new GitHub client with the given token and cache option.
func NewGitHubClient(token *string, cache bool) (*github.Client, error) {
	client, err := newGitHubClient(http.DefaultTransport, cache)
	if err != nil {
		return nil, fmt.Errorf("failed to create github client: %w", err)
	}

	if token == nil {
		return client, nil
	}

	return client.WithAuthToken(*token), nil
}

// NewGitHubClientForApp creates a new GitHub client for a GitHub App with the given credentials.
func NewGitHubClientForApp(appID int64, privateKey []byte, installationID int64, cache bool) (*github.Client, error) {
	tr := http.DefaultTransport

	atr, err := ghinstallation.NewAppsTransport(tr, appID, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create app transport: %w", err)
	}

	if installationID != -1 {
		tr = ghinstallation.NewFromAppsTransport(atr, installationID)
	} else {
		tr = atr
	}

	return newGitHubClient(tr, cache)
}

// newGitHubClient creates a new GitHub client with the given transport and cache option.
func newGitHubClient(tr http.RoundTripper, cache bool) (*github.Client, error) {
	tr = ratelimit.New(tr)

	if cache {
		ctr, err := cacheTransport(tr)
		if err != nil {
			return nil, fmt.Errorf("failed to create cache transport: %w", err)
		}
		tr = ctr
	}

	client := github.NewClient(&http.Client{Transport: tr})
	client.DisableRateLimitCheck = true

	return client, nil
}

// cacheTransport creates a new http.RoundTripper with a cache.
func cacheTransport(tr http.RoundTripper) (http.RoundTripper, error) {
	dir, err := os.MkdirTemp("", "terraform-provider-github")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	stor, err := bboltstorage.Open(filepath.Join(dir, "cache.db"), 0o644, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open cache: %w", err)
	}

	return ghcht.NewTransport(stor, tr), nil
}
