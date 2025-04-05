package tools

import (
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAllTools - Register all tools with the server
func RegisterAllTools(mcpServer *server.MCPServer, greeter Greeter) error {
	// Register greeting/hello tool
	if err := RegisterGreetingHelloTool(mcpServer, greeter); err != nil {
		return err
	}

	return nil
}
