package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientAppliesConfiguredAuthenticationHeader(t *testing.T) {
	tests := []struct {
		authType   string
		headerName string
		headerWant string
	}{
		{authType: "Bearer", headerName: "Authorization", headerWant: "Bearer secret"},
		{authType: "api-key", headerName: "api-key", headerWant: "secret"},
		{authType: "x-api-key", headerName: "x-api-key", headerWant: "secret"},
	}
	for _, tt := range tests {
		t.Run(tt.authType, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := r.Header.Get(tt.headerName); got != tt.headerWant {
					t.Errorf("%s = %q, want %q", tt.headerName, got, tt.headerWant)
				}
				_ = json.NewEncoder(w).Encode(map[string]any{
					"choices": []any{map[string]any{"message": map[string]any{"content": "ok"}}},
				})
			}))
			defer server.Close()

			client := NewClientWithAuth("secret", server.URL, "test-model", tt.authType)
			got, err := client.Complete(context.Background(), []Message{{Role: "user", Content: "hello"}})
			if err != nil {
				t.Fatal(err)
			}
			if got != "ok" {
				t.Fatalf("响应 = %q, want ok", got)
			}
		})
	}
}
