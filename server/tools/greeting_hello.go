package tools

import (
	"github.com/cockroachdb/errors"
	mcp "github.com/metoro-io/mcp-golang"
	"go.uber.org/zap"
)

// GreetingHelloArgs - Arguments for greeting/hello tool
type GreetingHelloArgs struct {
	Name string `json:"name" jsonschema:"description=Optional name for personalized greeting"`
}

// Greeter defines the interface for greeting generation
type Greeter interface {
	GenerateGreeting(name string) (string, error)
}

// RegisterGreetingHelloTool - Register the greeting/hello tool
func RegisterGreetingHelloTool(server *mcp.Server, greeter Greeter) error {
	zap.S().Debug("registering greeting/hello tool")
	err := server.RegisterTool("greeting/hello", "Generate a greeting message",
		func(args GreetingHelloArgs) (*mcp.ToolResponse, error) {
			zap.S().Debug("executing greeting/hello",
				zap.String("name", args.Name))

			// Generate greeting
			greeting, err := greeter.GenerateGreeting(args.Name)
			if err != nil {
				zap.S().Error("failed to generate greeting",
					zap.String("name", args.Name),
					zap.Error(err))
				return nil, errors.Wrap(err, "failed to generate greeting")
			}

			return mcp.NewToolResponse(mcp.NewTextContent(greeting)), nil
		})

	if err != nil {
		zap.S().Error("failed to register greeting/hello tool", zap.Error(err))
		return errors.Wrap(err, "failed to register greeting/hello tool")
	}

	return nil
}
