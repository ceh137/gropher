package gropher

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Node represents a vertex in the graph with generic data
type Node struct {
	ID   string      `json:"id"`
	Data interface{} `json:"data"`
}

// Edge represents a connection between two Nodes
type Edge struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Weight float64 `json:"weight"`
}

// Graph represents a directed weighted graph
type Graph struct {
	Nodes map[string]*Node
	Edges map[string]map[string]*Edge
	mu    sync.RWMutex
}

// New creates a new empty graph
func New() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
		Edges: make(map[string]map[string]*Edge),
	}
}

// AddNode adds a new node to the graph
func (g *Graph) AddNode(id string, data interface{}) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.Nodes[id]; exists {
		return fmt.Errorf("node with ID %s already exists", id)
	}

	g.Nodes[id] = &Node{
		ID:   id,
		Data: data,
	}

	g.Edges[id] = make(map[string]*Edge)
	return nil
}

// RemoveNode removes a node and all its Edges from the graph
func (g *Graph) RemoveNode(id string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.Nodes[id]; !exists {
		return fmt.Errorf("node with ID %s does not exist", id)
	}

	// Remove all Edges connected to this node
	delete(g.Edges, id)
	for _, edges := range g.Edges {
		delete(edges, id)
	}

	delete(g.Nodes, id)
	return nil
}

// AddEdge adds a new edge between two Nodes
func (g *Graph) AddEdge(from, to string, weight float64) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.Nodes[from]; !exists {
		return fmt.Errorf("source node %s does not exist", from)
	}
	if _, exists := g.Nodes[to]; !exists {
		return fmt.Errorf("destination node %s does not exist", to)
	}

	g.Edges[from][to] = &Edge{
		From:   from,
		To:     to,
		Weight: weight,
	}
	return nil
}

// RemoveEdge removes an edge between two Nodes
func (g *Graph) RemoveEdge(from, to string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.Edges[from][to]; !exists {
		return fmt.Errorf("edge from %s to %s does not exist", from, to)
	}

	delete(g.Edges[from], to)
	return nil
}

// GetNode returns a node by its ID
func (g *Graph) GetNode(id string) (*Node, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	node, exists := g.Nodes[id]
	if !exists {
		return nil, fmt.Errorf("node with ID %s does not exist", id)
	}
	return node, nil
}

// GetNeighbors returns all Nodes connected to the given node
func (g *Graph) GetNeighbors(id string) ([]*Node, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if _, exists := g.Nodes[id]; !exists {
		return nil, fmt.Errorf("node with ID %s does not exist", id)
	}

	var neighbors []*Node
	for toID := range g.Edges[id] {
		neighbors = append(neighbors, g.Nodes[toID])
	}
	return neighbors, nil
}

// graphData is used for JSON serialization
type graphData struct {
	Nodes []*Node                     `json:"Nodes"`
	Edges map[string]map[string]*Edge `json:"Edges"`
}

// SaveToFile saves the graph to a JSON file
func (g *Graph) SaveToFile(filename string) error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	data := graphData{
		Edges: g.Edges,
	}

	for _, node := range g.Nodes {
		data.Nodes = append(data.Nodes, node)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode graph: %v", err)
	}

	return nil
}

// LoadFromFile loads a graph from a JSON file
func (g *Graph) LoadFromFile(filename string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var data graphData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("failed to decode graph: %v", err)
	}

	// Clear existing graph
	g.Nodes = make(map[string]*Node)
	g.Edges = make(map[string]map[string]*Edge)

	// Restore Nodes
	for _, node := range data.Nodes {
		g.Nodes[node.ID] = node
		g.Edges[node.ID] = make(map[string]*Edge)
	}

	// Restore Edges
	for fromID, edges := range data.Edges {
		for toID, edge := range edges {
			g.Edges[fromID][toID] = edge
		}
	}

	return nil
}
