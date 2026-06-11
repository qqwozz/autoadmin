package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestRespondJSON(t *testing.T) {
	w := httptest.NewRecorder()
	respondJSON(w, 200, map[string]string{"key": "value"})

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["key"] != "value" {
		t.Errorf("expected key=value, got %v", result)
	}
}

func TestRespondError(t *testing.T) {
	w := httptest.NewRecorder()
	respondError(w, 404, "not found")

	if w.Code != 404 {
		t.Errorf("expected 404, got %d", w.Code)
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["error"] != "not found" {
		t.Errorf("expected error=not found, got %v", result)
	}
}

func TestGetID_Valid(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.SetPathValue("id", "42")

	id, ok := getID(req, "id")
	if !ok {
		t.Error("expected ok=true")
	}
	if id != 42 {
		t.Errorf("expected 42, got %d", id)
	}
}

func TestGetID_Invalid(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.SetPathValue("id", "abc")

	_, ok := getID(req, "id")
	if ok {
		t.Error("expected ok=false for non-numeric id")
	}
}

func TestDecodeJSON(t *testing.T) {
	body := bytes.NewBufferString(`{"name":"test","value":123}`)
	req := httptest.NewRequest("POST", "/", body)

	var result struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	err := decodeJSON(req, &result)
	if err != nil {
		t.Fatalf("decodeJSON failed: %v", err)
	}

	if result.Name != "test" {
		t.Errorf("expected name=test, got %s", result.Name)
	}
	if result.Value != 123 {
		t.Errorf("expected value=123, got %d", result.Value)
	}
}

func TestDecodeJSON_Invalid(t *testing.T) {
	body := bytes.NewBufferString(`not json`)
	req := httptest.NewRequest("POST", "/", body)

	var result struct {
		Name string `json:"name"`
	}
	err := decodeJSON(req, &result)
	if err == nil {
		t.Error("expected error for invalid json")
	}
}

func TestDerefStr(t *testing.T) {
	s := "hello"
	if derefStr(&s) != "hello" {
		t.Error("expected hello")
	}
	if derefStr(nil) != "" {
		t.Error("expected empty string for nil")
	}
}

func TestParseOptionalInt(t *testing.T) {
	v, err := parseOptionalInt("123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 123 {
		t.Errorf("expected 123, got %d", v)
	}

	_, err = parseOptionalInt("abc")
	if err == nil {
		t.Error("expected error for non-numeric string")
	}
}
