# PostgreSQL MCP 服务器
![Stars](https://img.shields.io/github/stars/leixiaotian1/pgsql-mcp-server)
![Forks](https://img.shields.io/github/forks/leixiaotian1/pgsql-mcp-server)

中文 | [English](README.md)

# PostgreSQL MCP 服务器

一个提供PostgreSQL数据库交互工具的模型上下文协议(MCP)服务。该服务使AI助手能够通过MCP协议执行SQL查询、解释sql语句、创建表和列出数据库表。

## ✨ 特性

*   **通过 AI 与数据库交互:** 使 LLM 能够通过结构化协议执行数据库操作。
*   **安全的工具集:** 将读写操作分离到不同的、可授权的工具中 (`read_query`, `write_query`)。
*   **模式管理:** 允许创建表 (`create_table`) 和查看表 (`list_tables`)。
*   **查询分析:** 提供分析查询执行计划的工具 (`explain_query`)。
*   **多种传输方式:** 支持 `stdio`、服务器发送事件 (`sse`) 和 `streamableHttp`，以实现灵活的客户端集成。
*   **基于环境的配置:** 使用 `.env` 文件即可轻松配置。

## 🛠️ 可用工具

服务器为 MCP 客户端暴露了以下可调用的工具：

| 工具名称        | 描述                                                       | 参数                                                                 |
| --------------- | ---------------------------------------------------------- | -------------------------------------------------------------------- |
| `read_query`    | 执行 `SELECT` SQL 查询。                                   | `query` (字符串, 必需): 要执行的 `SELECT` 语句。                     |
| `write_query`   | 执行 `INSERT`、`UPDATE` 或 `DELETE` SQL 查询。             | `query` (字符串, 必需): 要执行的 `INSERT/UPDATE/DELETE` 语句。       |
| `create_table`  | 执行 `CREATE TABLE` SQL 语句。                             | `schema` (字符串, 必需): `CREATE TABLE` 语句。                       |
| `list_tables`   | 列出数据库中所有用户创建的表。                             | `schema` (字符串, 可选): 用于过滤表的模式（schema）名称。            |
| `explain_query` | 返回给定 SQL 查询的执行计划。                              | `schema` (字符串, 必需): 需要解释的查询（必须以 `EXPLAIN` 开头）。 |

## 🚀 快速开始

### 环境要求

- Go 1.23或更高版本
- PostgreSQL数据库服务器

### 安装

1. 克隆仓库：
   ```bash
   git clone https://github.com/leixiaotian1/pgsql-mcp-server.git
   cd sql-mcp-server
   ```

2. 安装依赖：
   ```bash
   go mod download
   ```

3. 构建mcp server：
   ```bash
   go build -o sql-mcp-server
   ```

## 配置

pg-mcp-server需要通过环境变量提供数据库连接详情。在项目根目录创建`.env`文件，包含以下变量：

```
DB_HOST=localhost      # PostgreSQL服务器主机
DB_PORT=5432           # PostgreSQL服务器端口
DB_NAME=postgres       # 数据库名称
DB_USER=your_username  # 数据库用户
DB_PASSWORD=your_pass  # 数据库密码
DB_SSLMODE=disable     # SSL模式(disable, require, verify-ca, verify-full)
SERVER_MODE=stdio      # 通信方式(stdio, sse, streamableHttp)
```

## 使用

### 运行服务器

```bash
./sql-mcp-server
```

### MCP配置

要与支持MCP的AI助手一起使用此 pg-mcp-server ，请将以下内容添加到您的MCP配置中：

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
### DOCKER 部署
<details>
<summary><strong>点击展开 Docker 部署指南</strong></summary>


#### 先决条件

- 已安装 Docker

#### 部署步骤

1.  **克隆项目**
    ```bash
    git clone https://github.com/leixiaotian1/pgsql-mcp-server.git
    cd pgsql-mcp-server
    ```

2.  **配置 `.env` 文件**
    
    在项目根目录创建一个 `.env` 文件。此文件用于存放数据库连接信息。**请确保 `DB_HOST` 的值与后续启动的数据库容器名称一致。**
    
    ```properties
    DB_HOST=postgres
    DB_PORT=5432
    DB_NAME=postgres
    DB_USER=user
    DB_PASSWORD=password
    DB_SSLMODE=disable
    SERVER_MODE=sse
    ```

3.  **创建 Docker 网络**

    为了让应用容器和数据库容器能够相互通信，我们需要创建一个共享的 Docker 网络。此命令只需执行一次。
    ```bash
    docker network create sql-mcp-network
    ```

4.  **启动 PostgreSQL 数据库容器**

    使用以下命令启动一个 PostgreSQL 容器，并将其连接到我们刚创建的网络。
    
    > **注意:**
    > - `--name postgres-dbpsk`：容器的名称，必须与 `.env` 文件中的 `DB_HOST` 完全匹配。
    > - `--network sql-mcp-network`：连接到共享网络。
    > - `-p 5432:5432`：将主机的 `5432` 端口映射到容器的 `5432` 端口。这意味着您可以从您的电脑（例如使用 DBeaver）通过 `localhost:5432` 连接数据库，而应用容器将通过内部网络直接访问 `5432` 端口。

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

5.  **构建并运行应用**

    现在，您可以使用 `Makefile` 中的命令来管理应用。

    - **构建镜像并运行容器:**
      ```bash
      make build
      make run
      ```
      此命令会自动停止旧容器、构建新镜像并启动新容器。

    - **查看应用日志:**
      ```bash
      make logs
      ```
      如果看到 `Successfully connected to database`，说明一切正常。

    - **停止应用:**
      ```bash
      make stop
      ```

</details>

---

## 🔌 服务器模式

您可以通过设置 `SERVER_MODE` 变量来选择传输协议。

### `stdio`

服务器通过标准输入和输出进行通信。这是默认模式，非常适合本地测试或与基于命令行的 MCP 客户端直接集成。

### `sse`

服务器使用服务器发送事件 (SSE) 进行通信。启用此模式后，服务器将启动一个 HTTP 服务并监听连接。

*   **SSE 端点:** `http://localhost:8088/sse`
*   **消息端点:** `http://localhost:8088/message`

### `streamableHttp`

服务器使用 Streamable HTTP 传输，这是一种更现代、更灵活的基于 HTTP 的 MCP 传输方式。

*   **端点:** `http://localhost:8088/mcp`

## 🤝 贡献

欢迎任何形式的贡献！如果您发现任何错误、有功能请求或改进建议，请随时提交拉取请求 (Pull Request) 或开启一个问题 (Issue)。

1.  Fork 本仓库。
2.  创建您的功能分支 (`git checkout -b feature/AmazingFeature`)。
3.  提交您的更改 (`git commit -m 'Add some AmazingFeature'`)。
4.  将分支推送到远程 (`git push origin feature/AmazingFeature`)。
5.  开启一个拉取请求。

## 📄 许可证

本项目是开源的，基于 MIT 许可证发布。

