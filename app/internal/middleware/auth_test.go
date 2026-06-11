package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"autoadmin/internal/auth"
)

func TestAuthMiddleware_JWT(t *testing.T) {
	token, err := auth.GenerateToken(1, 12345)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		masterID, ok := GetMasterID(r)
		if !ok {
			t.Error("master_id not in context")
			return
		}
		if masterID != 1 {
			t.Errorf("expected master_id=1, got %d", masterID)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAuthMiddleware_TelegramID(t *testing.T) {
	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tgID, ok := GetTelegramID(r)
		if !ok {
			t.Error("telegram_id not in context")
			return
		}
		if tgID != 99999 {
			t.Errorf("expected telegram_id=99999, got %d", tgID)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Telegram-ID", "99999")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAuthMiddleware_NoAuth(t *testing.T) {
	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_BearerFormatRequired(t *testing.T) {
	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Token abc123")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetMasterID_NotSet(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	_, ok := GetMasterID(req)
	if ok {
		t.Error("expected ok=false when master_id not set")
	}
}

func TestGetTelegramID_NotSet(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	_, ok := GetTelegramID(req)
	if ok {
		t.Error("expected ok=false when telegram_id not set")
	}
}
