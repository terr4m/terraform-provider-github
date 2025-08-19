package ghutil

import (
	"context"
	"fmt"

	"github.com/google/go-github/v74/github"
	lru "github.com/hashicorp/golang-lru/v2"
)

// appClientCreator is responsible for creating GitHub clients using app authentication.
type appClientCreator struct {
	appID         int64
	privateKey    []byte
	cacheRequests bool
	clients       *lru.Cache[string, *github.Client]
}

// NewAppClientCreator creates a ClientCreator than can authenticate using a GitHub app.
func NewAppClientCreator(appID int64, privateKey []byte, capacity int, cacheRequests bool) (ClientCreator, error) {
	cache, err := lru.New[string, *github.Client](capacity)
	if err != nil {
		return nil, fmt.Errorf("failed to create client cache: %w", err)
	}

	cc := &appClientCreator{
		appID:         appID,
		privateKey:    privateKey,
		cacheRequests: cacheRequests,
		clients:       cache,
	}

	return cc, nil
}

// AppClient returns a GitHub app client.
func (cc *appClientCreator) AppClient() (*github.Client, error) {
	key := "-"
	c, ok := cc.clients.Get(key)
	if ok {
		return c, nil
	}

	c, err := NewGitHubClientForApp(cc.appID, cc.privateKey, -1, cc.cacheRequests)
	if err != nil {
		return nil, fmt.Errorf("failed to create github client: %w", err)
	}
	cc.clients.Add(key, c)

	return c, nil
}

// DefaultClient returns the default GitHub client.
func (cc *appClientCreator) DefaultClient(ctx context.Context) (*github.Client, error) {
	key := "_"
	c, ok := cc.clients.Get(key)
	if ok {
		return c, nil
	}

	ac, err := cc.AppClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get app client: %w", err)
	}

	insts, _, err := ac.Apps.ListInstallations(ctx, &github.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list installations: %w", err)
	}

	if len(insts) == 0 {
		return nil, fmt.Errorf("no app installations found")
	}

	inst := insts[0]
	c, err = NewGitHubClientForApp(cc.appID, cc.privateKey, inst.GetID(), cc.cacheRequests)
	if err != nil {
		return nil, fmt.Errorf("failed to create installation client: %w", err)
	}
	cc.clients.Add(key, c)
	cc.clients.Add(inst.GetAccount().GetLogin(), c)

	return c, nil
}

// OrganizationClient returns a GitHub client for an organization.
func (cc *appClientCreator) OrganizationClient(ctx context.Context, organization string) (*github.Client, error) {
	c, ok := cc.clients.Get(organization)
	if ok {
		return c, nil
	}

	ac, err := cc.AppClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get app client: %w", err)
	}

	inst, _, err := ac.Apps.FindOrganizationInstallation(ctx, organization)
	if err != nil {
		return nil, fmt.Errorf("failed to get installation ID: %w", err)
	}

	c, err = NewGitHubClientForApp(cc.appID, cc.privateKey, inst.GetID(), cc.cacheRequests)
	if err != nil {
		return nil, fmt.Errorf("failed to create installation client: %w", err)
	}
	cc.clients.Add(organization, c)

	return c, nil
}
