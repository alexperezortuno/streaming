package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexperezortuno/streaming/internal/model"
)

func TestRespondJSON(t *testing.T) {
	w := httptest.NewRecorder()
	respondJSON(w, http.StatusOK, map[string]string{"key": "value"})

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", resp.StatusCode)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["key"] != "value" {
		t.Fatalf("expected value, got: %s", body["key"])
	}
}

func TestRespondJSON_NilData(t *testing.T) {
	w := httptest.NewRecorder()
	respondJSON(w, http.StatusNoContent, nil)

	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got: %d", resp.StatusCode)
	}
}

func TestRespondError_NotFound(t *testing.T) {
	w := httptest.NewRecorder()
	respondError(w, model.ErrNotFound)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got: %d", resp.StatusCode)
	}

	var body ErrorResponse
	json.NewDecoder(resp.Body).Decode(&body)
	if body.Error != "resource not found" {
		t.Fatalf("expected 'resource not found', got: %s", body.Error)
	}
}

func TestRespondError_Unauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	respondError(w, model.ErrUnauthorized)

	resp := w.Result()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got: %d", resp.StatusCode)
	}
}

func TestRespondError_InvalidInput(t *testing.T) {
	w := httptest.NewRecorder()
	respondError(w, model.ErrInvalidInput)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got: %d", resp.StatusCode)
	}
}

func TestRespondError_Duplicate(t *testing.T) {
	w := httptest.NewRecorder()
	respondError(w, model.ErrDuplicate)

	resp := w.Result()
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected status 409, got: %d", resp.StatusCode)
	}
}

func TestRespondError_Internal(t *testing.T) {
	w := httptest.NewRecorder()
	respondError(w, model.ErrInternal)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got: %d", resp.StatusCode)
	}
}

func TestRespondError_Unknown(t *testing.T) {
	w := httptest.NewRecorder()
	respondError(w, model.ErrForbidden)

	resp := w.Result()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected status 403, got: %d", resp.StatusCode)
	}
}
