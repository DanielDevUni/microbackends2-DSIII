package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type Word struct {
	ID   int    `json:"id"`
	Word string `json:"word"`
}

// Conexión a Supabase usando variables de entorno
func GetDB() (*pgx.Conn, error) {
	url := os.Getenv("SUPABASE_HOST")  // ej: db.abcd.supabase.co
	user := os.Getenv("SUPABASE_USER") // normalmente "postgres"
	pass := os.Getenv("SUPABASE_PASSWORD")
	dbname := os.Getenv("SUPABASE_DB") // normalmente "postgres"
	port := os.Getenv("SUPABASE_PORT") // "5432"

	connStr := "postgres://" + user + ":" + pass + "@" + url + ":" + port + "/" + dbname
	return pgx.Connect(context.Background(), connStr)
}

// Esta es la función que Vercel detecta
func Handler(w http.ResponseWriter, r *http.Request) {
	db, err := GetDB()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer db.Close(context.Background())

	switch r.Method {
	case "GET":
		rows, _ := db.Query(context.Background(), "SELECT id, word FROM words")
		defer rows.Close()
		var words []Word
		for rows.Next() {
			var w Word
			rows.Scan(&w.ID, &w.Word)
			words = append(words, w)
		}
		json.NewEncoder(w).Encode(words)

	case "POST":
		var w Word
		json.NewDecoder(r.Body).Decode(&w)
		_, err := db.Exec(context.Background(), "INSERT INTO words (word) VALUES ($1)", w.Word)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Palabra creada"})

	case "PUT":
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		var w Word
		json.NewDecoder(r.Body).Decode(&w)
		_, err := db.Exec(context.Background(), "UPDATE words SET word=$1 WHERE id=$2", w.Word, id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Palabra actualizada"})

	case "DELETE":
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		_, err := db.Exec(context.Background(), "DELETE FROM words WHERE id=$1", id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Palabra eliminada"})
	}
}
