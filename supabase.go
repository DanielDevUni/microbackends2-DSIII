package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func GetDB() (*pgx.Conn, error) {
	connStr := os.Getenv("SUPABASE_URL")
	return pgx.Connect(context.Background(), connStr)
}
