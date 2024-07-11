package permify

import (
	"context"
	"errors"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
)

const (
	defaultPermifyImage        = "ghcr.io/permify/permify"
	defaultPermifyImageVersion = "v0.9.8"
	permifyRestPort            = "3476/tcp"
	permifyGrpcPort            = "3478/tcp"
	permifyStartupCommand      = "serve"
)

// PermifyContainer is a wrapper around testcontainers.Container
// that provides some conveince methods for working with Permify.
type PermifyContainer struct {
	testcontainers.Container
}

func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*PermifyContainer, error) {
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
			return nil, fmt.Errorf("error encountered while customizing: %w", err)
		}
	}

	if genericContainerReq.WaitingFor == nil {
	}
	return nil, errors.New("not implemented")
}
