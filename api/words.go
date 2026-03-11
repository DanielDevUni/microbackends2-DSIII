package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Word struct {
	ID   int    `json:"id"`
	Word string `json:"word"`
}

func connectDB(connStr string) (*pgx.Conn, error) {
	config, err := pgx.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	// Supabase pooler uses pgBouncer semantics; simple protocol avoids prepared statement cache issues.
	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	return pgx.ConnectConfig(context.Background(), config)
}

func GetDB() (*pgx.Conn, error) {
	if connStr := strings.TrimSpace(os.Getenv("SUPABASE_URL")); connStr != "" {
		return connectDB(connStr)
	}

	host := strings.TrimSpace(os.Getenv("SUPABASE_HOST"))
	user := strings.TrimSpace(os.Getenv("SUPABASE_USER"))
	pass := strings.TrimSpace(os.Getenv("SUPABASE_PASSWORD"))
	dbname := strings.TrimSpace(os.Getenv("SUPABASE_DB"))
	port := strings.TrimSpace(os.Getenv("SUPABASE_PORT"))

	missing := make([]string, 0, 5)
	if host == "" {
		missing = append(missing, "SUPABASE_HOST")
	}
	if user == "" {
		missing = append(missing, "SUPABASE_USER")
	}
	if pass == "" {
		missing = append(missing, "SUPABASE_PASSWORD")
	}
	if dbname == "" {
		missing = append(missing, "SUPABASE_DB")
	}
	if port == "" {
		missing = append(missing, "SUPABASE_PORT")
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing database configuration: define SUPABASE_URL or %s", strings.Join(missing, ", "))
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", user, pass, host, port, dbname)
	return connectDB(connStr)
}

func decodeWord(r *http.Request) (Word, error) {
	defer r.Body.Close()

	var input Word
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		if errors.Is(err, io.EOF) {
			return Word{}, errors.New("request body is required")
		}
		return Word{}, err
	}

	if strings.TrimSpace(input.Word) == "" {
		return Word{}, errors.New("field 'word' is required")
	}

	return input, nil
}

// Esta es la funcion que Vercel detecta.
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Debug-Method", r.Method)
	w.Header().Set("X-Debug-Path", r.URL.Path)

	db, err := GetDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close(context.Background())

	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query(context.Background(), "SELECT id, word FROM words")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var words []Word
		for rows.Next() {
			var word Word
			if err := rows.Scan(&word.ID, &word.Word); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			words = append(words, word)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(words)

	case http.MethodPost:
		input, err := decodeWord(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = db.Exec(context.Background(), "INSERT INTO words (word) VALUES ($1)", input.Word)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Palabra creada"})

	case http.MethodPut:
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		input, err := decodeWord(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = db.Exec(context.Background(), "UPDATE words SET word=$1 WHERE id=$2", input.Word, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Palabra actualizada"})

	case http.MethodDelete:
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))

		_, err := db.Exec(context.Background(), "DELETE FROM words WHERE id=$1", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Palabra eliminada"})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
