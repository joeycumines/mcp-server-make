package makeexec

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestNewExecutor(t *testing.T) {
	executor := NewExecutor("/usr/bin/make", "/tmp", 60, 4)

	if executor.MakePath != "/usr/bin/make" {
		t.Errorf("MakePath = %s, want /usr/bin/make", executor.MakePath)
	}
	if executor.DefaultWorkDir != "/tmp" {
		t.Errorf("DefaultWorkDir = %s, want /tmp", executor.DefaultWorkDir)
	}
	if executor.Timeout != 60*time.Second {
		t.Errorf("Timeout = %v, want %v", executor.Timeout, 60*time.Second)
	}
	if executor.Semaphore == nil {
		t.Error("Semaphore should not be nil")
	}
}

func TestExecute_SimpleCommand(t *testing.T) {
	// Skip if make is not available
	makePath, err := exec.LookPath("make")
	if err != nil {
		t.Skip("make not found in PATH")
	}

	// Create a temporary directory with a simple Makefile
	tmpDir, err := os.MkdirTemp("", "makeexec-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	makefileContent := `.PHONY: hello
hello:
	@echo "Hello from test"
`
	makefilePath := filepath.Join(tmpDir, "Makefile")
	if err := os.WriteFile(makefilePath, []byte(makefileContent), 0644); err != nil {
		t.Fatalf("Failed to write Makefile: %v", err)
	}

	executor := NewExecutor(makePath, tmpDir, 10, 1)
	params := MakeParams{Target: "hello"}

	result, err := executor.Execute(context.Background(), params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", result.ExitCode)
	}

	expectedOutput := "Hello from test\n"
	if result.Stdout != expectedOutput {
		t.Errorf("Stdout = %q, want %q", result.Stdout, expectedOutput)
	}
}

func TestExecute_WithTimeout(t *testing.T) {
	// Skip if make is not available
	makePath, err := exec.LookPath("make")
	if err != nil {
		t.Skip("make not found in PATH")
	}

	// Create a temporary directory with a Makefile that sleeps
	tmpDir, err := os.MkdirTemp("", "makeexec-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	makefileContent := `.PHONY: slow
slow:
	@sleep 10
`
	makefilePath := filepath.Join(tmpDir, "Makefile")
	if err := os.WriteFile(makefilePath, []byte(makefileContent), 0644); err != nil {
		t.Fatalf("Failed to write Makefile: %v", err)
	}

	// Use a very short timeout (1 second)
	executor := NewExecutor(makePath, tmpDir, 1, 1)
	params := MakeParams{Target: "slow"}

	result, err := executor.Execute(context.Background(), params)

	// Should return error due to timeout
	if err == nil {
		t.Error("Expected error due to timeout, got nil")
	}

	// Exit code should be -1 for timeout
	if result.ExitCode != -1 {
		t.Errorf("ExitCode = %d, want -1 for timeout", result.ExitCode)
	}
}

func TestExecute_NonExistentTarget(t *testing.T) {
	// Skip if make is not available
	makePath, err := exec.LookPath("make")
	if err != nil {
		t.Skip("make not found in PATH")
	}

	// Create a temporary directory with a simple Makefile
	tmpDir, err := os.MkdirTemp("", "makeexec-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	makefileContent := `.PHONY: hello
hello:
	@echo "Hello"
`
	makefilePath := filepath.Join(tmpDir, "Makefile")
	if err := os.WriteFile(makefilePath, []byte(makefileContent), 0644); err != nil {
		t.Fatalf("Failed to write Makefile: %v", err)
	}

	executor := NewExecutor(makePath, tmpDir, 10, 1)
	params := MakeParams{Target: "nonexistent"}

	result, err := executor.Execute(context.Background(), params)

	// Should return error because target doesn't exist
	if err == nil {
		t.Error("Expected error for non-existent target, got nil")
	}

	// Exit code should be non-zero
	if result.ExitCode == 0 {
		t.Errorf("ExitCode = %d, want non-zero for failed target", result.ExitCode)
	}
}

func TestSerializeResult(t *testing.T) {
	result := &Result{
		Stdout:     "Hello World",
		Stderr:     "",
		ExitCode:   0,
		DurationMs: 100,
	}

	jsonStr, err := SerializeResult(result)
	if err != nil {
		t.Fatalf("SerializeResult failed: %v", err)
	}

	// Check that the JSON contains expected fields
	if jsonStr == "" {
		t.Error("SerializeResult returned empty string")
	}

	// Should contain expected content
	expectedContent := []string{
		`"stdout":"Hello World"`,
		`"exit_code":0`,
		`"duration_ms":100`,
	}

	for _, expected := range expectedContent {
		if !contains(jsonStr, expected) {
			t.Errorf("JSON output %q does not contain %q", jsonStr, expected)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
