# MCP Greeting Server

MCP Greeting Server is a Go-based MCP server implementation that provides basic greeting functionality, allowing MCP clients (e.g., Claude Desktop) to generate greeting messages.

## Features

* MCP Compliance: Provides a JSON‐RPC based interface for tool execution according to the MCP specification.
* Greeting Operations: Supports generating greeting messages, with options for personalization.
* Transport Options: Supports both STDIO and Streamable HTTP transports.

## Requirements

- Docker (recommended)

For local development:

- Go 1.24 or later

## Using with Docker (Recommended)

```bash
docker pull cnosuke/mcp-greeting:latest

docker run -i --rm cnosuke/mcp-greeting:latest
```

### Using with Claude Desktop (Docker)

To integrate with Claude Desktop using Docker, add an entry to your `claude_desktop_config.json` file:

```json
{
  "mcpServers": {
    "greeting": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "cnosuke/mcp-greeting:latest"]
    }
  }
}
```

## Building and Running (Go Binary)

Alternatively, you can build and run the Go binary directly:

```bash
# Build the server
make bin/mcp-greeting
```

### STDIO Transport

Run the server with STDIO transport (used with Claude Desktop):

```bash
./bin/mcp-greeting stdioserver --config=config.yml
```

### HTTP Transport (Streamable HTTP)

Run the server with Streamable HTTP transport:

```bash
./bin/mcp-greeting httpserver --config=config.yml
```

### Using with Claude Desktop (Go Binary)

To integrate with Claude Desktop using the Go binary, add an entry to your `claude_desktop_config.json` file:

```json
{
  "mcpServers": {
    "greeting": {
      "command": "./bin/mcp-greeting",
      "args": ["stdioserver"],
      "env": {
        "LOG_PATH": "mcp-greeting.log",
        "LOG_LEVEL": "info",
        "GREETING_DEFAULT_MESSAGE": "こんにちは"
      }
    }
  }
}
```

## Configuration

The server is configured via a YAML file (default: config.yml). For example:

```yaml
log: 'path/to/mcp-greeting.log' # Log file path, if empty no log will be produced
log_level: 'info'               # debug, info, warn, error

http:
  binding: 'localhost:8080'
  endpoint_path: '/mcp'
  auth_token: ''                # Bearer token for authentication (optional)
  allowed_origins: []           # CORS allowed origins (optional)

greeting:
  default_message: "こんにちは！"
```

Note: The default greeting message can also be injected via an environment variable `GREETING_DEFAULT_MESSAGE`. If this environment variable is set, it will override the value in the configuration file.

You can override configurations using environment variables:
- `LOG_PATH`: Path to log file
- `LOG_LEVEL`: Log level (debug, info, warn, error)
- `HTTP_BINDING`: HTTP server binding address
- `HTTP_ENDPOINT_PATH`: HTTP endpoint path
- `HTTP_AUTH_TOKEN`: Bearer token for authentication
- `HTTP_ALLOWED_ORIGINS`: Comma-separated list of allowed CORS origins
- `GREETING_DEFAULT_MESSAGE`: Default greeting message

## Logging

Logging behavior is controlled through configuration:

- If `log` is set in the config file, logs will be written to the specified file
- If `log` is empty, no logs will be produced
- Set `log_level` to control verbosity: `debug`, `info` (default), `warn`, `error`
- In STDIO mode, console output is suppressed to preserve JSON-RPC protocol integrity
- In HTTP mode, logs are output to console (info/below to stdout, warn/above to stderr)

## MCP Server Usage

MCP clients interact with the server by sending JSON‐RPC requests to execute various tools. The following MCP tools are supported:

* `greeting/hello`: Generates a greeting message, with an optional name parameter for personalization.

## Command-Line Reference

```
COMMANDS:
   stdioserver, stdio, s  Run MCP server with STDIO transport
   httpserver, http       Run MCP server with Streamable HTTP transport

OPTIONS (both commands):
   --config value, -c value  Path to the configuration file (default: "config.yml")
```

## Contributing

Contributions are welcome! Please fork the repository and submit pull requests for improvements or bug fixes. For major changes, open an issue first to discuss your ideas.

## License

This project is licensed under the MIT License.

Author: cnosuke ( x.com/cnosuke )
