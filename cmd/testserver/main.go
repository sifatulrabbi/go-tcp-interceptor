package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		res := map[string]string{
			"message": "Hello world",
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Fatalln(err)
		}
	})

	http.ListenAndServe(":9002", mux)
}
