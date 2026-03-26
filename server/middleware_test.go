package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func setupTestLogger(t *testing.T) {
	t.Helper()
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)
}

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
})

// --- statusWriter tests ---

func TestStatusWriter_CapturesStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	sw := &statusWriter{ResponseWriter: rec, status: http.StatusOK}

	sw.WriteHeader(http.StatusNotFound)

	assert.Equal(t, http.StatusNotFound, sw.status)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestStatusWriter_DefaultStatus(t *testing.T) {
	sw := statusWriter{status: http.StatusOK}

	assert.Equal(t, http.StatusOK, sw.status)
}

type fakeFlusher struct {
	http.ResponseWriter
	flushed bool
}

func (f *fakeFlusher) Flush() { f.flushed = true }

func TestStatusWriter_Flush(t *testing.T) {
	ff := &fakeFlusher{ResponseWriter: httptest.NewRecorder()}
	sw := &statusWriter{ResponseWriter: ff, status: http.StatusOK}

	sw.Flush()

	assert.True(t, ff.flushed)
}

// --- peekJSONRPCRequest tests ---

func TestPeekJSONRPCRequest_NilBody(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Body = nil

	info := peekJSONRPCRequest(r)

	assert.Empty(t, info.Method)
	assert.Empty(t, info.ToolName)
	assert.Empty(t, info.Params)
}

func TestPeekJSONRPCRequest_EmptyBody(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))

	info := peekJSONRPCRequest(r)

	assert.Empty(t, info.Method)
}

func TestPeekJSONRPCRequest_InvalidJSON(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{invalid"))

	info := peekJSONRPCRequest(r)

	assert.Empty(t, info.Method)
}

func TestPeekJSONRPCRequest_ValidRequest(t *testing.T) {
	body := `{"method":"initialize","params":{"name":"test"}}`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	info := peekJSONRPCRequest(r)

	assert.Equal(t, "initialize", info.Method)
	assert.Empty(t, info.ToolName)
	assert.NotEmpty(t, info.Params)
}

func TestPeekJSONRPCRequest_ToolsCall(t *testing.T) {
	body := `{"method":"tools/call","params":{"name":"greeting_hello"}}`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	info := peekJSONRPCRequest(r)

	assert.Equal(t, "tools/call", info.Method)
	assert.Equal(t, "greeting_hello", info.ToolName)
}

func TestPeekJSONRPCRequest_BodyRestored(t *testing.T) {
	body := `{"method":"initialize","params":{}}`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	_ = peekJSONRPCRequest(r)

	restored, err := io.ReadAll(r.Body)
	require.NoError(t, err)
	assert.Equal(t, body, string(restored))
}

// --- withAuthMiddleware tests ---

func TestAuthMiddleware_NoToken(t *testing.T) {
	handler := withAuthMiddleware(okHandler, "")
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	handler := withAuthMiddleware(okHandler, "secret123")
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer secret123")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	handler := withAuthMiddleware(okHandler, "secret123")

	tests := []struct {
		name   string
		header string
	}{
		{"wrong token", "Bearer wrong"},
		{"no header", ""},
		{"not bearer", "Basic secret123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				r.Header.Set("Authorization", tt.header)
			}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

// --- withOriginValidation tests ---

func TestOriginValidation_NoAllowedOrigins(t *testing.T) {
	handler := withOriginValidation(okHandler, nil)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Origin", "https://evil.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOriginValidation_NoOriginHeader(t *testing.T) {
	handler := withOriginValidation(okHandler, []string{"https://allowed.com"})
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOriginValidation_AllowedOrigin(t *testing.T) {
	allowed := []string{"https://allowed.com", "https://also-ok.com"}
	handler := withOriginValidation(okHandler, allowed)

	for _, origin := range allowed {
		t.Run(origin, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Origin", origin)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestOriginValidation_ForbiddenOrigin(t *testing.T) {
	handler := withOriginValidation(okHandler, []string{"https://allowed.com"})
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Origin", "https://evil.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// --- withRequestLogging tests ---

func TestRequestLogging_PassThrough(t *testing.T) {
	setupTestLogger(t)

	handler := withRequestLogging(okHandler)
	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}

func TestRequestLogging_CapturesStatus(t *testing.T) {
	setupTestLogger(t)

	errHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	handler := withRequestLogging(errHandler)
	r := httptest.NewRequest(http.MethodPost, "/rpc", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
