package ui

import (
	"testing"
	"time"
)

func TestFormatTokens(t *testing.T) {
	tests := []struct {
		input int64
		want  string
	}{
		{0, "0"},
		{500, "500"},
		{1500, "1.5K"},
		{15000, "15K"},
		{150000, "150K"},
		{1500000, "1.5M"},
		{15000000, "15.0M"},
	}
	for _, tt := range tests {
		got := FormatTokens(tt.input)
		if got != tt.want {
			t.Errorf("FormatTokens(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatCost(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{0.0, "0.00"},
		{0.50, "0.50"},
		{5.123, "5.12"},
		{15.5, "15.5"},
		{150.0, "150"},
	}
	for _, tt := range tests {
		got := FormatCost(tt.input)
		if got != tt.want {
			t.Errorf("FormatCost(%f) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		input time.Duration
		want  string
	}{
		{0, "—"},
		{-1 * time.Second, "—"},
		{30 * time.Second, "<1m"},
		{5 * time.Minute, "5m"},
		{90 * time.Minute, "1h 30m"},
		{3*time.Hour + 15*time.Minute, "3h 15m"},
	}
	for _, tt := range tests {
		got := FormatDuration(tt.input)
		if got != tt.want {
			t.Errorf("FormatDuration(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestRenderHTML(t *testing.T) {
	data := WidgetData{
		SessionInputTokens:  "1.5M",
		SessionOutputTokens: "500K",
		SessionCacheTokens:  "200K",
		SessionCost:         "5.25",
		SessionRequests:     42,
		SessionModel:        "claude-sonnet-4-6",
		SessionDuration:     "1h 30m",
		SessionActive:       true,
		TodayCost:           "12.50",
		TodayRequests:       100,
		TodaySessions:       3,
		TodayTokens:         "5.0M",
		TodayPctOfPlan:      25.0,
		WeeklyCost:          "45.00",
		WeeklyRequests:      500,
		WeeklySessions:      15,
		WeeklyTokens:        "25.0M",
		WeeklyPctOfPlan:     45.0,
		PlanLabel:           "Max 5x",
		DailyLimit:          "$100",
		WeeklyLimit:         "$500",
	}

	html := RenderHTML(data)

	if len(html) == 0 {
		t.Fatal("RenderHTML returned empty string")
	}
	if !contains(html, "Claude Usage") {
		t.Error("HTML missing title")
	}
	if !contains(html, "5.25") {
		t.Error("HTML missing session cost")
	}
	if !contains(html, "var(--success)") {
		t.Error("HTML missing active status indicator")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr) && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
