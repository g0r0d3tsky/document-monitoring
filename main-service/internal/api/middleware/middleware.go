package middleware

import (
	"context"
	"log/slog"
	"main-service/internal/domain"
	"net/http"
	"strconv"
	"strings"
)

type ContextKey string

const (
	UserCtx ContextKey = "user_info"
)

type UserService interface {
	GenerateToken(ctx context.Context, login string, password string) (string, error)
	ParseToken(token string) (*usecase.UserInfo, error)
}

type UserMiddleware struct {
	service UserService
}

func NewUserMiddleware(service UserService) *UserMiddleware {
	return &UserMiddleware{service: service}
}

func (m *UserMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			handlers.NewErrorResponse(w, http.StatusUnauthorized, "empty auth header")
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			handlers.NewErrorResponse(w, http.StatusUnauthorized, "invalid auth header")
			return
		}

		if len(headerParts[1]) == 0 || len(strings.Split(headerParts[1], "")) == 0 {
			handlers.NewErrorResponse(w, http.StatusUnauthorized, "token is empty")
			return
		}

		userInfo, err := m.service.ParseToken(headerParts[1])

		if err != nil {
			handlers.NewErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), UserCtx, userInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *UserMiddleware) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(UserCtx).(*usecase.UserInfo)
		if !ok {
			slog.Error("User context not found")
			handlers.NewErrorResponse(w, http.StatusUnauthorized, "User context not found")
			return
		}

		if user == nil {
			slog.Error("User information is nil")
			handlers.NewErrorResponse(w, http.StatusUnauthorized, "User information not found")
			return
		}
		if user.Role != domain.ADMIN {
			slog.Error("Access denied for non-admin user")
			handlers.NewErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *UserMiddleware) LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		mw := NewResponseWriter(w)
		next.ServeHTTP(mw, r)
		statusCode := mw.StatusCode()
		responseSize := mw.Size()

		slog.Info(
			"request ", "Method", r.Method, "Path", r.URL.Path, "Status", strconv.Itoa(statusCode), "Size", responseSize,
		)
	})
}

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK, 0}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriter) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	rw.size += size
	return size, err
}

func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

func (rw *ResponseWriter) Size() int {
	return rw.size
}
