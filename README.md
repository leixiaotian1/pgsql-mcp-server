# PostgreSQL MCP Server

[![Build](https://github.com/mark3labs/mcp-go/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/mark3labs/mcp-go/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mark3labs/mcp-go?cache)](https://goreportcard.com/report/github.com/mark3labs/mcp-go)
[![GoDoc](https://pkg.go.dev/badge/github.com/mark3labs/mcp-go.svg)](https://pkg.go.dev/github.com/mark3labs/mcp-go)
![Stars](https://img.shields.io/github/stars/leixiaotian1/pgsql-mcp-server)
![Forks](https://img.shields.io/github/forks/leixiaotian1/pgsql-mcp-server)

English | [中文](readme_zh_CH.md)



A Model Context Protocol (MCP) server that provides tools for interacting with a PostgreSQL database. This server enables AI assistants to execute SQL queries, create tables, and list database tables through the MCP protocol.

## Features

The server provides the following tools:

- **read_query**: Execute SELECT queries on the PostgreSQL database
- **write_query**: Execute INSERT, UPDATE, or DELETE queries on the PostgreSQL database
- **create_table**: Create a new table in the PostgreSQL database
- **list_tables**: List all user tables in the database (with optional schema filtering)

## Installation

### Prerequisites

- Go 1.23 or later
- PostgreSQL database server

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/sql-mcp-server.git
   cd sql-mcp-server
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the server:
   ```bash
   go build -o sql-mcp-server
   ```

## Configuration

The server requires database connection details through environment variables. Create a `.env` file in the project root with the following variables:

```
DB_HOST=localhost      # PostgreSQL server host
DB_PORT=5432           # PostgreSQL server port
DB_NAME=postgres       # Database name
DB_USER=your_username  # Database user
DB_PASSWORD=your_pass  # Database password
DB_SSLMODE=disable     # SSL mode (disable, require, verify-ca, verify-full)
```

## Usage

### Running the Server

```bash
./sql-mcp-server
```

### MCP Configuration

To use this server with an AI assistant that supports MCP, add the following to your MCP configuration:

```json
{
  "mcpServers": {
    "pgsql-mcp-server": {
      "command": "/path/to/sql-mcp-server",
      "args": [],
      "env": {
        "DB_HOST": "localhost",
        "DB_PORT": "5432",
        "DB_NAME": "postgres",
        "DB_USER": "your_username",
        "DB_PASSWORD": "your_password",
        "DB_SSLMODE": "disable"
      },
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

### Tool Examples

#### List Tables

List all user tables in the database:

```json
{
  "server_name": "pgsql-mcp-server",
  "tool_name": "list_tables",
  "arguments": {}
}
```

List tables in a specific schema:

```json
{
  "server_name": "pgsql-mcp-server",
  "tool_name": "list_tables",
  "arguments": {
    "schema": "public"
  }
}
```

#### Create Table

Create a new table:

```json
{
  "server_name": "pgsql-mcp-server",
  "tool_name": "create_table",
  "arguments": {
    "schema": "CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(100), email VARCHAR(100) UNIQUE, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)"
  }
}
```

#### Read Query

Execute a SELECT query:

```json
{
  "server_name": "pgsql-mcp-server",
  "tool_name": "read_query",
  "arguments": {
    "query": "SELECT * FROM users LIMIT 10"
  }
}
```

#### Write Query

Execute an INSERT query:

```json
{
  "server_name": "pgsql-mcp-server",
  "tool_name": "write_query",
  "arguments": {
    "query": "INSERT INTO users (name, email) VALUES ('John Doe', 'john@example.com')"
  }
}
```

Execute an UPDATE query:

```json
{
  "server_name": "pgsql-mcp-server",
  "tool_name": "write_query",
  "arguments": {
    "query": "UPDATE users SET name = 'Jane Doe' WHERE id = 1"
  }
}
```

Execute a DELETE query:

```json
{
  "server_name": "pgsql-mcp-server",
  "tool_name": "write_query",
  "arguments": {
    "query": "DELETE FROM users WHERE id = 1"
  }
}
```

## Security Considerations

- The server validates query types to ensure that only appropriate operations are performed with each tool.
- Input sanitization is performed for schema names to prevent SQL injection.
- Consider using a dedicated database user with limited permissions for this server.
- In production environments, enable SSL by setting `DB_SSLMODE` to `require` or higher.

## Dependencies

- [github.com/joho/godotenv](https://github.com/joho/godotenv) - For loading environment variables from .env file
- [github.com/lib/pq](https://github.com/lib/pq) - PostgreSQL driver for Go
- [github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) - Go SDK for Model Context Protocol

## License

[Add license information here]

## Contributing

[Add contribution guidelines here]
