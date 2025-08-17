package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

var (
	db          *sql.DB
	writeStmt   = regexp.MustCompile(`(?i)^\s*(INSERT|UPDATE|DELETE)`)
	selectStmt  = regexp.MustCompile(`(?i)^\s*SELECT`)
	createStmt  = regexp.MustCompile(`(?i)^\s*CREATE TABLE`)
	explainStmt = regexp.MustCompile(`(?i)^\s*EXPLAIN`)
)

type ToolHandlerFunc func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)

func readQueryTool() mcp.Tool {
	return mcp.NewTool("read_query",
		mcp.WithDescription("Execute a SELECT query on the postgres database"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("SELECT SQL query to execute"),
		),
	)
}

func writeQueryTool() mcp.Tool {
	return mcp.NewTool("write_query",
		mcp.WithDescription("Execute an INSERT, UPDATE, or DELETE query on the postgres database"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("SQL query to execute"),
		),
	)
}

func createTableTool() mcp.Tool {
	return mcp.NewTool("create_table",
		mcp.WithDescription("Create a new table in the postgres database"),
		mcp.WithString("schema",
			mcp.Required(),
			mcp.Description("CREATE TABLE SQL statement"),
		),
	)
}

func listTablesTool() mcp.Tool {
	return mcp.NewTool("list_tables",
		mcp.WithDescription("List all user tables in the database"),
		mcp.WithString("schema",
			mcp.Description("Optional schema name to filter tables"),
		),
	)
}

// 创建explain查询工具定义
func explainQueryTool() mcp.Tool {
	return mcp.NewTool("explain_query",
		mcp.WithDescription("Explain a query execution plan on the postgres database"),
		mcp.WithString("schema",
			mcp.Required(),
			mcp.Description("SQL query to explain,start with EXPLAIN"),
		),
	)
}

func readQueryToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, ok := request.GetArguments()["query"].(string)
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
	query, ok := request.GetArguments()["query"].(string)
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
	schema, ok := request.GetArguments()["schema"].(string)
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
	if schema, ok := request.GetArguments()["schema"].(string); ok {
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

// explain查询处理函数
func explainQueryToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, ok := request.GetArguments()["schema"].(string)
	if !ok {
		return nil, errors.New("invalid schema parameter")
	}

	if !explainStmt.MatchString(query) {
		return nil, errors.New("invalid explain schema parameter")
	}

	// 执行EXPLAIN查询

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Explain error: %v\n", err)
		return nil, fmt.Errorf("explain execution failed: %v", err)
	}
	defer rows.Close()

	// 解析结果
	var plan strings.Builder
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			return nil, fmt.Errorf("error scanning explain result: %v", err)
		}
		plan.WriteString(line)
		plan.WriteString("\n")
	}

	return mcp.NewToolResultText(fmt.Sprintf("Execution plan:\n%s", plan.String())), nil
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
