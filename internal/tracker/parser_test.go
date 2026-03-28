package tracker

import (
	"math"
	"testing"
)

func TestCalculateCost(t *testing.T) {
	tests := []struct {
		name       string
		input      int64
		output     int64
		model      string
		cacheRead  int64
		cacheWrite int64
		wantMin    float64
		wantMax    float64
	}{
		{
			name:    "zero tokens",
			input:   0, output: 0, model: "claude-sonnet-4-6",
			wantMin: 0, wantMax: 0,
		},
		{
			name:    "1M input tokens sonnet",
			input:   1_000_000, output: 0, model: "claude-sonnet-4-6",
			wantMin: 2.99, wantMax: 3.01,
		},
		{
			name:    "1M output tokens sonnet",
			input:   0, output: 1_000_000, model: "claude-sonnet-4-6",
			wantMin: 14.99, wantMax: 15.01,
		},
		{
			name:    "opus pricing",
			input:   1_000_000, output: 1_000_000, model: "claude-opus-4-6",
			wantMin: 89.9, wantMax: 90.1,
		},
		{
			name:    "unknown model uses default",
			input:   1_000_000, output: 0, model: "unknown-model",
			wantMin: 2.99, wantMax: 3.01,
		},
		{
			name:       "cache tokens",
			input:      0, output: 0, model: "claude-sonnet-4-6",
			cacheRead:  1_000_000, cacheWrite: 1_000_000,
			wantMin:    4.04, wantMax: 4.06,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateCost(tt.input, tt.output, tt.model, tt.cacheRead, tt.cacheWrite)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("CalculateCost() = %f, want between %f and %f", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestParseTimestamp(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		valid bool
	}{
		{"unix float", 1711612800.0, true},
		{"RFC3339", "2024-03-28T12:00:00Z", true},
		{"ISO no TZ", "2024-03-28T12:00:00", true},
		{"nil", nil, false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts, ok := parseTimestamp(tt.input)
			if ok != tt.valid {
				t.Errorf("parseTimestamp(%v) valid = %v, want %v", tt.input, ok, tt.valid)
			}
			if ok && ts.IsZero() {
				t.Error("parseTimestamp returned zero time but ok=true")
			}
		})
	}
}

func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
