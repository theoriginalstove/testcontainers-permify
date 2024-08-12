<!-- markdownlint-configure-file {
  "MD033": false,
  "MD041": false
} -->
<div align="center">

# Permify Testcontainer - [testcontainers](https://www.testcontainers.org/) implementation for [Permify](https://permify.co)

</div>

# testcontainers-aerospike-go


## Install

Use `go get` to install the latest version of the library.

```bash
go get -u github.com/theoriginalstove/testcontainers-permify@latest
```

## Usage

```go
import (
    "context"
    "testing"

    "github.com/stretchr/testify/require"
	permifygrpc "github.com/Permify/permify-go/grpc"
    permifytest "github.com/theoriginalstove/testcontainers-permify"
)

func TestPermify(t *testing.T) {
    permifyClient := setupPermify(t)
    // your code here
}

func setupPermify(t *testing.T) *permify_gcpc.Client {
    ctx := context.Background()

    container, err := permifytest.RunContainer(ctx)
    require.NoError(t, err)
    t.Cleanup(func() {
        err := container.Terminate(ctx)
        require.NoError(t, err)
    })

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

    return client
}
```
