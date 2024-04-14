package api

import (
	"fmt"
	"net/http"

	"helloscribe.ai/shadow-tracker/internals/logger"
)

func NewApi() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header)

		logger.NewHTTPLog(r)
		w.WriteHeader(http.StatusTeapot)
	})
	return mux
}
