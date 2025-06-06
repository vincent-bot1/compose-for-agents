package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// MockResponse represents a mock response configuration
type MockResponse struct {
	Tool     string                 `json:"tool"`
	Method   string                 `json:"method"`
	Response map[string]interface{} `json:"response"`
}

// MockGateway handles mock responses for MCP tools
type MockGateway struct {
	mocks map[string]MockResponse
}

// NewMockGateway creates a new mock gateway instance
func NewMockGateway() *MockGateway {
	return &MockGateway{
		mocks: make(map[string]MockResponse),
	}
}

// LoadMocks loads mock configurations from JSON files in the mocks directory
func (mg *MockGateway) LoadMocks(mocksDir string) error {
	return filepath.Walk(mocksDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read mock file %s: %w", path, err)
		}

		var mock MockResponse
		if err := json.Unmarshal(data, &mock); err != nil {
			return fmt.Errorf("failed to parse mock file %s: %w", path, err)
		}

		key := fmt.Sprintf("%s:%s", mock.Tool, mock.Method)
		mg.mocks[key] = mock
		log.Printf("Loaded mock for %s", key)

		return nil
	})
}

// HandleMCPCall handles incoming MCP tool calls
func (mg *MockGateway) HandleMCPCall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request path to extract tool and method
	// Expected format: /mcp/{tool}/{method}
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 || pathParts[0] != "mcp" {
		http.Error(w, "Invalid path format. Expected: /mcp/{tool}/{method}", http.StatusBadRequest)
		return
	}

	tool := pathParts[1]
	method := pathParts[2]
	key := fmt.Sprintf("%s:%s", tool, method)

	log.Printf("Received MCP call: %s", key)

	// Read request body for logging
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	log.Printf("Request body: %s", string(body))

	// Find matching mock
	mock, exists := mg.mocks[key]
	if !exists {
		log.Printf("No mock found for %s, returning default response", key)
		// Return a default response if no mock is found
		defaultResponse := map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("No mock configured for %s", key),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(defaultResponse)
		return
	}

	// Return the mock response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mock.Response); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Printf("Returned mock response for %s", key)
}

// HandleHealth provides a health check endpoint
func (mg *MockGateway) HandleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status": "healthy",
		"mocks":  len(mg.mocks),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleListMocks lists all loaded mocks
func (mg *MockGateway) HandleListMocks(w http.ResponseWriter, r *http.Request) {
	mocks := make([]string, 0, len(mg.mocks))
	for key := range mg.mocks {
		mocks = append(mocks, key)
	}

	response := map[string]interface{}{
		"mocks": mocks,
		"count": len(mocks),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	tcpPort := os.Getenv("TCP_PORT")
	if tcpPort == "" {
		tcpPort = "8081"
	}

	mocksDir := os.Getenv("MOCKS_DIR")
	if mocksDir == "" {
		mocksDir = "./mocks"
	}

	// Create HTTP gateway
	gateway := NewMockGateway()
	if err := gateway.LoadMocks(mocksDir); err != nil {
		log.Printf("Warning: Failed to load mocks: %v", err)
	}

	// Create TCP MCP server
	mcpServer := NewMCPServer()
	if err := mcpServer.LoadMocks(mocksDir); err != nil {
		log.Printf("Warning: Failed to load mocks for MCP server: %v", err)
	}

	// Start HTTP server in a goroutine
	go func() {
		// Setup HTTP routes
		http.HandleFunc("/health", gateway.HandleHealth)
		http.HandleFunc("/mocks", gateway.HandleListMocks)
		http.HandleFunc("/mcp/", gateway.HandleMCPCall)

		log.Printf("HTTP Mock Gateway starting on port %s", httpPort)
		log.Printf("Mocks directory: %s", mocksDir)
		log.Printf("Loaded %d mocks", len(gateway.mocks))

		if err := http.ListenAndServe(":"+httpPort, nil); err != nil {
			log.Fatal("HTTP Server failed to start:", err)
		}
	}()

	// Start TCP MCP server
	log.Printf("MCP TCP Server starting on port %s", tcpPort)
	if err := mcpServer.StartTCPServer(tcpPort); err != nil {
		log.Fatal("TCP Server failed to start:", err)
	}
}
