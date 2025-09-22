package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Message struct {
	Text string `json:"message"`
}

func main() {
	http.HandleFunc("/button1", func(w http.ResponseWriter, r *http.Request) {
		msg := Message{Text: "Button 1 Clicked"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)
	})

	http.HandleFunc("/button2", func(w http.ResponseWriter, r *http.Request) {
		msg := Message{Text: "Button 2 Clicked"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)
	})

	log.Println("Server Started!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}