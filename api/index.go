package handler

import (
	"embed"
	"io/fs"
	"net/http"

	internalapi "metric-hell/pkg/api"
	"metric-hell/pkg/content"
)

//go:embed data/nodes.json web/*
var embeddedFiles embed.FS

var h = newEmbeddedHandler()

func Handler(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTP(w, r)
}

func newEmbeddedHandler() http.Handler {
	nodes, err := content.LoadNodesFS(embeddedFiles, "data/nodes.json")
	if err != nil {
		return embeddedErrorHandler{err: err}
	}
	webFS, err := fs.Sub(embeddedFiles, "web")
	if err != nil {
		return embeddedErrorHandler{err: err}
	}
	return internalapi.NewHandlerFS(nodes, webFS)
}

type embeddedErrorHandler struct {
	err error
}

func (h embeddedErrorHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(`{"error":"` + h.err.Error() + `"}`))
}
