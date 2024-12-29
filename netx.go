package gropher

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
)

// NetworkXJSON represents the JSON structure used by NetworkX
type NetworkXJSON struct {
	Directed   bool                   `json:"directed"`
	Multigraph bool                   `json:"multigraph"`
	Graph      map[string]interface{} `json:"graph"`
	Nodes      []NetworkXNode         `json:"Nodes"`
	Links      []NetworkXLink         `json:"links"`
}

// NetworkXNode represents a node in NetworkX JSON format
type NetworkXNode struct {
	ID   string                 `json:"id"`
	Data map[string]interface{} `json:"data,omitempty"`
}

// NetworkXLink represents an edge in NetworkX JSON format
type NetworkXLink struct {
	Source string  `json:"source"`
	Target string  `json:"target"`
	Weight float64 `json:"weight,omitempty"`
}

// SaveNetworkXJSON saves the graph in NetworkX-compatible JSON format
func (g *Graph) SaveNetworkXJSON(filename string) error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	nx := NetworkXJSON{
		Directed:   true,
		Multigraph: false,
		Graph:      make(map[string]interface{}),
	}

	// Convert Nodes
	for id, node := range g.Nodes {
		nxNode := NetworkXNode{
			ID: id,
		}

		// Convert node data to map if possible
		if data, ok := node.Data.(map[string]interface{}); ok {
			nxNode.Data = data
		} else {
			// If data is not a map, store it under a "value" key
			nxNode.Data = map[string]interface{}{
				"value": node.Data,
			}
		}

		nx.Nodes = append(nx.Nodes, nxNode)
	}

	// Convert Edges
	for from, edges := range g.Edges {
		for to, edge := range edges {
			nx.Links = append(nx.Links, NetworkXLink{
				Source: from,
				Target: to,
				Weight: edge.Weight,
			})
		}
	}

	// Write to file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(nx); err != nil {
		return fmt.Errorf("failed to encode NetworkX JSON: %v", err)
	}

	return nil
}

// LoadNetworkXJSON loads a graph from NetworkX-compatible JSON format
func (g *Graph) LoadNetworkXJSON(filename string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var nx NetworkXJSON
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&nx); err != nil {
		return fmt.Errorf("failed to decode NetworkX JSON: %v", err)
	}

	// Clear existing graph
	g.Nodes = make(map[string]*Node)
	g.Edges = make(map[string]map[string]*Edge)

	// Add Nodes
	for _, nxNode := range nx.Nodes {
		var data interface{} = nxNode.Data
		// If data has only a "value" key, simplify it
		if val, ok := nxNode.Data["value"]; ok && len(nxNode.Data) == 1 {
			data = val
		}

		g.Nodes[nxNode.ID] = &Node{
			ID:   nxNode.ID,
			Data: data,
		}
		g.Edges[nxNode.ID] = make(map[string]*Edge)
	}

	// Add Edges
	for _, link := range nx.Links {
		g.Edges[link.Source][link.Target] = &Edge{
			From:   link.Source,
			To:     link.Target,
			Weight: link.Weight,
		}
	}

	return nil
}

// GraphML support can be added here similarly
type GraphML struct {
	XMLName xml.Name `xml:"graphml"`
	// Add GraphML structure
}

// Future implementation for GraphML format
func (g *Graph) SaveGraphML(filename string) error {
	// Implement GraphML export
	return nil
}

func (g *Graph) LoadGraphML(filename string) error {
	// Implement GraphML import
	return nil
}

// MultiGraphCollection represents a collection of graphs in NetworkX format
type MultiGraphCollection struct {
	Graphs []NetworkXJSON `json:"graphs"`
}

// ConvertToNetworkX converts a single graph to NetworkX format
func ConvertToNetworkX(g *Graph) (NetworkXJSON, error) {
	if g == nil {
		return NetworkXJSON{}, fmt.Errorf("nil graph provided")
	}

	nx := NetworkXJSON{
		Directed:   true,
		Multigraph: false,
		Graph:      make(map[string]interface{}),
	}

	// Convert Nodes
	for id, node := range g.Nodes {
		nxNode := NetworkXNode{
			ID: id,
		}

		// Convert node data to map if possible
		if data, ok := node.Data.(map[string]interface{}); ok {
			nxNode.Data = data
		} else {
			// If data is not a map, store it under a "value" key
			nxNode.Data = map[string]interface{}{
				"value": node.Data,
			}
		}

		nx.Nodes = append(nx.Nodes, nxNode)
	}

	// Convert Edges
	for from, edges := range g.Edges {
		for to, edge := range edges {
			nx.Links = append(nx.Links, NetworkXLink{
				Source: from,
				Target: to,
				Weight: edge.Weight,
			})
		}
	}

	return nx, nil
}

// addGraphsToFile converts multiple graphs to NetworkX format and saves them to a file
func SaveGraphsToFile(graphs []*Graph, filename string) error {
	collection := MultiGraphCollection{
		Graphs: make([]NetworkXJSON, 0, len(graphs)),
	}

	// Convert each graph
	for i, g := range graphs {
		if g == nil {
			return fmt.Errorf("nil graph at index %d", i)
		}

		// Lock the graph for reading
		g.mu.RLock()
		nx, err := ConvertToNetworkX(g)
		g.mu.RUnlock()

		if err != nil {
			return fmt.Errorf("failed to convert graph %d: %v", i, err)
		}

		collection.Graphs = append(collection.Graphs, nx)
	}

	// Create or truncate the file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Write with pretty formatting
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(collection); err != nil {
		return fmt.Errorf("failed to encode graphs: %v", err)
	}

	return nil
}

// LoadGraphsFromFile loads multiple graphs from a NetworkX format file
func LoadGraphsFromFile(filename string) ([]*Graph, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var collection MultiGraphCollection
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&collection); err != nil {
		return nil, fmt.Errorf("failed to decode graphs: %v", err)
	}

	// Convert each NetworkX graph to Gropher graph
	graphs := make([]*Graph, 0, len(collection.Graphs))
	for _, nx := range collection.Graphs {
		g := New()

		// Add Nodes
		for _, nxNode := range nx.Nodes {
			var data interface{} = nxNode.Data
			// If data has only a "value" key, simplify it
			if val, ok := nxNode.Data["value"]; ok && len(nxNode.Data) == 1 {
				data = val
			}

			g.Nodes[nxNode.ID] = &Node{
				ID:   nxNode.ID,
				Data: data,
			}
			g.Edges[nxNode.ID] = make(map[string]*Edge)
		}

		// Add Edges
		for _, link := range nx.Links {
			g.Edges[link.Source][link.Target] = &Edge{
				From:   link.Source,
				To:     link.Target,
				Weight: link.Weight,
			}
		}

		graphs = append(graphs, g)
	}

	return graphs, nil
}
