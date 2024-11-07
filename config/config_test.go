package config

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestConnectDatabase(t *testing.T) {
	err := godotenv.Load("/workspaces/Invoicerator/.env.test")
	if err != nil {
		t.Fatalf("Failed to load .env.test: %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Fatal("DATABASE_URL is not set in the environment")
	}

	ConnectDatabase()

	if DB == nil {
		t.Fatal("Database connection failed: DB object is nil")
	}
}
