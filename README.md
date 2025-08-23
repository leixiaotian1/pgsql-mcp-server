# PostgreSQL MCP Server
[![GoDoc](https://pkg.go.dev/badge/github.com/sql-mcp-server/.svg)](https://pkg.go.dev/github.com/leixiaotian1/pgsql-mcp-server/)
![Stars](https://img.shields.io/github/stars/leixiaotian1/pgsql-mcp-server)
![Forks](https://img.shields.io/github/forks/leixiaotian1/pgsql-mcp-server)

[中文](readme_zh_CH.md) | English

A Model Context Protocol (MCP) server that provides tools for interacting with a PostgreSQL database. It enables AI assistants to execute SQL queries, explain statements, create tables, and list database tables via the MCP protocol.

## ✨ Features

*   **Interact with Databases via AI:** Enables LLMs to perform database operations through a structured protocol.
*   **Secure Toolset:** Separates read and write operations into distinct, authorizable tools (`read_query`, `write_query`).
*   **Schema Management:** Allows for table creation (`create_table`) and listing (`list_tables`).
*   **Query Analysis:** Provides a tool to analyze query execution plans (`explain_query`).
*   **Multiple Transport Modes:** Supports `stdio`, Server-Sent Events (`sse`), and `streamableHttp` for flexible client integration.
*   **Environment-Based Configuration:** Easily configurable using a `.env` file.

## 🛠️ Available Tools

The server exposes the following tools for MCP clients to invoke:

| Tool Name       | Description                                                | Parameters                                                                   |
| --------------- | ---------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `read_query`    | Executes a `SELECT` SQL query.                             | `query` (string, required): The `SELECT` statement to execute.               |
| `write_query`   | Executes an `INSERT`, `UPDATE`, or `DELETE` SQL query.     | `query` (string, required): The `INSERT/UPDATE/DELETE` statement to execute. |
| `create_table`  | Executes a `CREATE TABLE` SQL statement.                   | `schema` (string, required): The `CREATE TABLE` statement.                   |
| `list_tables`   | Lists all user-created tables in the database.             | `schema` (string, optional): The schema name to filter tables by.            |
| `explain_query` | Returns the execution plan for a given SQL query.          | `query` (string, required): The query to explain (must start with `EXPLAIN`).|

## 🚀 Quick Start

### Prerequisites

- Go 1.23 or later
- A PostgreSQL database server

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/sql-mcp-server.git
   cd sql-mcp-server
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the MCP server:
   ```bash
   go build -o sql-mcp-server
   ```

## Configuration

The `pg-mcp-server` requires database connection details to be provided via environment variables. Create a `.env` file in the project root with the following variables:

```
DB_HOST=localhost      # PostgreSQL server host
DB_PORT=5432           # PostgreSQL server port
DB_NAME=postgres       # Database name
DB_USER=your_username  # Database user
DB_PASSWORD=your_pass  # Database password
DB_SSLMODE=disable     # SSL mode (disable, require, verify-ca, verify-full)
SERVER_MODE=stdio      # Server mode (stdio, sse, streamableHttp)
```

## Usage

### Running the Server

```bash
./sql-mcp-server
```

### MCP Configuration

To use this server with an MCP-enabled AI assistant, add the following to your MCP configuration:

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
        "DB_SSLMODE": "disable",
        "SERVER_MODE": "stdio"
      },
      "disabled": false,
      "autoApprove": []
    }
  }
}
```
---


### DOCKER DEPLOYMENT

<details>
<summary><strong>Click to expand Docker Deployment Guide</strong></summary>


#### Prerequisites

- Docker installed

#### Deployment Steps

1.  **Clone the project**
    ```bash
    git clone https://github.com/leixiaotian1/pgsql-mcp-server.git
    cd pgsql-mcp-server
    ```

2.  **Configure `.env` file**
    
    Create a `.env` file in the project root directory. This file stores database connection information. **Ensure the `DB_HOST` value matches the database container name you'll start later.**
    
    ```properties
    DB_HOST=postgres
    DB_PORT=5432
    DB_NAME=postgres
    DB_USER=user
    DB_PASSWORD=password
    DB_SSLMODE=disable
    SERVER_MODE=sse
    ```

3.  **Create Docker network**

    To enable communication between the application container and database container, create a shared Docker network. This command only needs to run once.
    ```bash
    docker network create sql-mcp-network
    ```

4.  **Start PostgreSQL database container**

    Use this command to start a PostgreSQL container and connect it to our network.
    
    > **Note:**
    > - `--name postgres`: Container name, must exactly match the `DB_HOST` in your `.env` file.
    > - `--network sql-mcp-network`: Connect to the shared network.
    > - `-p 5432:5432`: Maps host's `5432` port to container's `5432` port. This means you can connect from your computer (e.g., using DBeaver) via `localhost:5432`, while the app container will access `5432` port directly through the internal network.

    ```bash
    docker run -d \
      --name postgres \
      --network sql-mcp-network \
      -e POSTGRES_USER=user \
      -e POSTGRES_PASSWORD=password \
      -e POSTGRES_DB=postgres \
      -p 5432:5432 \
      postgres
    ```

5.  **Build and run the application**

    Now you can use commands from the `Makefile` to manage the application.

    - **Build image and run container:**
      ```bash
      make build
      make run
      ```
      This will automatically stop old containers, build a new image, and start a new container.

    - **View application logs:**
      ```bash
      make logs
      ```
      If you see `Successfully connected to database`, everything is working correctly.

    - **Stop the application:**
      ```bash
      make stop
      ```

</details>


---

## 🔌 Server Modes

You can select the transport protocol by setting the `SERVER_MODE` environment variable.

### `stdio`

The server communicates over standard input and output. This is the default mode and is ideal for local testing or direct integration with command-line-based MCP clients.

### `sse`

The server communicates using Server-Sent Events (SSE). When this mode is enabled, the server will start an HTTP service and listen for connections.

*   **SSE Endpoint:** `http://localhost:8088/sse`
*   **Message Endpoint:** `http://localhost:8088/message`

### `streamableHttp`

The server uses the Streamable HTTP transport, a more modern and flexible HTTP-based transport for MCP.

*   **Endpoint:** `http://localhost:8088/mcp`

## 🤝 Contributing

Contributions are welcome! If you find any bugs, have feature requests, or suggestions for improvement, please feel free to submit a Pull Request or open an Issue.

1.  Fork the Project.
2.  Create your Feature Branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your Changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the Branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.

## 📄 License

This project is open source and is licensed under the MIT License.
