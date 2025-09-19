package ghutil

import (
	"context"
	"fmt"

	"github.com/google/go-github/v74/github"
)

// ClientCreator provides an interface to dynamically return GitHub clients with the correct authentication.
type ClientCreator interface {
	AppClient() (*github.Client, error)
	DefaultClient(ctx context.Context) (*github.Client, error)
	OrganizationClient(ctx context.Context, organization string) (*github.Client, error)
}

// clientCreator is responsible for creating GitHub clients using optional token authentication.
type clientCreator struct {
	token         *string
	cacheRequests bool
	client        *github.Client
}

// NewClientCreator creates a ClientCreator that can optionally authenticate using a token.
func NewClientCreator(token *string, cacheRequests bool) (ClientCreator, error) {
	cc := &clientCreator{
		token:         token,
		cacheRequests: cacheRequests,
	}

	return cc, nil
}

// AppClient returns a GitHub app client.
func (cc *clientCreator) AppClient() (*github.Client, error) {
	return nil, fmt.Errorf("not an app client")
}

// DefaultClient returns the default GitHub client.
func (cc *clientCreator) DefaultClient(ctx context.Context) (*github.Client, error) {
	if cc.client != nil {
		return cc.client, nil
	}

	c, err := NewGitHubClient(cc.token, cc.cacheRequests)
	if err != nil {
		return nil, fmt.Errorf("failed to create github client: %w", err)
	}
	cc.client = c

	return c, nil
}

// OrganizationClient returns a GitHub client for an organization.
func (cc *clientCreator) OrganizationClient(ctx context.Context, organization string) (*github.Client, error) {
	return cc.DefaultClient(ctx)
}
