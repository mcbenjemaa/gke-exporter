package gke

import (
	"context"

	"cloud.google.com/go/container"
)

type GKEMetadataFetcher interface {
	FetchMetadata(ctx context.Context) ([]*container.Resource, error)
}

type GKEClient struct {
	Client    *container.Client
	Zone      string
	ProjectID string
}

// NewGKEClient returns gke client
// uses Application Default Credentials to authenticate.
func NewGKEClient(ctx context.Context, projectID string) (*container.Client, error) {
	return container.NewClient(ctx, projectID)
}

// FetchMetadata returns all cluster information
func (gkeClient GKEClient) FetchMetadata(ctx context.Context) ([]*container.Resource, error) {
	return gkeClient.Client.Clusters(ctx, gkeClient.Zone)
}
