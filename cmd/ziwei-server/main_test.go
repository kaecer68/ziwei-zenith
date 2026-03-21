package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRestPortFromContract(t *testing.T) {
	t.Parallel()

	// 從環境變數讀取期望的 port
	expectedPort := os.Getenv("REST_PORT")
	if expectedPort == "" {
		t.Skip("REST_PORT not set, skipping test")
	}

	tests := []struct {
		name      string
		content   string
		wantPort  string
		wantError string
	}{
		{
			name: "parse first server url port",
			content: `openapi: 3.1.0
servers:
  - url: http://localhost:` + expectedPort + `
paths:
  /api/v1/calculate:
    post: {}
`,
			wantPort: expectedPort,
		},
		{
			name: "support quoted url",
			content: `openapi: 3.1.0
servers:
  - url: "https://example.com:9443"
paths: {}
`,
			wantPort: "9443",
		},
		{
			name: "error when server url has no port",
			content: `openapi: 3.1.0
servers:
  - url: http://localhost
paths: {}
`,
			wantError: "has no explicit port",
		},
		{
			name: "error when servers block has no url",
			content: `openapi: 3.1.0
info:
  title: test
paths: {}
`,
			wantError: "no server url found",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			file := filepath.Join(t.TempDir(), "contract.yaml")
			if err := os.WriteFile(file, []byte(tc.content), 0o644); err != nil {
				t.Fatalf("write contract failed: %v", err)
			}

			got, err := restPortFromContract(file)
			if tc.wantError != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantError)
				}
				if !strings.Contains(err.Error(), tc.wantError) {
					t.Fatalf("expected error containing %q, got %q", tc.wantError, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("restPortFromContract returned error: %v", err)
			}
			if got != tc.wantPort {
				t.Fatalf("expected port %q, got %q", tc.wantPort, got)
			}
		})
	}
}
