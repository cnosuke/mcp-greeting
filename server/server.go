package server

import (
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"

	"github.com/cnosuke/mcp-greeting/config"
	"github.com/cnosuke/mcp-greeting/server/tools"
	"github.com/cockroachdb/errors"
)

// Run - Execute the MCP server
func Run(cfg *config.Config) error {
	zap.S().Infow("starting MCP Greeting Server")

	// Create Greeting server
	zap.S().Debugw("creating Greeting server")
	greetingServer, err := NewGreetingServer(cfg)
	if err != nil {
		zap.S().Errorw("failed to create Greeting server", "error", err)
		return err
	}

	// Create MCP server with server name and version
	zap.S().Debugw("creating MCP server")
	mcpServer := server.NewMCPServer(
		"MCP Greeting Server",
		"1.0.0",
		server.WithLogging(),
	)

	// Register all tools
	zap.S().Debugw("registering tools")
	if err := tools.RegisterAllTools(mcpServer, greetingServer); err != nil {
		zap.S().Errorw("failed to register tools", "error", err)
		return err
	}

	// Start the server with stdio transport
	zap.S().Infow("starting MCP server")
	err = server.ServeStdio(mcpServer)
	if err != nil {
		zap.S().Errorw("failed to start server", "error", err)
		return errors.Wrap(err, "failed to start server")
	}

	// ServeStdio will block until the server is terminated
	zap.S().Infow("server shutting down")
	return nil
}
