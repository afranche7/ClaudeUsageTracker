package ui

import (
	"fmt"
	"html/template"
	"strings"
)

// WidgetData holds all data passed to the HTML template
type WidgetData struct {
	// Session stats
	SessionInputTokens  string
	SessionOutputTokens string
	SessionCacheTokens  string
	SessionCost         string
	SessionRequests     int
	SessionModel        string
	SessionDuration     string
	SessionActive       bool

	// Today stats
	TodayCost      string
	TodayRequests  int
	TodaySessions  int
	TodayTokens    string
	TodayPctOfPlan float64

	// Weekly stats
	WeeklyCost      string
	WeeklyRequests  int
	WeeklySessions  int
	WeeklyTokens    string
	WeeklyPctOfPlan float64

	// Plan info
	PlanLabel    string
	DailyLimit   string
	WeeklyLimit  string
}

// RenderHTML generates the widget HTML with current data
func RenderHTML(data WidgetData) string {
	var buf strings.Builder
	tmpl := template.Must(template.New("widget").Parse(widgetTemplate))
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Sprintf("<html><body style='color:white;background:#1B1B1F'>Error: %s</body></html>", err)
	}
	return buf.String()
}

func clampPct(v float64) float64 {
	if v > 100 {
		return 100
	}
	if v < 0 {
		return 0
	}
	return v
}

const widgetTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }

  :root {
    --bg-dark: #1B1B1F;
    --bg-card: #252529;
    --bg-card-alt: #2A2A30;
    --accent: #D4A27C;
    --accent-light: #E8C9A8;
    --accent-deep: #C4956C;
    --text-primary: #F5F0EB;
    --text-secondary: #A0A0A0;
    --text-dim: #6B6B6B;
    --border: #333338;
    --success: #7DB88B;
    --warning: #E5A84B;
    --danger: #D4736D;
    --bar-bg: #2A2A30;
  }

  html, body {
    font-family: "Segoe UI", -apple-system, BlinkMacSystemFont, sans-serif;
    font-size: 12px;
    color: var(--text-primary);
    background: transparent;
    -webkit-user-select: none;
    user-select: none;
    overflow: hidden;
    height: 100%;
  }

  .widget {
    background: var(--bg-dark);
    border: 1px solid var(--border);
    border-radius: 14px;
    padding: 14px;
    width: 268px;
    min-height: 100px;
    opacity: 0.95;
  }

  /* ── Header ─────────────────── */
  .header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 12px;
    padding-bottom: 8px;
    border-bottom: 1px solid var(--border);
  }
  .header-left {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .logo {
    width: 20px;
    height: 20px;
    border-radius: 5px;
    background: linear-gradient(135deg, var(--accent), var(--accent-deep));
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 700;
    font-size: 11px;
    color: var(--bg-dark);
  }
  .title {
    font-size: 13px;
    font-weight: 600;
    color: var(--text-primary);
    letter-spacing: -0.3px;
  }
  .status-dot {
    width: 7px;
    height: 7px;
    border-radius: 50%;
    background: {{if .SessionActive}}var(--success){{else}}var(--text-dim){{end}};
    {{if .SessionActive}}box-shadow: 0 0 6px rgba(125, 184, 139, 0.5);{{end}}
  }

  /* ── Section ────────────────── */
  .section {
    background: var(--bg-card);
    border-radius: 10px;
    padding: 10px 11px;
    margin-bottom: 8px;
  }
  .section:last-child { margin-bottom: 0; }

  .section-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 8px;
  }
  .section-label {
    font-size: 10px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.8px;
    color: var(--text-secondary);
  }
  .section-badge {
    font-size: 9px;
    padding: 2px 6px;
    border-radius: 4px;
    background: var(--bg-card-alt);
    color: var(--text-dim);
  }

  /* ── Cost display ───────────── */
  .cost-row {
    display: flex;
    align-items: baseline;
    gap: 4px;
    margin-bottom: 6px;
  }
  .cost-value {
    font-size: 22px;
    font-weight: 700;
    color: var(--accent);
    letter-spacing: -0.5px;
    line-height: 1;
  }
  .cost-currency {
    font-size: 11px;
    color: var(--accent-deep);
    font-weight: 500;
  }
  .cost-limit {
    font-size: 10px;
    color: var(--text-dim);
    margin-left: auto;
  }

  /* ── Progress bar ───────────── */
  .progress-bar {
    width: 100%;
    height: 4px;
    background: var(--bar-bg);
    border-radius: 2px;
    margin-bottom: 8px;
    overflow: hidden;
  }
  .progress-fill {
    height: 100%;
    border-radius: 2px;
    transition: width 0.5s ease;
  }
  .fill-normal { background: var(--accent); }
  .fill-warning { background: var(--warning); }
  .fill-danger { background: var(--danger); }

  /* ── Stat grid ──────────────── */
  .stat-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 4px 12px;
  }
  .stat-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 2px 0;
  }
  .stat-label {
    font-size: 10px;
    color: var(--text-dim);
  }
  .stat-value {
    font-size: 10px;
    font-weight: 600;
    color: var(--text-secondary);
    font-variant-numeric: tabular-nums;
  }

  /* ── Tokens display ─────────── */
  .token-row {
    display: flex;
    gap: 8px;
    margin-top: 6px;
    padding-top: 6px;
    border-top: 1px solid var(--border);
  }
  .token-chip {
    font-size: 9px;
    padding: 2px 6px;
    border-radius: 4px;
    background: var(--bg-card-alt);
  }
  .token-chip .label { color: var(--text-dim); }
  .token-chip .val { color: var(--text-secondary); font-weight: 600; }

  /* ── Footer ─────────────────── */
  .footer {
    text-align: center;
    padding-top: 6px;
    font-size: 9px;
    color: var(--text-dim);
  }

  /* ── Drag area ──────────────── */
  .drag-area {
    -webkit-app-region: drag;
    cursor: grab;
  }
  .no-drag {
    -webkit-app-region: no-drag;
  }
</style>
</head>
<body>
<div class="widget drag-area">

  <!-- Header -->
  <div class="header">
    <div class="header-left">
      <div class="logo">C</div>
      <span class="title">Claude Usage</span>
    </div>
    <div class="status-dot" title="{{if .SessionActive}}Session active{{else}}No active session{{end}}"></div>
  </div>

  <!-- Current Session -->
  <div class="section">
    <div class="section-header">
      <span class="section-label">Session</span>
      <span class="section-badge">{{.SessionDuration}}</span>
    </div>
    <div class="cost-row">
      <span class="cost-currency">$</span>
      <span class="cost-value">{{.SessionCost}}</span>
      <span class="cost-limit">{{.SessionRequests}} reqs</span>
    </div>
    <div class="token-row">
      <div class="token-chip"><span class="label">IN </span><span class="val">{{.SessionInputTokens}}</span></div>
      <div class="token-chip"><span class="label">OUT </span><span class="val">{{.SessionOutputTokens}}</span></div>
      <div class="token-chip"><span class="label">CACHE </span><span class="val">{{.SessionCacheTokens}}</span></div>
    </div>
  </div>

  <!-- Today -->
  <div class="section">
    <div class="section-header">
      <span class="section-label">Today</span>
      <span class="section-badge">{{.TodaySessions}} sessions</span>
    </div>
    <div class="cost-row">
      <span class="cost-currency">$</span>
      <span class="cost-value">{{.TodayCost}}</span>
      <span class="cost-limit">/ {{.DailyLimit}}</span>
    </div>
    <div class="progress-bar">
      <div class="progress-fill {{if gt .TodayPctOfPlan 85.0}}fill-danger{{else if gt .TodayPctOfPlan 60.0}}fill-warning{{else}}fill-normal{{end}}"
           style="width: {{printf "%.1f" .TodayPctOfPlan}}%"></div>
    </div>
    <div class="stat-grid">
      <div class="stat-item"><span class="stat-label">Tokens</span><span class="stat-value">{{.TodayTokens}}</span></div>
      <div class="stat-item"><span class="stat-label">Requests</span><span class="stat-value">{{.TodayRequests}}</span></div>
    </div>
  </div>

  <!-- Weekly -->
  <div class="section">
    <div class="section-header">
      <span class="section-label">This Week</span>
      <span class="section-badge">{{.WeeklySessions}} sessions</span>
    </div>
    <div class="cost-row">
      <span class="cost-currency">$</span>
      <span class="cost-value">{{.WeeklyCost}}</span>
      <span class="cost-limit">/ {{.WeeklyLimit}}</span>
    </div>
    <div class="progress-bar">
      <div class="progress-fill {{if gt .WeeklyPctOfPlan 85.0}}fill-danger{{else if gt .WeeklyPctOfPlan 60.0}}fill-warning{{else}}fill-normal{{end}}"
           style="width: {{printf "%.1f" .WeeklyPctOfPlan}}%"></div>
    </div>
    <div class="stat-grid">
      <div class="stat-item"><span class="stat-label">Tokens</span><span class="stat-value">{{.WeeklyTokens}}</span></div>
      <div class="stat-item"><span class="stat-label">Requests</span><span class="stat-value">{{.WeeklyRequests}}</span></div>
    </div>
  </div>

  <!-- Footer -->
  <div class="footer">
    {{.PlanLabel}} Plan · Right-click tray to exit
  </div>

</div>
</body>
</html>`
