package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Message struct {
	ID      int    `json:"id" db:"id"`
	Content string `json:"content" db:"content"`
}

func main() {
	// Get user and password from env
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	// Database connection
	host := "postgres.default.svc.cluster.local"
	connectionStr := fmt.Sprintf("host=%s user=%s password=%s dbname=mydb sslmode=disable", host, user, password)
	db, err := sqlx.Connect("postgres", connectionStr)
	if err != nil {
		log.Fatal(err)
	}

	// handlers
	mux := mux.NewRouter()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// fetch data from database
		messages := []Message{}
		err := db.SelectContext(
			r.Context(),
			&messages,
			"SELECT * FROM message",
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		// return it as json response
		err = json.NewEncoder(w).Encode(messages)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}).Methods(http.MethodGet)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// get content from json body
		message := &Message{}
		err := json.NewDecoder(r.Body).Decode(message)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		// save to database
		_, err = db.NamedExecContext(
			r.Context(),
			"INSERT INTO message (content) VALUES (:content)",
			message,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("SUCCESS"))
	}).Methods(http.MethodPost)

	// start server
	fmt.Println("Server is running on port 8080...")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
