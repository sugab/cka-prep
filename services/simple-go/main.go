package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Message string `json:"message,omitempty"`
	Hash    string `json:"hash,omitempty"`
}

func main() {
	// get default message from env
	message := os.Getenv("DEFAULT_MESSAGE")
	if message == "" {
		// fallback message
		message = "Success"
	}

	srv := http.NewServeMux()
	srv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hasher := sha256.New()
		hasher.Write([]byte(message))
		hash := fmt.Sprintf("%x", hasher.Sum(nil))

		res := Response{
			Message: message,
			Hash:    hash,
		}

		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Fatal(err)
		}
	})

	fmt.Println("Server is running on port 8080...")
	err := http.ListenAndServe(":8080", srv)
	if err != nil {
		log.Fatal(err)
	}
}
