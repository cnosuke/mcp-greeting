package server

import (
	"context"

	"github.com/cnosuke/mcp-greeting/config"
	"github.com/cnosuke/mcp-greeting/greeter"
	ierrors "github.com/cnosuke/mcp-greeting/internal/errors"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// RunStdio - Execute the MCP server with STDIO transport
func RunStdio(cfg *config.Config, name string, version string, revision string) error {
	zap.S().Infow("starting MCP Greeting Server with STDIO transport")

	mcpServer, err := createMCPServer(cfg, name, version, revision)
	if err != nil {
		return err
	}

	zap.S().Infow("starting MCP server with STDIO")
	err = mcpServer.Run(context.Background(), &mcp.StdioTransport{})
	if err != nil {
		zap.S().Errorw("failed to start STDIO server", "error", err)
		return ierrors.Wrap(err, "failed to start STDIO server")
	}

	zap.S().Infow("STDIO server shutting down")
	return nil
}

// createMCPServer - Create MCP server instance with common configuration
func createMCPServer(cfg *config.Config, name string, version string, revision string) (*mcp.Server, error) {
	versionString := version
	if revision != "" && revision != "xxx" {
		versionString = versionString + " (" + revision + ")"
	}

	zap.S().Debugw("creating Greeter")
	greeterInstance, err := greeter.NewGreeter(cfg)
	if err != nil {
		zap.S().Errorw("failed to create Greeter", "error", err)
		return nil, err
	}

	zap.S().Debugw("creating MCP server",
		"name", name,
		"version", versionString,
	)
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    name,
		Version: versionString,
	}, nil)

	zap.S().Debugw("registering tools")
	if err := RegisterAllTools(mcpServer, greeterInstance); err != nil {
		zap.S().Errorw("failed to register tools", "error", err)
		return nil, err
	}

	return mcpServer, nil
}
