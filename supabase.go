package main

import (
    "context"
    "os"

    "github.com/jackc/pgx/v5"
)

func GetDB() (*pgx.Conn, error) {
    url := os.Getenv("SUPABASE_URL")     // ej: db.abcd.supabase.co
    user := os.Getenv("SUPABASE_USER")   // postgres
    pass := os.Getenv("SUPABASE_PASSWORD")
    dbname := os.Getenv("SUPABASE_DB")   // postgres
    port := os.Getenv("SUPABASE_PORT")   // 5432

    connStr := "postgres://" + user + ":" + pass + "@" + url + ":" + port + "/" + dbname
    return pgx.Connect(context.Background(), connStr)
}
