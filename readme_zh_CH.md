# PostgreSQL MCP 服务器
![Stars](https://img.shields.io/github/stars/leixiaotian1/pgsql-mcp-server)
![Forks](https://img.shields.io/github/forks/leixiaotian1/pgsql-mcp-server)

中文 | [English](README.md)

一个提供PostgreSQL数据库交互工具的模型上下文协议(MCP)服务。该服务使AI助手能够通过MCP协议执行SQL查询、解释sql语句、创建表和列出数据库表。

## 功能

服务器提供以下工具：

- **read_query**: 在PostgreSQL数据库上执行SELECT查询
- **write_query**: 在PostgreSQL数据库上执行INSERT、UPDATE或DELETE查询
- **create_table**: 在PostgreSQL数据库中创建新表
- **list_tables**: 列出数据库中的所有用户表(可选模式过滤)
- **explain_query**: 解释PostgreSQL数据库上执行的SQL语句

## 安装

### 先决条件

- Go 1.23或更高版本
- PostgreSQL数据库服务器

### 步骤

1. 克隆仓库：
   ```bash
   git clone https://github.com/sql-mcp-server.git
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
        "DB_SSLMODE": "disable"
      },
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

### 工具示例

#### 列出表

列出数据库中的所有用户表：

```json
{
  "server_name": "pgsql-mcp-server",
  "tool_name": "list_tables",
  "arguments": {}
}
```

列出特定模式中的表：

```json
{
  "server_name": "pgsql-mcp-server",
  "tool_name": "list_tables",
  "arguments": {
    "schema": "public"
  }
}
```

#### 创建表

创建新表：

```json
{
  "server_name": "pgsql-mcp-server",
  "tool_name": "create_table",
  "arguments": {
    "schema": "CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(100), email VARCHAR(100) UNIQUE, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)"
  }
}
```

#### 读取查询

执行SELECT查询：

```json
{
  "server_name": "pgsql-mcp-server",
  "tool_name": "read_query",
  "arguments": {
    "query": "SELECT * FROM users LIMIT 10"
  }
}
```

#### 写入查询

执行INSERT查询：

```json
{
  "server_name": "pgsql-mcp-server",
  "tool_name": "write_query",
  "arguments": {
    "query": "INSERT INTO users (name, email) VALUES ('John Doe', 'john@example.com')"
  }
}
```

执行UPDATE查询：

```json
{
  "server_name": "pgsql-mcp-server",
  "tool_name": "write_query",
  "arguments": {
    "query": "UPDATE users SET name = 'Jane Doe' WHERE id = 1"
  }
}
```

执行DELETE查询：

```json
{
  "server_name": "pgsql-mcp-server",
  "tool_name": "write_query",
  "arguments": {
    "query": "DELETE FROM users WHERE id = 1"
  }
}
```

## 安全考虑

- pg-mcp-server 验证查询类型以确保每个工具只执行适当的操作
- 对模式名称执行输入清理以防止SQL注入
- 考虑为此 pg-mcp-server 使用具有有限权限的专用数据库用户
- 在生产环境中，通过将`DB_SSLMODE`设置为`require`或更高来启用SSL

## 依赖项

- [github.com/joho/godotenv](https://github.com/joho/godotenv) - 用于从.env文件加载环境变量
- [github.com/lib/pq](https://github.com/lib/pq) - Go的PostgreSQL驱动
- [github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) - 模型上下文协议的Go SDK

