package tools

import (
	"testing"
)

func TestFormatTime(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", "Not available"},
		{"Thu, 15 Jun 2023 10:30:00 UTC", "15 Jun 2023 10:30:00"},
		{"2023-06-15T10:30:00Z", "15 Jun 2023 10:30:00"},
		{"completely invalid", "Unknown format (completely invalid)"},
	}
	for _, tt := range tests {
		if got := FormatTime(tt.input); got != tt.want {
			t.Errorf("FormatTime(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
