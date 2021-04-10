package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message" : "get method called"}`))

	case "POST":
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message" : "post method called}`))

	case "PUT":
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"message" : "put method called}`))

	case "DELETE":
		w.WriteHeader((http.StatusOK))
		w.Write([]byte(`{"message" : "delete method called}`))

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message" : "method not found}`))
	}
}

func main() {
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
