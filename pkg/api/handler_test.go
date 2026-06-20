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

func TestHandlerActionAcceptsSubmissionAndReturnsAuditRecord(t *testing.T) {
	scoreMin := 700.0
	handler := NewHandlerFS([]game.Node{
		{
			ID:    game.InitialNodeID,
			Title: "高考成绩 Benchmark",
			Stage: "高中",
			Input: game.InputSpec{Type: game.InputTypeNumber, Prompt: "你的高考分数是多少？"},
			Options: []game.AnswerOption{
				{
					ID:      "score_700_plus",
					Label:   "700+",
					Min:     &scoreMin,
					Verdict: "你要是 700 以上还认识这个开发者？系统怀疑样本来源异常。",
					Proof:   "系统把你标记为罕见样本。",
					Effects: game.Effects{BenchScore: 20},
					Unlocks: []string{"university_tier"},
				},
			},
		},
		{
			ID:    "university_tier",
			Title: "大学层次过滤器",
			Stage: "大学",
		},
	}, fstest.MapFS{
		"index.html": {Data: []byte("<!doctype html><title>WorkflowBench</title>")},
	})

	initialReq := httptest.NewRequest(http.MethodPost, "/api/new", nil)
	initialRec := httptest.NewRecorder()
	handler.ServeHTTP(initialRec, initialReq)
	var initial game.Result
	if err := json.NewDecoder(initialRec.Body).Decode(&initial); err != nil {
		t.Fatalf("decode /api/new response: %v", err)
	}

	body := `{"state":` + mustJSON(t, initial.State) + `,"submission":{"node_id":"gaokao_score","numeric_value":701}}`
	req := httptest.NewRequest(http.MethodPost, "/api/action", strings.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("POST /api/action status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var result game.Result
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode /api/action response: %v", err)
	}
	if result.AuditRecord == nil {
		t.Fatal("audit_record = nil, want record")
	}
	if result.AuditRecord.SubmittedLabel != "701" {
		t.Fatalf("submitted label = %q, want 701", result.AuditRecord.SubmittedLabel)
	}
	if !strings.Contains(result.AuditRecord.Verdict, "700 以上") {
		t.Fatalf("verdict = %q, want 700+ copy", result.AuditRecord.Verdict)
	}
}

func mustJSON(t *testing.T, value any) string {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal value: %v", err)
	}
	return string(data)
}
