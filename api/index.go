package handler

import (
	"net/http"

	internalapi "metric-hell/internal/api"
)

var h = internalapi.MustNewDefaultHandler()

func Handler(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTP(w, r)
}
