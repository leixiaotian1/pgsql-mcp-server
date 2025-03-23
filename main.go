package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var (
	db         *sql.DB
	writeStmt  = regexp.MustCompile(`(?i)^\s*(INSERT|UPDATE|DELETE)`)
	selectStmt = regexp.MustCompile(`(?i)^\s*SELECT`)
	createStmt = regexp.MustCompile(`(?i)^\s*CREATE TABLE`)
)

type DB struct {
	*sql.DB
}

func init() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file:", err)
		}
	}
}

func main() {
	// 初始化数据库连接池
	if err := initConnectionPool(); err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	s := server.NewMCPServer(
		"sql-mcp-server 🚀",
		"1.0.0",
	)

	s.AddTool(createReadQueryTool(), readQueryToolHandler)
	s.AddTool(createWriteQueryTool(), writeQueryToolHandler)
	s.AddTool(createCreateTableTool(), createTableToolHandler)
	s.AddTool(createListTablesTool(), listTableToolHandler)

	if err := server.ServeStdio(s); err != nil {
		log.Printf("Server error: %v\n", err)
	}
}

func createReadQueryTool() mcp.Tool {
	return mcp.NewTool("read_query",
		mcp.WithDescription("Execute a SELECT query on the postgres database"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("SELECT SQL query to execute"),
		),
	)
}

func createWriteQueryTool() mcp.Tool {
	return mcp.NewTool("write_query",
		mcp.WithDescription("Execute an INSERT, UPDATE, or DELETE query on the postgres database"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("SQL query to execute"),
		),
	)
}

func createCreateTableTool() mcp.Tool {
	return mcp.NewTool("create_table",
		mcp.WithDescription("Create a new table in the postgres database"),
		mcp.WithString("schema",
			mcp.Required(),
			mcp.Description("CREATE TABLE SQL statement"),
		),
	)
}

func createListTablesTool() mcp.Tool {
	return mcp.NewTool("list_tables",
		mcp.WithDescription("List all user tables in the database"),
		mcp.WithString("schema",
			mcp.Description("Optional schema name to filter tables"),
		),
	)
}

func initConnectionPool() error {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return fmt.Errorf("invalid DB_PORT: %w", err)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		port,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	if err = db.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	log.Println("Successfully connected to database")
	return nil
}

func readQueryToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, ok := request.Params.Arguments["query"].(string)
	if !ok {
		return nil, errors.New("invalid query parameter")
	}

	if !selectStmt.MatchString(query) {
		return nil, errors.New("only SELECT queries are allowed")
	}

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Query error: %v\n", err)
		return nil, fmt.Errorf("query execution failed")
	}
	defer rows.Close()

	results, err := parseSQLRows(rows)
	if err != nil {
		return nil, fmt.Errorf("result parsing failed")
	}

	return mcp.NewToolResultText(fmt.Sprintf("Query results: %v", results)), nil
}

func writeQueryToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, ok := request.Params.Arguments["query"].(string)
	if !ok {
		return nil, errors.New("invalid query parameter")
	}

	if !writeStmt.MatchString(query) {
		return nil, errors.New("only INSERT/UPDATE/DELETE queries are allowed")
	}

	result, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Write error: %v\n", err)
		return nil, fmt.Errorf("write operation failed")
	}

	rowsAffected, _ := result.RowsAffected()
	return mcp.NewToolResultText(fmt.Sprintf("Operation successful. Rows affected: %d", rowsAffected)), nil
}

func createTableToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	schema, ok := request.Params.Arguments["schema"].(string)
	if !ok {
		return nil, errors.New("invalid schema parameter")
	}

	if !createStmt.MatchString(schema) {
		return nil, errors.New("invalid CREATE TABLE statement")
	}

	if _, err := db.ExecContext(ctx, schema); err != nil {
		log.Printf("Create table error: %v\n", err)
		return nil, fmt.Errorf("table creation failed")
	}

	return mcp.NewToolResultText("Table created successfully"), nil
}

func listTableToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	schemaFilter := ""
	if schema, ok := request.Params.Arguments["schema"].(string); ok {
		schemaFilter = fmt.Sprintf(" AND schemaname = '%s'", sanitizeInput(schema))
	}

	query := fmt.Sprintf(`
		SELECT tablename 
		FROM pg_catalog.pg_tables 
		WHERE schemaname NOT IN ('pg_catalog', 'information_schema') %s
	`, schemaFilter)

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("List tables error: %v\n", err)
		return nil, fmt.Errorf("failed to list tables")
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, fmt.Errorf("error scanning table name")
		}
		tables = append(tables, table)
	}

	return mcp.NewToolResultText(fmt.Sprintf("Tables: %v", tables)), nil
}

func parseSQLRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(cols))
		pointers := make([]interface{}, len(cols))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range cols {
			row[col] = values[i]
		}
		results = append(results, row)
	}
	return results, nil
}

func sanitizeInput(input string) string {
	return strings.ReplaceAll(input, "'", "''")
}
