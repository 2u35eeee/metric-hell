package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"metric-hell/pkg/game"
)

func TestHandlerServesFromFSWithoutProjectRoot(t *testing.T) {
	handler := NewHandlerFS([]game.Node{
		{ID: game.InitialNodeID, Title: "高考成绩 Benchmark", Stage: "高中"},
	}, fstest.MapFS{
		"index.html": {Data: []byte("<!doctype html><title>WorkflowBench</title>")},
		"style.css":  {Data: []byte("body{background:#111}")},
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET / status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "WorkflowBench") {
		t.Fatalf("GET / body = %q, want embedded index content", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodPost, "/api/new", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("POST /api/new status = %d, want 200", rec.Code)
	}
	var result game.Result
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode /api/new response: %v", err)
	}
	if result.CurrentNode == nil || result.CurrentNode.ID != game.InitialNodeID {
		t.Fatalf("current node = %#v, want %s", result.CurrentNode, game.InitialNodeID)
	}
}
