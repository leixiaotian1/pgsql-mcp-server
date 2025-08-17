package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mark3labs/mcp-go/server"
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

	// Register tools for both server types
	s.AddTool(readQueryTool(), readQueryToolHandler)
	s.AddTool(writeQueryTool(), writeQueryToolHandler)
	s.AddTool(createTableTool(), createTableToolHandler)
	s.AddTool(listTablesTool(), listTableToolHandler)
	s.AddTool(explainQueryTool(), explainQueryToolHandler)

	serverMode := os.Getenv("SERVER_MODE")
	if serverMode == "" {
		serverMode = "stdio" // Default to stdio mode
	}

	switch serverMode {
	case "stdio":
		if err := server.ServeStdio(s); err != nil {
			log.Printf("Stdio server error: %v\n", err)
		}
	case "sse":
		sse := server.NewSSEServer(s)
		log.Printf("sse server listening on :8088/sse")
		if err := sse.Start(":8088"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	case "streamableHttp":
		sse := server.NewStreamableHTTPServer(s)
		log.Printf("streamableHttp server listening on :8088/mcp")
		if err := sse.Start(":8088"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	default:
		log.Fatalf("Unknown SERVER_MODE: %s. Use 'stdio' or 'http'.", serverMode)
	}
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
