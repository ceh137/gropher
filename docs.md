# Gropher API Documentation

## Table of Contents

- [Types](#types)
- [Functions](#functions)
- [Methods](#methods)
- [Error Handling](#error-handling)
- [Examples](#examples)

## Types

### Graph

```go
type Graph struct {
    // Contains filtered or unexported fields
}
```

The main graph structure. Thread-safe and supports concurrent operations.

### Node

```go
type Node struct {
    ID   string      `json:"id"`
    Data interface{} `json:"data"`
}
```

Represents a vertex in the graph. Can store any JSON-serializable data.

### Edge

```go
type Edge struct {
    From   string  `json:"from"`
    To     string  `json:"to"`
    Weight float64 `json:"weight"`
}
```

Represents a directed weighted edge between two Nodes.

## Functions

### New

```go
func New() *Graph
```

Creates a new empty graph. This is the main entry point for using the library.

**Example:**
```go
g := gropher.New()
```

## Methods

### AddNode

```go
func (g *Graph) AddNode(id string, data interface{}) error
```

Adds a new node to the graph with the specified ID and data.

**Parameters:**
- `id`: Unique identifier for the node
- `data`: Any JSON-serializable data to store in the node

**Returns:**
- `error`: If node already exists or invalid parameters

**Example:**
```go
err := g.AddNode("1", map[string]interface{}{
    "name": "Node 1",
    "value": 42,
})
```

### RemoveNode

```go
func (g *Graph) RemoveNode(id string) error
```

Removes a node and all its connected Edges from the graph.

**Parameters:**
- `id`: ID of the node to remove

**Returns:**
- `error`: If node doesn't exist

### AddEdge

```go
func (g *Graph) AddEdge(from string, to string, weight float64) error
```

Adds a directed weighted edge between two Nodes.

**Parameters:**
- `from`: Source node ID
- `to`: Destination node ID
- `weight`: Edge weight (can be negative)

**Returns:**
- `error`: If either node doesn't exist

### RemoveEdge

```go
func (g *Graph) RemoveEdge(from string, to string) error
```

Removes an edge between two Nodes.

**Parameters:**
- `from`: Source node ID
- `to`: Destination node ID

**Returns:**
- `error`: If edge doesn't exist

### GetNode

```go
func (g *Graph) GetNode(id string) (*Node, error)
```

Retrieves a node by its ID.

**Parameters:**
- `id`: Node ID to retrieve

**Returns:**
- `*Node`: Node if found
- `error`: If node doesn't exist

### GetNeighbors

```go
func (g *Graph) GetNeighbors(id string) ([]*Node, error)
```

Gets all Nodes connected to the specified node by outgoing Edges.

**Parameters:**
- `id`: Node ID to get neighbors for

**Returns:**
- `[]*Node`: Slice of neighbor Nodes
- `error`: If node doesn't exist

### SaveToFile

```go
func (g *Graph) SaveToFile(filename string) error
```

Saves the graph to a JSON file.

**Parameters:**
- `filename`: Path to save the file

**Returns:**
- `error`: If file operations fail

### LoadFromFile

```go
func (g *Graph) LoadFromFile(filename string) error
```

Loads a graph from a JSON file.

**Parameters:**
- `filename`: Path to load the file from

**Returns:**
- `error`: If file operations fail or invalid format

## Error Handling

All methods return specific error types that can be checked using standard Go error handling:

```go
if err := g.AddNode("1", data); err != nil {
    switch {
    case strings.Contains(err.Error(), "already exists"):
        // Handle duplicate node
    case strings.Contains(err.Error(), "invalid"):
        // Handle invalid parameters
    default:
        // Handle other errors
    }
}
```

## Examples

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/yourusername/gropher"
)

func main() {
    g := gropher.New()

    // Add Nodes
    g.AddNode("1", "Node 1")
    g.AddNode("2", "Node 2")
    g.AddNode("3", "Node 3")

    // Add Edges
    g.AddEdge("1", "2", 1.0)
    g.AddEdge("2", "3", 2.0)

    // Get neighbors
    neighbors, _ := g.GetNeighbors("1")
    for _, neighbor := range neighbors {
        fmt.Printf("Node %s is connected to %s\n", "1", neighbor.ID)
    }
}
```

### Custom Data Types

```go
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    g := gropher.New()

    // Add Nodes with custom data
    g.AddNode("1", Person{Name: "Alice", Age: 30})
    g.AddNode("2", Person{Name: "Bob", Age: 25})

    // Retrieve and use custom data
    node, _ := g.GetNode("1")
    person := node.Data.(Person)
    fmt.Printf("Name: %s, Age: %d\n", person.Name, person.Age)
}
```

### File Operations

```go
func main() {
    g := gropher.New()

    // Create graph...
    g.AddNode("1", "Data")

    // Save to file
    err := g.SaveToFile("graph.json")
    if err != nil {
        log.Fatal(err)
    }

    // Load in new graph
    newGraph := gropher.New()
    err = newGraph.LoadFromFile("graph.json")
    if err != nil {
        log.Fatal(err)
    }
}
```

### Concurrent Operations

```go
func main() {
    g := gropher.New()

    // Safe for concurrent use
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            id := fmt.Sprintf("node%d", i)
            g.AddNode(id, fmt.Sprintf("Data %d", i))
        }(i)
    }
    wg.Wait()
}
```