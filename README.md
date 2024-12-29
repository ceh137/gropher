# gropher

gropher is a high-performance, thread-safe graph library for Go that provides flexible data storage and serialization capabilities. Similar to Python's NetworkX, it allows storing arbitrary data in Nodes while maintaining type safety and performance.

## Features

- 🚀 High-performance graph operations
- 🔒 Thread-safe implementation
- 💾 JSON serialization support
- 🎯 Generic data storage in Nodes
- 📦 Easy to use API
- 🧪 Comprehensive test coverage
- 📝 Well-documented code

## Installation

```bash
go get github.com/yourusername/gropher
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/yourusername/gropher"
)

func main() {
    // Create a new graph
    g := gropher.New()

    // Add Nodes with any type of data
    g.AddNode("1", "String data")
    g.AddNode("2", 42)
    g.AddNode("3", map[string]interface{}{
        "name": "Complex Data",
        "value": 100,
    })

    // Add Edges
    g.AddEdge("1", "2", 1.0)
    g.AddEdge("2", "3", 2.5)

    // Get node data
    node, _ := g.GetNode("1")
    fmt.Println("Node 1 data:", node.Data)

    // Get neighbors
    neighbors, _ := g.GetNeighbors("1")
    fmt.Println("Node 1 neighbors:", neighbors)

    // Save to file
    g.SaveToFile("my_graph.json")

    // Load from file
    newGraph := gropher.New()
    newGraph.LoadFromFile("my_graph.json")
}
```

## API Documentation

### Graph Creation

```go
g := gropher.New()
```

### Node Operations

```go
// Add a node
err := g.AddNode(id string, data interface{})

// Get a node
node, err := g.GetNode(id string)

// Remove a node
err := g.RemoveNode(id string)
```

### Edge Operations

```go
// Add an edge
err := g.AddEdge(from string, to string, weight float64)

// Remove an edge
err := g.RemoveEdge(from string, to string)

// Get neighbors
neighbors, err := g.GetNeighbors(id string)
```

### Serialization

```go
// Save to file
err := g.SaveToFile(filename string)

// Load from file
err := g.LoadFromFile(filename string)
```

## Custom Data Types

gropher supports any JSON-serializable data type in Nodes:

```go
type PersonData struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Country string `json:"country"`
}

// Add node with custom data
g.AddNode("person1", PersonData{
    Name: "Alice",
    Age: 30,
    Country: "USA",
})
```

## Thread Safety

All operations in gropher are thread-safe by default. The library uses mutex locks to ensure safe concurrent access:

```go
// Safe for concurrent use
go func() {
    g.AddNode("1", "data")
}()
go func() {
    g.GetNode("1")
}()
```

## Contributing

We welcome contributions! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Development Setup

1. Clone the repository
2. Install dependencies (if any)
3. Run tests: `go test ./...`

### Contribution Guidelines

- Write tests for new features
- Follow Go best practices and coding conventions
- Update documentation for significant changes
- Add comments for complex logic
- Consider performance implications

## Future Improvements

- [ ] Add graph traversal algorithms (BFS, DFS)
- [ ] Implement path finding algorithms
- [ ] Add support for undirected graphs
- [ ] Create visualization tools
- [ ] Add graph metrics and analysis
- [ ] Implement graph partitioning
- [ ] Add support for weighted paths
- [ ] Create graph database integration

## Performance Considerations

- Node operations: O(1) average case
- Edge operations: O(1) average case
- Memory usage: O(|V| + |E|) where V is vertices and E is Edges
- Thread-safe operations may have slight overhead due to mutex locks

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by Python's NetworkX library
- Thanks to all contributors and users

## Support

- Create an issue for bug reports or feature requests
- Star the repository if you find it useful
- Contribute code or documentation improvements

## Code of Conduct

Please note that this project is released with a Contributor Code of Conduct. By participating in this project you agree to abide by its terms.