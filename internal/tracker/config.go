package tracker

import (
	"os"
	"path/filepath"
)

// Model pricing in USD per 1M tokens
type ModelPricing struct {
	Input  float64
	Output float64
}

var PricingTable = map[string]ModelPricing{
	"claude-opus-4-6":   {Input: 15.00, Output: 75.00},
	"claude-sonnet-4-6": {Input: 3.00, Output: 15.00},
	"claude-haiku-4-5":  {Input: 0.80, Output: 4.00},
	"default":           {Input: 3.00, Output: 15.00},
}

// PlanLimits defines usage budget for a subscription plan
type PlanLimits struct {
	Label     string
	DailyUSD  float64
	WeeklyUSD float64
}

var Plans = map[string]PlanLimits{
	"pro":     {Label: "Pro", DailyUSD: 20.0, WeeklyUSD: 100.0},
	"max_5x":  {Label: "Max 5x", DailyUSD: 100.0, WeeklyUSD: 500.0},
	"max_20x": {Label: "Max 20x", DailyUSD: 400.0, WeeklyUSD: 2000.0},
}

// ClaudeHomePath returns the path to Claude Code's config directory
func ClaudeHomePath() string {
	if env := os.Getenv("CLAUDE_HOME"); env != "" {
		return env
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".claude")
}

// ActivePlan returns the configured plan name
func ActivePlan() string {
	if p := os.Getenv("CLAUDE_PLAN"); p != "" {
		return p
	}
	return "max_5x"
}

// GetPlanLimits returns limits for the active plan
func GetPlanLimits() PlanLimits {
	plan := ActivePlan()
	if l, ok := Plans[plan]; ok {
		return l
	}
	return Plans["max_5x"]
}
