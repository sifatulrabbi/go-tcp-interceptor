package main

import (
	"fmt"
	"log"
	"net/http"

	"helloscribe.ai/shadow-tracker/internals/api"
)

func main() {
	r := api.NewApi()
	fmt.Println("starting server on port 9001")
	if err := http.ListenAndServe(":9001", r); err != nil {
		log.Fatalln(err)
	}
}
