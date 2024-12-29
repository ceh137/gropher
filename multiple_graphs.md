# Multiple Graphs Example

## Go Usage

```go
package main

import (
	"fmt"
    "github.com/ceh137/gropher"
)

func main() {
    // Create multiple graphs
    g1 := gropher.New()
    g1.AddNode("1", "Graph 1 Node 1")
    g1.AddNode("2", "Graph 1 Node 2")
    g1.AddEdge("1", "2", 1.0)

    g2 := gropher.New()
    g2.AddNode("A", "Graph 2 Node A")
    g2.AddNode("B", "Graph 2 Node B")
    g2.AddEdge("A", "B", 2.0)

    // Save multiple graphs to file
    graphs := []*gropher.Graph{g1, g2}
    err := gropher.SaveGraphsToFile(graphs, "multiple_graphs.json")
    if err != nil {
        panic(err)
    }

    // Load multiple graphs from file
    loadedGraphs, err := gropher.LoadGraphsFromFile("multiple_graphs.json")
    if err != nil {
        panic(err)
    }

    // Use the loaded graphs
    for i, g := range loadedGraphs {
        fmt.Printf("Graph %d has %d Nodes\n", i+1, len(g.Nodes))
    }
}
```

## Python (NetworkX) Usage

```python
import json
import networkx as nx

# Load the graphs from the file
with open('multiple_graphs.json', 'r') as f:
    data = json.load(f)

# Convert each graph in the collection
graphs = []
for graph_data in data['graphs']:
    G = nx.node_link_graph(graph_data)
    graphs.append(G)

# Use the graphs
for i, G in enumerate(graphs):
    print(f"Graph {i+1} info:")
    print(f"modes: {G.Nodes(data=True)}")
    print(f"Edges: {G.Edges(data=True)}")
    print()

# Create graphs in NetworkX and save them for Gropher
def save_graphs_for_gropher(graphs, filename):
    collection = {
        'graphs': [nx.node_link_data(G) for G in graphs]
    }
    with open(filename, 'w') as f:
        json.dump(collection, f, indent=2)

# Example usage
G1 = nx.DiGraph()
G1.add_edge(1, 2, weight=1.0)
G1.Nodes[1]['data'] = {'name': 'Node 1'}
G1.Nodes[2]['data'] = {'name': 'Node 2'}

G2 = nx.DiGraph()
G2.add_edge('A', 'B', weight=2.0)
G2.Nodes['A']['data'] = {'name': 'Node A'}
G2.Nodes['B']['data'] = {'name': 'Node B'}

save_graphs_for_gropher([G1, G2], 'nx_multiple_graphs.json')
```

## File Format

The file format used for multiple graphs looks like this:

```json
{
  "graphs": [
    {
      "directed": true,
      "multigraph": false,
      "graph": {},
      "Nodes": [
        {
          "id": "1",
          "data": {
            "value": "Graph 1 Node 1"
          }
        },
        {
          "id": "2",
          "data": {
            "value": "Graph 1 Node 2"
          }
        }
      ],
      "links": [
        {
          "source": "1",
          "target": "2",
          "weight": 1.0
        }
      ]
    },
    {
      "directed": true,
      "multigraph": false,
      "graph": {},
      "Nodes": [
        {
          "id": "A",
          "data": {
            "value": "Graph 2 Node A"
          }
        },
        {
          "id": "B",
          "data": {
            "value": "Graph 2 Node B"
          }
        }
      ],
      "links": [
        {
          "source": "A",
          "target": "B",
          "weight": 2.0
        }
      ]
    }
  ]
}
```