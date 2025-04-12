[![Go Reference](https://pkg.go.dev/badge/github.com/ruizu/go-pgxshard.svg)](https://pkg.go.dev/github.com/ruizu/go-pgxshard)

# go-pgxshard - PostgreSQL sharding for Golang

go-pgxshard is a Golang module that provides a simple and efficient way to manage and interact with multiple PostgreSQL database shards using the [`pgxpool`](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool) library. It allows developers to distribute data across multiple shards and provides utility methods for shard management, connectivity, and querying.

## Features

- **Shard Management**: Easily manage multiple PostgreSQL shards.
- **Custom Shard Indexing**: Define custom shard indexing logic based on your application's requirements.
- **Connection Pooling**: Leverages `pgxpool` for efficient connection pooling.

## Installation

To install the module, use:

```bash
go get github.com/ruizu/go-pgxshard
```

## Usage

### Creating a ShardManager

```go
package main

import (
	"context"
	"log"
	"github.com/ruizu/go-pgxshard"
)

func main() {
	ctx := context.Background()
	connectionStrings := []string{
		"postgres://user:password@localhost:5432/shard1",
		"postgres://user:password@localhost:5432/shard2",
	}

	shardManager, err := pgxshard.New(ctx, connectionStrings)
	if err != nil {
		log.Fatalf("Failed to create ShardManager: %v", err)
	}
	defer shardManager.Close(ctx)

	// Use shardManager to interact with shards
}
```

### Setting a Custom Shard Index Function

```go
shardManager.SetShardIndexFunc(ctx, func(key any, count int) (int, error) {
	// Custom logic to determine shard index
})
```

### Accessing a Shard

```go
shard, err := shardManager.Shard(ctx, "some-key")
if err != nil {
	log.Fatalf("Failed to get shard: %v", err)
}

// Use the shard (pgxpool.Pool) for database operations
```

### Checking Connectivity

```go
if err := shardManager.Ping(ctx); err != nil {
	log.Fatalf("Shard connectivity issue: %v", err)
}
```

## License

This project is licensed under the MIT License. See the [LICENSE.md](LICENSE.md) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## Acknowledgments

- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit for Go.