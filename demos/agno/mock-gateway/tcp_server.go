package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

// MCPRequest represents an MCP JSON-RPC request
type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// MCPResponse represents an MCP JSON-RPC response
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents an MCP JSON-RPC error
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCPServer handles MCP protocol over TCP
type MCPServer struct {
	mocks map[string]MockResponse
	mu    sync.RWMutex
}

// NewMCPServer creates a new MCP server
func NewMCPServer() *MCPServer {
	return &MCPServer{
		mocks: make(map[string]MockResponse),
	}
}

// LoadMocks loads mock configurations
func (s *MCPServer) LoadMocks(mocksDir string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	gateway := NewMockGateway()
	if err := gateway.LoadMocks(mocksDir); err != nil {
		return err
	}

	s.mocks = gateway.mocks
	return nil
}

// HandleConnection handles a single TCP connection
func (s *MCPServer) HandleConnection(conn net.Conn) {
	defer conn.Close()

	log.Printf("New MCP connection from %s", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	encoder := json.NewEncoder(conn)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		log.Printf("Received: %s", line)

		var request MCPRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			log.Printf("Failed to parse request: %v", err)
			response := MCPResponse{
				JSONRPC: "2.0",
				ID:      nil,
				Error: &MCPError{
					Code:    -32700,
					Message: "Parse error",
				},
			}
			encoder.Encode(response)
			continue
		}

		response := s.handleRequest(request)
		if err := encoder.Encode(response); err != nil {
			log.Printf("Failed to send response: %v", err)
			break
		}

		log.Printf("Sent response for method: %s", request.Method)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Connection error: %v", err)
	}

	log.Printf("Connection closed: %s", conn.RemoteAddr())
}

// handleRequest processes an MCP request and returns a response
func (s *MCPServer) handleRequest(request MCPRequest) MCPResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	switch request.Method {
	case "initialize":
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
				"serverInfo": map[string]interface{}{
					"name":    "mock-gateway",
					"version": "1.0.0",
				},
			},
		}

	case "tools/list":
		tools := make([]map[string]interface{}, 0)
		for key, _ := range s.mocks {
			parts := strings.Split(key, ":")
			if len(parts) == 2 {
				// Use just the method name (second part) as the tool name
				// This matches what the agent expects after splitting "mcp/tool:method"
				toolName := parts[1]
				tools = append(tools, map[string]interface{}{
					"name":        toolName,
					"description": fmt.Sprintf("Mock tool for %s", key),
					"inputSchema": map[string]interface{}{
						"type":       "object",
						"properties": map[string]interface{}{},
					},
				})
			}
		}

		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Result: map[string]interface{}{
				"tools": tools,
			},
		}

	case "tools/call":
		var params struct {
			Name      string          `json:"name"`
			Arguments json.RawMessage `json:"arguments"`
		}

		if err := json.Unmarshal(request.Params, &params); err != nil {
			return MCPResponse{
				JSONRPC: "2.0",
				ID:      request.ID,
				Error: &MCPError{
					Code:    -32602,
					Message: "Invalid params",
				},
			}
		}

		// Find the mock by method name (params.Name should be just the method like "list_issues")
		var foundMock *MockResponse
		for key, mock := range s.mocks {
			parts := strings.Split(key, ":")
			if len(parts) == 2 && parts[1] == params.Name {
				foundMock = &mock
				break
			}
		}

		if foundMock != nil {
			return MCPResponse{
				JSONRPC: "2.0",
				ID:      request.ID,
				Result: map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"type": "text",
							"text": fmt.Sprintf("%v", foundMock.Response),
						},
					},
				},
			}
		}

		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Tool not found: %s", params.Name),
			},
		}

	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", request.Method),
			},
		}
	}
}

// StartTCPServer starts the MCP TCP server
func (s *MCPServer) StartTCPServer(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to start TCP server: %w", err)
	}
	defer listener.Close()

	log.Printf("MCP TCP server listening on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go s.HandleConnection(conn)
	}
}
