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

func GetDB() (*pgx.Conn, error) {
	url := os.Getenv("SUPABASE_HOST")
	user := os.Getenv("SUPABASE_USER")
	pass := os.Getenv("SUPABASE_PASSWORD")
	dbname := os.Getenv("SUPABASE_DB")
	port := os.Getenv("SUPABASE_PORT")

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
			var word Word
			rows.Scan(&word.ID, &word.Word)
			words = append(words, word)
		}
		json.NewEncoder(w).Encode(words)

	case "POST":
		var input Word
		json.NewDecoder(r.Body).Decode(&input)
		_, err := db.Exec(context.Background(), "INSERT INTO words (word) VALUES ($1)", input.Word)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Palabra creada"})

	case "PUT":
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		var input Word
		json.NewDecoder(r.Body).Decode(&input)
		_, err := db.Exec(context.Background(), "UPDATE words SET word=$1 WHERE id=$2", input.Word, id)
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
