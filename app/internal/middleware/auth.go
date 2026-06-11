package middleware

import (
	"autoadmin/internal/auth"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string

const MasterIDKey contextKey = "master_id"
const TelegramIDKey contextKey = "telegram_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tgID := r.Header.Get("X-Telegram-ID")

		if tgID != "" {
			var tid int
			if _, err := fmt.Sscanf(tgID, "%d", &tid); err == nil {
				ctx := context.WithValue(r.Context(), TelegramIDKey, tid)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		if authHeader == "" {
			http.Error(w, `{"error":"missing authorization"}`, http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateToken(token)
		if err != nil {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), MasterIDKey, claims.MasterID)
		ctx = context.WithValue(ctx, TelegramIDKey, claims.TelegramID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetMasterID(r *http.Request) (int, bool) {
	id, ok := r.Context().Value(MasterIDKey).(int)
	return id, ok
}

func GetTelegramID(r *http.Request) (int, bool) {
	id, ok := r.Context().Value(TelegramIDKey).(int)
	return id, ok
}
