# NetworkX Compatibility Examples

## Python (NetworkX) to Gropher

```python
import networkx as nx
import json

# Create a NetworkX graph
G = nx.DiGraph()

# Add Nodes with attributes
G.add_node(1, name="Node 1", value=42)
G.add_node(2, name="Node 2", value=100)
G.add_node(3, name="Node 3", value=75)

# Add Edges with weights
G.add_edge(1, 2, weight=1.5)
G.add_edge(2, 3, weight=2.0)
G.add_edge(3, 1, weight=0.5)

# Save to JSON
data = nx.node_link_data(G)
with open('graph.json', 'w') as f:
    json.dump(data, f, indent=2)
```

```go
package main

import (
    "fmt"
    "github.com/ceh137/gropher"
)

func main() {
    // Create a new graph
    g := gropher.New()

    // Load the NetworkX-generated JSON file
    err := g.LoadNetworkXJSON("graph.json")
    if err != nil {
        panic(err)
    }

    // Access the data
    node1, _ := g.GetNode("1")
    data := node1.Data.(map[string]interface{})
    fmt.Printf("Node 1 name: %s, value: %v\n", data["name"], data["value"])
}
```

## Gropher to NetworkX

```go
package main

import (
    "github.com/ceh137/gropher"
)

func main() {
    // Create a new graph
    g := gropher.New()

    // Add Nodes with data
    g.AddNode("1", map[string]interface{}{
        "name": "Node 1",
        "value": 42,
    })
    g.AddNode("2", map[string]interface{}{
        "name": "Node 2",
        "value": 100,
    })

    // Add Edges
    g.AddEdge("1", "2", 1.5)

    // Save in NetworkX-compatible format
    err := g.SaveNetworkXJSON("graph.json")
    if err != nil {
        panic(err)
    }
}
```

```python
import networkx as nx
import json

# Load the Gropher-generated JSON file
with open('graph.json', 'r') as f:
    data = json.load(f)

# Create NetworkX graph from the data
G = nx.node_link_graph(data)

# Access the data
print(f"Node 1 attributes: {G.Nodes['1']}")
print(f"Edge weight: {G.Edges[('1', '2')]['weight']}")
```

## Data Type Mapping

When converting between Gropher and NetworkX, the following data type mappings are used:

| Gropher | NetworkX |
|---------|----------|
| Node.ID (string) | node id (str) |
| Node.Data (map) | node attributes (dict) |
| Edge.Weight (float64) | edge weight (float) |

## File Format Compatibility

The JSON format used is compatible with NetworkX's `node_link_data()` and `node_link_graph()` functions. The format includes:

```json
{
  "directed": true,
  "multigraph": false,
  "graph": {},
  "Nodes": [
    {
      "id": "1",
      "data": {
        "name": "Node 1",
        "value": 42
      }
    }
  ],
  "links": [
    {
      "source": "1",
      "target": "2",
      "weight": 1.5
    }
  ]
}
```

## Notes

1. NetworkX uses numeric node IDs by default, but they are converted to strings in Gropher.
2. Complex Python objects in NetworkX node/edge attributes should be serializable to JSON.
3. Gropher preserves all NetworkX node and edge attributes in the data field.