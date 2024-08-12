package permify

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
)

const (
	defaultPermifyImage        = "ghcr.io/permify/permify"
	defaultPermifyImageVersion = "v0.10.2"
	permifyRestPort            = "3476/tcp"
	permifyGrpcPort            = "3478/tcp"
	permifyStartupCommand      = "serve"
)

// PermifyContainer is a wrapper around testcontainers.Container
// that provides some conveince methods for working with Permify.
type PermifyContainer struct {
	testcontainers.Container
}

// Run creates an instance of the Permify container type.
func Run(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*PermifyContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", defaultPermifyImage, defaultPermifyImageVersion),
		ExposedPorts: []string{permifyRestPort, permifyGrpcPort},
		Cmd:          []string{permifyStartupCommand},
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	for _, opt := range opts {
		if err := opt.Customize(&genericContainerReq); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return nil, fmt.Errorf("failed to start permify: %w", err)
	}

	return &PermifyContainer{Container: container}, nil
}

// RESTPort returns the port which the Rest API for Permify is listening on.
func (p PermifyContainer) RESTPort(ctx context.Context) (int, error) {
	port, err := p.Container.MappedPort(ctx, permifyRestPort)
	if err != nil {
		return 0, err
	}
	return port.Int(), nil
}

// GRPCPort returns the port which the GRPC API for Permify is listening on.
func (p PermifyContainer) GRPCPort(ctx context.Context) (int, error) {
	port, err := p.Container.MappedPort(ctx, permifyGrpcPort)
	if err != nil {
		return 0, err
	}
	return port.Int(), nil
}
