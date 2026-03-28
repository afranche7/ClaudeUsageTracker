package ui

import (
	"fmt"
	"time"

	"github.com/afranche7/claudeusagetracker/internal/tracker"
)

// FormatTokens formats a token count for display (e.g., "1.2M", "450K", "3,200")
func FormatTokens(n int64) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	case n >= 10_000:
		return fmt.Sprintf("%.0fK", float64(n)/1_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}

// FormatCost formats a USD amount for display
func FormatCost(usd float64) string {
	if usd >= 100 {
		return fmt.Sprintf("%.0f", usd)
	}
	if usd >= 10 {
		return fmt.Sprintf("%.1f", usd)
	}
	return fmt.Sprintf("%.2f", usd)
}

// FormatDuration formats a duration like "2h 15m" or "45m"
func FormatDuration(d time.Duration) string {
	if d <= 0 {
		return "—"
	}
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	}
	return "<1m"
}

// BuildWidgetData constructs the template data from tracker stats
func BuildWidgetData() WidgetData {
	session := tracker.GetCurrentSessionStats()
	today := tracker.GetTodayStats()
	weekly := tracker.GetWeeklyStats()
	plan := tracker.GetPlanLimits()

	sessionActive := !session.LastActivity.IsZero() &&
		time.Since(session.LastActivity) < 5*time.Minute

	var sessionDur time.Duration
	if !session.StartTime.IsZero() && !session.LastActivity.IsZero() {
		sessionDur = session.LastActivity.Sub(session.StartTime)
	}

	todayPct := 0.0
	if plan.DailyUSD > 0 {
		todayPct = (today.TotalCostUSD / plan.DailyUSD) * 100
	}
	weeklyPct := 0.0
	if plan.WeeklyUSD > 0 {
		weeklyPct = (weekly.TotalCostUSD / plan.WeeklyUSD) * 100
	}

	return WidgetData{
		SessionInputTokens:  FormatTokens(session.InputTokens),
		SessionOutputTokens: FormatTokens(session.OutputTokens),
		SessionCacheTokens:  FormatTokens(session.CacheReadTokens + session.CacheWriteTokens),
		SessionCost:         FormatCost(session.TotalCostUSD),
		SessionRequests:     session.RequestCount,
		SessionModel:        session.Model,
		SessionDuration:     FormatDuration(sessionDur),
		SessionActive:       sessionActive,

		TodayCost:      FormatCost(today.TotalCostUSD),
		TodayRequests:  today.RequestCount,
		TodaySessions:  today.SessionCount,
		TodayTokens:    FormatTokens(today.InputTokens + today.OutputTokens),
		TodayPctOfPlan: clampPct(todayPct),

		WeeklyCost:      FormatCost(weekly.TotalCostUSD),
		WeeklyRequests:  weekly.RequestCount,
		WeeklySessions:  weekly.SessionCount,
		WeeklyTokens:    FormatTokens(weekly.InputTokens + weekly.OutputTokens),
		WeeklyPctOfPlan: clampPct(weeklyPct),

		PlanLabel:   plan.Label,
		DailyLimit:  fmt.Sprintf("$%.0f", plan.DailyUSD),
		WeeklyLimit: fmt.Sprintf("$%.0f", plan.WeeklyUSD),
	}
}
