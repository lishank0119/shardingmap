[![Go Reference](https://pkg.go.dev/badge/github.com/lishank0119/shardingmap.svg)](https://pkg.go.dev/github.com/lishank0119/shardingmap)
[![go.mod](https://img.shields.io/github/go-mod/go-version/lishank0119/shardingmap)](go.mod)

# ShardingMap

[中文](README.zh-TW.md)

A simple and efficient sharding map (distributed hash map) implementation in Go. It allows splitting your key-value
pairs across multiple shards to achieve thread-safe, high-performance access to large datasets.

## Features

- **Sharded**: Distribute your data across multiple shards for better scalability and concurrency.
- **Thread-Safe**: Use of `sync.RWMutex` ensures concurrent reads and writes are handled safely.
- **Customizable**: Set custom shard count and sharding function.
- **FNV Hashing**: Default sharding function based on FNV-1a hash.

## Installation

```bash
go get github.com/lishank0119/shardingmap
```

## Usage

### Creating a Sharded Map

```go
package main

import (
	"fmt"
	"github.com/lishank0119/shardingmap"
)

func main() {
	// Create a new sharded map with custom options
	sm := shardingmap.New[int, string](
		shardingmap.WithShardCount[int, string](32),                                      // Set the number of shards
		shardingmap.WithShardingFunc[int, string](func(key int) int { return key }), // Custom sharding function
	)

	// Set some key-value pairs
	sm.Set(1, "value1")
	sm.Set(2, "value2")

	// Retrieve values
	if value, ok := sm.Get(1); ok {
		fmt.Println(value) // Output: value1
	}

	// Delete a key
	sm.Delete(2)

	// Iterate over all key-value pairs
	sm.ForEach(func(k int, v string) {
		fmt.Printf("%d: %s\n", k, v)
	})

	fmt.Println("Total items:", sm.Len()) // Output: Total items: 1
}
```

### Methods

- **New**: Creates a new sharded map with configurable options (`WithShardCount` and `WithShardingFunc`).
- **Set**: Adds or updates a key-value pair in the map.
- **Get**: Retrieves the value associated with the key.
- **Delete**: Removes the key-value pair from the map.
- **Len**: Returns the total number of key-value pairs in the map.
- **ForEach**: Iterates over all key-value pairs.

## Customization

You can customize the following:

- **Shard Count**: Define the number of shards (default is 16).
- **Sharding Function**: Provide a custom sharding function for more complex data distributions.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
