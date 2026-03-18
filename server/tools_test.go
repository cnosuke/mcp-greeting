package server

import (
	"testing"

	"github.com/cnosuke/mcp-greeting/greeter"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// Test for GreetingHelloArgs
func TestGreetingHelloArgs(t *testing.T) {
	// When name is empty
	argsEmpty := GreetingHelloArgs{
		Name: "",
	}
	assert.Equal(t, "", argsEmpty.Name)

	// When name is set
	argsWithName := GreetingHelloArgs{
		Name: "Test User",
	}
	assert.Equal(t, "Test User", argsWithName.Name)
}

// TestRegisterAllTools - Test tool registration
func TestRegisterAllTools(t *testing.T) {
	// Set up test logger
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	// Create MCP server
	mcpServer := mcp.NewServer(&mcp.Implementation{Name: "test-server", Version: "0.0.1"}, nil)

	// Create mock greeter instance
	cfg := &greeter.Greeter{
		DefaultMessage: "Hello!",
	}

	// Register tools
	err := RegisterAllTools(mcpServer, cfg)
	assert.NoError(t, err)
}
