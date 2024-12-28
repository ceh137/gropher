package gropher

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestNodeOperations(t *testing.T) {
	g := New()

	t.Run("Add Nodes", func(t *testing.T) {
		// Test adding nodes with different data types
		testCases := []struct {
			id   string
			data interface{}
		}{
			{"1", "String data"},
			{"2", 42},
			{"3", map[string]interface{}{"name": "Test", "value": 100}},
			{"4", struct{ Value string }{Value: "Custom struct"}},
		}

		for _, tc := range testCases {
			err := g.AddNode(tc.id, tc.data)
			if err != nil {
				t.Errorf("Failed to add node %s: %v", tc.id, err)
			}

			// Verify node exists
			node, err := g.GetNode(tc.id)
			if err != nil {
				t.Errorf("Failed to get node %s: %v", tc.id, err)
			}
			// Compare data using JSON marshaling to handle complex types
			expectedJSON, err := json.Marshal(tc.data)
			if err != nil {
				t.Errorf("Failed to marshal expected data for node %s: %v", tc.id, err)
			}
			actualJSON, err := json.Marshal(node.Data)
			if err != nil {
				t.Errorf("Failed to marshal actual data for node %s: %v", tc.id, err)
			}
			if string(expectedJSON) != string(actualJSON) {
				t.Errorf("Node %s data mismatch: expected %s, got %s", tc.id, expectedJSON, actualJSON)
			}
		}
	})

	t.Run("Add Duplicate Node", func(t *testing.T) {
		err := g.AddNode("1", "Duplicate")
		if err == nil {
			t.Error("Expected error when adding duplicate node")
		}
	})

	t.Run("Remove Node", func(t *testing.T) {
		// Add edges to test cascade deletion
		err := g.AddEdge("1", "2", 1.0)
		if err != nil {
			t.Errorf("Failed to add edge: %v", err)
		}
		err = g.AddEdge("2", "1", 1.0)
		if err != nil {
			t.Errorf("Failed to add edge: %v", err)
		}

		// Remove node and verify
		err = g.RemoveNode("1")
		if err != nil {
			t.Errorf("Failed to remove node: %v", err)
		}

		// Verify node is gone
		_, err = g.GetNode("1")
		if err == nil {
			t.Error("Expected error when getting removed node")
		}

		// Verify edges are removed
		neighbors, err := g.GetNeighbors("2")
		if err != nil {
			t.Errorf("Failed to get neighbors: %v", err)
		}
		if len(neighbors) != 0 {
			t.Error("Expected no neighbors after node removal")
		}
	})

	t.Run("Remove Non-existent Node", func(t *testing.T) {
		err := g.RemoveNode("nonexistent")
		if err == nil {
			t.Error("Expected error when removing non-existent node")
		}
	})
}

func TestEdgeOperations(t *testing.T) {
	g := New()

	// Setup
	err := g.AddNode("1", "Node 1")
	if err != nil {
		t.Fatal(err)
	}
	err = g.AddNode("2", "Node 2")
	if err != nil {
		t.Fatal(err)
	}
	err = g.AddNode("3", "Node 3")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Add Edges", func(t *testing.T) {
		testCases := []struct {
			from   string
			to     string
			weight float64
		}{
			{"1", "2", 1.0},
			{"2", "3", 2.5},
			{"3", "1", 0.5},
		}

		for _, tc := range testCases {
			err := g.AddEdge(tc.from, tc.to, tc.weight)
			if err != nil {
				t.Errorf("Failed to add edge %s->%s: %v", tc.from, tc.to, err)
			}

			// Verify edge exists
			neighbors, err := g.GetNeighbors(tc.from)
			if err != nil {
				t.Errorf("Failed to get neighbors: %v", err)
			}
			found := false
			for _, neighbor := range neighbors {
				if neighbor.ID == tc.to {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Edge %s->%s not found in neighbors", tc.from, tc.to)
			}
		}
	})

	t.Run("Add Invalid Edges", func(t *testing.T) {
		testCases := []struct {
			from string
			to   string
		}{
			{"1", "nonexistent"},
			{"nonexistent", "1"},
			{"nonexistent1", "nonexistent2"},
		}

		for _, tc := range testCases {
			err := g.AddEdge(tc.from, tc.to, 1.0)
			if err == nil {
				t.Errorf("Expected error when adding edge %s->%s", tc.from, tc.to)
			}
		}
	})

	t.Run("Remove Edges", func(t *testing.T) {
		err := g.RemoveEdge("1", "2")
		if err != nil {
			t.Errorf("Failed to remove edge: %v", err)
		}

		// Verify edge is removed
		neighbors, err := g.GetNeighbors("1")
		if err != nil {
			t.Errorf("Failed to get neighbors: %v", err)
		}
		for _, neighbor := range neighbors {
			if neighbor.ID == "2" {
				t.Error("Edge still exists after removal")
			}
		}
	})

	t.Run("Remove Non-existent Edge", func(t *testing.T) {
		err := g.RemoveEdge("1", "nonexistent")
		if err == nil {
			t.Error("Expected error when removing non-existent edge")
		}
	})
}

func TestGraphSerialization(t *testing.T) {
	g := New()
	filename := "test_graph.json"

	t.Run("Save and Load Empty Graph", func(t *testing.T) {
		err := g.SaveToFile(filename)
		if err != nil {
			t.Errorf("Failed to save empty graph: %v", err)
		}

		newGraph := New()
		err = newGraph.LoadFromFile(filename)
		if err != nil {
			t.Errorf("Failed to load empty graph: %v", err)
		}
	})

	t.Run("Save and Load Complex Graph", func(t *testing.T) {
		// Add various types of data
		testData := []struct {
			id   string
			data interface{}
		}{
			{"1", "String data"},
			{"2", 42},
			{"3", map[string]interface{}{"name": "Test", "value": 100}},
			{"4", struct {
				Name  string `json:"name"`
				Value int    `json:"value"`
			}{"Test Struct", 200}},
		}

		for _, td := range testData {
			err := g.AddNode(td.id, td.data)
			if err != nil {
				t.Fatal(err)
			}
		}

		// Add edges
		edges := []struct {
			from   string
			to     string
			weight float64
		}{
			{"1", "2", 1.0},
			{"2", "3", 2.5},
			{"3", "4", 0.5},
			{"4", "1", 1.5},
		}

		for _, edge := range edges {
			err := g.AddEdge(edge.from, edge.to, edge.weight)
			if err != nil {
				t.Fatal(err)
			}
		}

		// Save graph
		err := g.SaveToFile(filename)
		if err != nil {
			t.Errorf("Failed to save graph: %v", err)
		}

		// Load into new graph
		newGraph := New()
		err = newGraph.LoadFromFile(filename)
		if err != nil {
			t.Errorf("Failed to load graph: %v", err)
		}

		// Verify nodes
		for _, td := range testData {
			node, err := newGraph.GetNode(td.id)
			if err != nil {
				t.Errorf("Failed to get node %s from loaded graph: %v", td.id, err)
			}

			// Compare data using JSON marshaling to handle complex types
			originalJSON, _ := json.Marshal(td.data)
			loadedJSON, _ := json.Marshal(node.Data)
			if string(originalJSON) != string(loadedJSON) {
				t.Errorf("Node %s data mismatch: expected %s, got %s", td.id, originalJSON, loadedJSON)
			}
		}

		// Verify edges
		for _, edge := range edges {
			neighbors, err := newGraph.GetNeighbors(edge.from)
			if err != nil {
				t.Errorf("Failed to get neighbors for node %s: %v", edge.from, err)
			}

			found := false
			for _, neighbor := range neighbors {
				if neighbor.ID == edge.to {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Edge %s->%s not found in loaded graph", edge.from, edge.to)
			}
		}
	})

	t.Run("Load Invalid File", func(t *testing.T) {
		// Create invalid JSON file
		err := os.WriteFile(filename, []byte("invalid json"), 0644)
		if err != nil {
			t.Fatal(err)
		}

		err = g.LoadFromFile(filename)
		if err == nil {
			t.Error("Expected error when loading invalid file")
		}
	})

	t.Run("Load Non-existent File", func(t *testing.T) {
		err := g.LoadFromFile("nonexistent.json")
		if err == nil {
			t.Error("Expected error when loading non-existent file")
		}
	})

	t.Run("Verify File Content", func(t *testing.T) {
		// Create a simple graph with known structure
		g := New()
		err := g.AddNode("test1", "Test Data 1")
		if err != nil {
			t.Fatal(err)
		}
		err = g.AddNode("test2", 42)
		if err != nil {
			t.Fatal(err)
		}
		err = g.AddEdge("test1", "test2", 1.5)
		if err != nil {
			t.Fatal(err)
		}

		// Save the graph
		tempFile := "verify_content.json"
		err = g.SaveToFile(tempFile)
		if err != nil {
			t.Fatal(err)
		}

		// Read the raw file content
		content, err := os.ReadFile(tempFile)
		if err != nil {
			t.Fatal(err)
		}

		// Verify JSON structure
		var data map[string]interface{}
		err = json.Unmarshal(content, &data)
		if err != nil {
			t.Fatal(err)
		}

		// Check for required fields
		if _, ok := data["nodes"]; !ok {
			t.Error("File content missing 'nodes' field")
		}
		if _, ok := data["edges"]; !ok {
			t.Error("File content missing 'edges' field")
		}

		os.Remove(tempFile)
	})

	t.Run("Save With Different File Extensions", func(t *testing.T) {
		g := New()
		err := g.AddNode("1", "test")
		if err != nil {
			t.Fatal(err)
		}

		extensions := []string{".json", ".graph", ".txt"}
		for _, ext := range extensions {
			filename := "test_graph" + ext
			err := g.SaveToFile(filename)
			if err != nil {
				t.Errorf("Failed to save with extension %s: %v", ext, err)
			}
			os.Remove(filename)
		}
	})

	t.Run("Load Corrupted File", func(t *testing.T) {
		// Create a file with invalid JSON
		tempFile := "corrupted.json"
		err := os.WriteFile(tempFile, []byte(`{"nodes": [{"id": "1", "data": "test"}, {invalid}`), 0644)
		if err != nil {
			t.Fatal(err)
		}

		err = g.LoadFromFile(tempFile)
		if err == nil {
			t.Error("Expected error when loading corrupted file")
		}

		os.Remove(tempFile)
	})

	t.Run("Save To Read-Only Directory", func(t *testing.T) {
		// Create a read-only directory
		dirName := "readonly_dir"
		err := os.Mkdir(dirName, 0444)
		if err != nil {
			t.Fatal(err)
		}

		filename := dirName + "/test_graph.json"
		err = g.SaveToFile(filename)
		if err == nil {
			t.Error("Expected error when saving to read-only directory")
		}

		os.RemoveAll(dirName)
	})

	t.Run("Large Graph Serialization", func(t *testing.T) {
		g := New()

		// Create a large graph
		for i := 0; i < 1000; i++ {
			err := g.AddNode(fmt.Sprintf("node%d", i), map[string]interface{}{
				"data":      fmt.Sprintf("large data %d", i),
				"timestamp": time.Now().Unix(),
			})
			if err != nil {
				t.Fatal(err)
			}
		}

		// Add some edges
		for i := 0; i < 999; i++ {
			err := g.AddEdge(
				fmt.Sprintf("node%d", i),
				fmt.Sprintf("node%d", i+1),
				float64(i),
			)
			if err != nil {
				t.Fatal(err)
			}
		}

		tempFile := "large_graph.json"
		err := g.SaveToFile(tempFile)
		if err != nil {
			t.Errorf("Failed to save large graph: %v", err)
		}

		// Load and verify
		newGraph := New()
		err = newGraph.LoadFromFile(tempFile)
		if err != nil {
			t.Errorf("Failed to load large graph: %v", err)
		}

		// Verify node count
		nodeCount := 0
		for range newGraph.nodes {
			nodeCount++
		}
		if nodeCount != 1000 {
			t.Errorf("Expected 1000 nodes, got %d", nodeCount)
		}

		os.Remove(tempFile)
	})

	// Cleanup
	os.Remove(filename)
}

func TestConcurrency(t *testing.T) {
	g := New()

	// Setup initial nodes
	for i := 0; i < 5; i++ {
		err := g.AddNode(string(rune('A'+i)), i)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Run("Concurrent Node Operations", func(t *testing.T) {
		done := make(chan bool)
		for i := 0; i < 100; i++ {
			go func(i int) {
				id := string(rune('A' + (i % 5)))

				// Randomly get or modify nodes
				if i%2 == 0 {
					_, _ = g.GetNode(id)
				} else {
					_ = g.AddNode(string(rune('Z'-i%5)), i)
				}
				done <- true
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("Concurrent Edge Operations", func(t *testing.T) {
		done := make(chan bool)
		for i := 0; i < 100; i++ {
			go func(i int) {
				from := string(rune('A' + (i % 5)))
				to := string(rune('A' + ((i + 1) % 5)))

				// Randomly add or remove edges
				if i%2 == 0 {
					_ = g.AddEdge(from, to, float64(i))
				} else {
					_ = g.RemoveEdge(from, to)
				}
				done <- true
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("Concurrent Mixed Operations", func(t *testing.T) {
		done := make(chan bool)
		for i := 0; i < 100; i++ {
			go func(i int) {
				id := string(rune('A' + (i % 5)))

				switch i % 4 {
				case 0:
					_, _ = g.GetNode(id)
				case 1:
					_ = g.AddNode(string(rune('Z'-i%5)), i)
				case 2:
					to := string(rune('A' + ((i + 1) % 5)))
					_ = g.AddEdge(id, to, float64(i))
				case 3:
					_, _ = g.GetNeighbors(id)
				}
				done <- true
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < 100; i++ {
			<-done
		}
	})
}

func TestGraphTraversal(t *testing.T) {
	g := New()

	// Setup a test graph
	nodes := []string{"A", "B", "C", "D", "E"}
	for _, node := range nodes {
		err := g.AddNode(node, node)
		if err != nil {
			t.Fatal(err)
		}
	}

	edges := []struct {
		from string
		to   string
	}{
		{"A", "B"},
		{"B", "C"},
		{"C", "D"},
		{"D", "E"},
		{"E", "A"},
		{"A", "C"},
		{"B", "D"},
	}

	for _, edge := range edges {
		err := g.AddEdge(edge.from, edge.to, 1.0)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Run("Get All Neighbors", func(t *testing.T) {
		expected := map[string]int{
			"A": 2, // B, C
			"B": 2, // C, D
			"C": 1, // D
			"D": 1, // E
			"E": 1, // A
		}

		for node, expectedCount := range expected {
			neighbors, err := g.GetNeighbors(node)
			if err != nil {
				t.Errorf("Failed to get neighbors for node %s: %v", node, err)
			}
			if len(neighbors) != expectedCount {
				t.Errorf("Expected %d neighbors for node %s, got %d", expectedCount, node, len(neighbors))
			}
		}
	})
}
