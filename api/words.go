package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Word struct {
	ID   int    `json:"id"`
	Word string `json:"word"`
}

func handler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		id, _ := strconv.Atoi(ps.ByName("id"))
		var w Word
		json.NewDecoder(r.Body).Decode(&w)
		_, err := db.Exec(context.Background(), "UPDATE words SET word=$1 WHERE id=$2", w.Word, id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Palabra actualizada"})

	case "DELETE":
		id, _ := strconv.Atoi(ps.ByName("id"))
		_, err := db.Exec(context.Background(), "DELETE FROM words WHERE id=$1", id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Palabra eliminada"})
	}
}

func main() {
	router := httprouter.New()
	router.GET("/api/words", handler)
	router.POST("/api/words", handler)
	router.PUT("/api/words/:id", handler)
	router.DELETE("/api/words/:id", handler)

	http.ListenAndServe(":3000", router)
}
