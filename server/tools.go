package server

import (
	"context"

	"github.com/cnosuke/mcp-greeting/greeter"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// GreetingHelloArgs - Arguments for greeting_hello tool
type GreetingHelloArgs struct {
	Name string `json:"name" jsonschema:"Optional name for personalized greeting"`
}

// RegisterAllTools - Register all tools with the server
func RegisterAllTools(mcpServer *mcp.Server, greeter *greeter.Greeter) error {
	if err := registerGreetingHelloTool(mcpServer, greeter); err != nil {
		return err
	}

	return nil
}

// registerGreetingHelloTool - Register the greeting_hello tool
func registerGreetingHelloTool(mcpServer *mcp.Server, greeter *greeter.Greeter) error {
	zap.S().Debugw("registering greeting_hello tool")

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "greeting_hello",
		Description: "Generate a greeting message",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GreetingHelloArgs) (*mcp.CallToolResult, struct{}, error) {
		zap.S().Debugw("executing greeting_hello",
			"name", input.Name)

		greeting, err := greeter.GenerateGreeting(input.Name)
		if err != nil {
			zap.S().Errorw("failed to generate greeting",
				"name", input.Name,
				"error", err)
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
			}, struct{}{}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: greeting}},
		}, struct{}{}, nil
	})

	return nil
}
