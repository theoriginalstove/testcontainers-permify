package permify

import (
	"context"
	"fmt"
	"testing"

	permify_payload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
	permify_grpc "github.com/Permify/permify-go/grpc"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestRunContainer(t *testing.T) {
	ctx := context.Background()

	container, err := RunContainer(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		err = container.Terminate(ctx)
		require.NoError(t, err, "failed to terminate container")
	})

	host, err := container.Host(ctx)
	require.NoError(t, err, "failed to fetch permify host")

	grpcPort, err := container.GRPCPort(ctx)
	require.NoError(t, err, "failed to fetch permify grpc api port")

	client, err := permify_grpc.NewClient(permify_grpc.Config{
		Endpoint: fmt.Sprintf("%s:%d", host, grpcPort),
	}, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "failed to initialize permify grpc clinet")

	// create a tenant
	tenant, err := client.Tenancy.Create(ctx, &permify_payload.TenantCreateRequest{
		Id:   "t1",
		Name: "test tenant 1",
	})
	require.NoError(t, err, "failed to create a test tenant")
	require.Equal(t, "t1", tenant.Tenant.Id)
	require.Equal(t, "test tenant 1", tenant.Tenant.Name)

	// write a schema
	schema, err := client.Schema.Write(ctx, &permify_payload.SchemaWriteRequest{
		TenantId: "t1",
		Schema: `
		entity user {}

		entity document {
			relation viewer @user
			action view = viewer
		}`,
	})
	require.NoError(t, err, "failed to write a schema to permify service")
	require.NotEmpty(t, schema.SchemaVersion)

	// write a relationship
	relationship, err := client.Data.WriteRelationships(ctx, &permify_payload.RelationshipWriteRequest{
		TenantId: "t1",
		Metadata: &permify_payload.RelationshipWriteRequestMetadata{
			SchemaVersion: schema.SchemaVersion,
		},
		Tuples: []*permify_payload.Tuple{
			{
				Entity: &permify_payload.Entity{
					Type: "document",
					Id:   "1",
				},
				Relation: "viewer",
				Subject: &permify_payload.Subject{
					Type: "user",
					Id:   "1",
				},
			},
			{
				Entity: &permify_payload.Entity{
					Type: "document",
					Id:   "3",
				},
				Relation: "viewer",
				Subject: &permify_payload.Subject{
					Type: "user",
					Id:   "1",
				},
			},
		},
	})
	require.NoError(t, err, "failed to write relationships")
	require.NotEmpty(t, relationship.SnapToken, "snap token was empty")
}
