package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMustRuntimeValuePrefersProcessEnv(t *testing.T) {
	t.Setenv("GRPC_PORT", "50053")
	t.Setenv("ENV_PORTS_FILE", filepath.Join(t.TempDir(), "missing.env.ports"))

	if got := mustRuntimeValue("GRPC_PORT"); got != "50053" {
		t.Fatalf("mustRuntimeValue(GRPC_PORT) = %q, want %q", got, "50053")
	}
}

func TestMustRuntimeValueFallsBackToEnvPortsFile(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env.ports")
	if err := os.WriteFile(envFile, []byte("REST_PORT=8083\nGRPC_PORT=50053\n"), 0o644); err != nil {
		t.Fatalf("write env file failed: %v", err)
	}

	t.Setenv("ENV_PORTS_FILE", envFile)
	t.Setenv("REST_PORT", "")

	if got := mustRuntimeValue("REST_PORT"); got != "8083" {
		t.Fatalf("mustRuntimeValue(REST_PORT) = %q, want %q", got, "8083")
	}
}
