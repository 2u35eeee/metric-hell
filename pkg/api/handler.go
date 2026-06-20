package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"metric-hell/pkg/content"
	"metric-hell/pkg/game"
)

type Handler struct {
	engine *game.Engine
	nodes  []game.Node
	webFS  fs.FS
}

type ActionRequest struct {
	State      game.State       `json:"state"`
	Submission *game.Submission `json:"submission"`
}

func NewHandler(nodes []game.Node, webDir string) *Handler {
	return NewHandlerFS(nodes, os.DirFS(webDir))
}

func NewHandlerFS(nodes []game.Node, webFS fs.FS) *Handler {
	return &Handler{
		engine: game.NewEngine(nodes),
		nodes:  nodes,
		webFS:  webFS,
	}
}

func MustNewDefaultHandler() http.Handler {
	root, err := FindProjectRoot()
	if err != nil {
		return errorHandler{err: fmt.Errorf("find project root: %w", err)}
	}
	nodes, err := content.LoadNodes(filepath.Join(root, "data", "nodes.json"))
	if err != nil {
		return errorHandler{err: err}
	}
	return NewHandler(nodes, filepath.Join(root, "web"))
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/api/new" && r.Method == http.MethodPost:
		h.handleNew(w, r)
	case r.URL.Path == "/api/action" && r.Method == http.MethodPost:
		h.handleAction(w, r)
	case r.URL.Path == "/api/nodes" && r.Method == http.MethodGet:
		writeJSON(w, http.StatusOK, h.nodes)
	case strings.HasPrefix(r.URL.Path, "/api/"):
		writeError(w, http.StatusNotFound, "api route not found")
	default:
		h.serveStatic(w, r)
	}
}

func (h *Handler) handleNew(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.engine.InitialResult(game.NewRandomSeed()))
}

func (h *Handler) handleAction(w http.ResponseWriter, r *http.Request) {
	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request")
		return
	}
	if req.Submission == nil {
		writeError(w, http.StatusBadRequest, "missing submission")
		return
	}
	result, err := h.engine.StepSubmission(req.State, *req.Submission)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) serveStatic(w http.ResponseWriter, r *http.Request) {
	path := filepath.Clean(r.URL.Path)
	if path == "." || path == "/" {
		path = "index.html"
	}
	path = strings.TrimPrefix(path, "/")
	if strings.Contains(path, "..") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	if !fs.ValidPath(path) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	data, err := fs.ReadFile(h.webFS, path)
	if errors.Is(err, fs.ErrNotExist) {
		path = "index.html"
		data, err = fs.ReadFile(h.webFS, path)
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("read static asset: %v", err))
		return
	}
	if contentType := mime.TypeByExtension(filepath.Ext(path)); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func FindProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	candidates := []string{
		wd,
		filepath.Dir(wd),
		filepath.Dir(filepath.Dir(wd)),
	}
	for _, candidate := range candidates {
		if fileExists(filepath.Join(candidate, "data", "nodes.json")) && fileExists(filepath.Join(candidate, "web", "index.html")) {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("could not find data/nodes.json and web/index.html from %s", wd)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

type errorHandler struct {
	err error
}

func (e errorHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	writeError(w, http.StatusInternalServerError, e.err.Error())
}
