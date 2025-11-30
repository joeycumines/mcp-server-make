package helpproc

import (
	"testing"
)

func TestProcessHelpOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "no Notes line",
			input:    "Usage:\n  make <target>\n\nTargets:\n  build  Build the project",
			expected: "Usage:\n  make <target>\n\nTargets:\n  build  Build the project",
		},
		{
			name:     "Notes at end",
			input:    "Usage:\n  make <target>\n\nNotes\n\nThis is a note",
			expected: "Usage:\n  make <target>\n",
		},
		{
			name:     "Notes in middle",
			input:    "Header\nNotes\nSome content\nMore content",
			expected: "Header",
		},
		{
			name:     "Notes with extra whitespace should not match",
			input:    "Header\n Notes\nContent",
			expected: "Header\n Notes\nContent",
		},
		{
			name:     "NotesExtended should not match",
			input:    "Header\nNotesExtended\nContent",
			expected: "Header\nNotesExtended\nContent",
		},
		{
			name:     "Notes only",
			input:    "Notes",
			expected: "",
		},
		{
			name: "complex example from requirements",
			input: `
Usage:
  gmake <target>

Standard Targets
  all               Builds all, and (non-standard per GNU) runs all checks.
  clean             Cleans up outputs of other targets, e.g. removing coverage files.

Other Targets
  help              Display this help.
  h                 Alias for help.

Notes

Makefile:
  Extensible multi-module Makefile tailored for Go projects.
  
  This Makefile is designed to be used in a monorepo, where multiple Go modules
  are contained within a single repository.`,
			expected: `
Usage:
  gmake <target>

Standard Targets
  all               Builds all, and (non-standard per GNU) runs all checks.
  clean             Cleans up outputs of other targets, e.g. removing coverage files.

Other Targets
  help              Display this help.
  h                 Alias for help.
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessHelpOutput(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessHelpOutput() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatHelpPreamble(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace only",
			input:    "   \t\n",
			expected: "",
		},
		{
			name:     "simple text",
			input:    "Hello",
			expected: "Hello\n\n",
		},
		{
			name:     "text with trailing whitespace",
			input:    "Hello  \n\n",
			expected: "Hello\n\n",
		},
		{
			name:     "text with trailing tabs",
			input:    "Hello\t\t",
			expected: "Hello\n\n",
		},
		{
			name:     "multi-line text",
			input:    "Line 1\nLine 2",
			expected: "Line 1\nLine 2\n\n",
		},
		{
			name:     "multi-line text with trailing newlines",
			input:    "Line 1\nLine 2\n\n\n",
			expected: "Line 1\nLine 2\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatHelpPreamble(tt.input)
			if result != tt.expected {
				t.Errorf("FormatHelpPreamble() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestNotesRegex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "exact Notes",
			input:    "Notes",
			expected: true,
		},
		{
			name:     "Notes with leading space",
			input:    " Notes",
			expected: false,
		},
		{
			name:     "Notes with trailing space",
			input:    "Notes ",
			expected: false,
		},
		{
			name:     "Notes with content after",
			input:    "Notes:",
			expected: false,
		},
		{
			name:     "NotesExtended",
			input:    "NotesExtended",
			expected: false,
		},
		{
			name:     "lowercase notes",
			input:    "notes",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NotesRegex.MatchString(tt.input)
			if result != tt.expected {
				t.Errorf("NotesRegex.MatchString(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
