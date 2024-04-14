package logger

import (
	"fmt"
	"net/http"
)

func NewHTTPLog(r *http.Request) {
	constructFilename(r)
}

func constructFilename(r *http.Request) string {
	return fmt.Sprintf("%s", r.URL.String())
}
