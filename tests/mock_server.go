package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
)

// MockResponse holds a canned response for a specific method+path.
type MockResponse struct {
	Status int
	Body   any
}

// MockServer is a lightweight HTTP mock for the Jellyfin API.
type MockServer struct {
	mu        sync.Mutex
	responses map[string]*MockResponse // "METHOD /path" -> response
	server    *httptest.Server
}

// NewMockServer creates and starts a mock Jellyfin API server.
func NewMockServer() *MockServer {
	ms := &MockServer{
		responses: make(map[string]*MockResponse),
	}
	ms.server = httptest.NewServer(http.HandlerFunc(ms.handler))
	return ms
}

// URL returns the base URL of the mock server.
func (ms *MockServer) URL() string {
	return ms.server.URL
}

// Close shuts down the mock server.
func (ms *MockServer) Close() {
	ms.server.Close()
}

// Reset clears all registered responses.
func (ms *MockServer) Reset() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.responses = make(map[string]*MockResponse)
}

// On registers a canned response for the given method and path.
func (ms *MockServer) On(method, path string, status int, body any) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	key := method + " " + path
	ms.responses[key] = &MockResponse{Status: status, Body: body}
}

func (ms *MockServer) handler(w http.ResponseWriter, r *http.Request) {
	ms.mu.Lock()
	// Try exact match first
	key := r.Method + " " + r.URL.Path
	resp, ok := ms.responses[key]
	if !ok {
		// Try matching with query string stripped, then with prefix matching
		for k, v := range ms.responses {
			parts := strings.SplitN(k, " ", 2)
			if len(parts) == 2 && parts[0] == r.Method {
				if strings.HasPrefix(r.URL.Path, parts[1]) || r.URL.Path == parts[1] {
					resp = v
					ok = true
					break
				}
			}
		}
	}
	ms.mu.Unlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "no mock registered for " + key})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	if resp.Body != nil {
		json.NewEncoder(w).Encode(resp.Body)
	}
}
