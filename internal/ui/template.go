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
    --bg-dark: rgba(15, 15, 20, 0.55);
    --bg-card: rgba(255, 255, 255, 0.06);
    --bg-card-alt: rgba(255, 255, 255, 0.04);
    --bg-card-hover: rgba(255, 255, 255, 0.09);
    --accent: #D4A27C;
    --accent-light: #E8C9A8;
    --accent-deep: #C4956C;
    --accent-glow: rgba(212, 162, 124, 0.25);
    --text-primary: rgba(245, 240, 235, 0.95);
    --text-secondary: rgba(180, 180, 185, 0.8);
    --text-dim: rgba(140, 140, 150, 0.5);
    --border: rgba(255, 255, 255, 0.08);
    --border-glow: rgba(255, 255, 255, 0.12);
    --success: #7DB88B;
    --warning: #E5A84B;
    --danger: #D4736D;
    --bar-bg: rgba(255, 255, 255, 0.05);
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
    backdrop-filter: blur(20px) saturate(1.3);
    -webkit-backdrop-filter: blur(20px) saturate(1.3);
    border: 1px solid var(--border);
    border-radius: 16px;
    padding: 16px;
    width: 268px;
    min-height: 100px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3),
                inset 0 0.5px 0 rgba(255, 255, 255, 0.06);
  }

  /* ── Header ─────────────────── */
  .header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 14px;
    padding-bottom: 10px;
    border-bottom: 1px solid var(--border);
  }
  .header-left {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .logo {
    width: 22px;
    height: 22px;
    border-radius: 6px;
    background: linear-gradient(135deg, var(--accent), var(--accent-deep));
    box-shadow: 0 2px 8px rgba(212, 162, 124, 0.3);
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 700;
    font-size: 11px;
    color: #0F0F14;
  }
  .title {
    font-size: 13px;
    font-weight: 600;
    color: var(--text-primary);
    letter-spacing: -0.2px;
  }
  .status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
  }
  .status-dot.active {
    background: var(--success);
    box-shadow: 0 0 8px rgba(125, 184, 139, 0.5);
    animation: pulse 2.5s ease-in-out infinite;
  }
  .status-dot.inactive {
    background: var(--text-dim);
  }
  @keyframes pulse {
    0%, 100% { box-shadow: 0 0 6px rgba(125, 184, 139, 0.4); }
    50%      { box-shadow: 0 0 12px rgba(125, 184, 139, 0.7), 0 0 4px rgba(125, 184, 139, 0.3); }
  }

  /* ── Section ────────────────── */
  .section {
    background: var(--bg-card);
    backdrop-filter: blur(12px);
    -webkit-backdrop-filter: blur(12px);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 12px 14px;
    margin-bottom: 10px;
    transition: background 0.2s ease, border-color 0.2s ease;
  }
  .section:hover {
    background: var(--bg-card-hover);
    border-color: var(--border-glow);
  }
  .section:last-child { margin-bottom: 0; }

  .section-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 10px;
  }
  .section-label {
    font-size: 10px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--text-secondary);
  }
  .section-badge {
    font-size: 9px;
    padding: 3px 8px;
    border-radius: 6px;
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid rgba(255, 255, 255, 0.05);
    color: var(--text-dim);
  }

  /* ── Cost display ───────────── */
  .cost-row {
    display: flex;
    align-items: baseline;
    gap: 3px;
    margin-bottom: 8px;
  }
  .cost-value {
    font-size: 24px;
    font-weight: 700;
    color: var(--accent);
    letter-spacing: -0.5px;
    line-height: 1;
  }
  .cost-currency {
    font-size: 13px;
    color: var(--accent-deep);
    font-weight: 600;
  }
  .cost-limit {
    font-size: 10px;
    color: var(--text-dim);
    margin-left: auto;
  }

  /* ── Progress bar ───────────── */
  .progress-bar {
    width: 100%;
    height: 6px;
    background: var(--bar-bg);
    border-radius: 3px;
    margin-bottom: 10px;
    overflow: hidden;
  }
  .progress-fill {
    height: 100%;
    border-radius: 3px;
    transition: width 0.8s cubic-bezier(0.22, 1, 0.36, 1);
  }
  .fill-normal {
    background: linear-gradient(90deg, var(--accent-deep), var(--accent));
    box-shadow: 0 0 8px var(--accent-glow);
  }
  .fill-warning {
    background: linear-gradient(90deg, #c48a30, var(--warning));
    box-shadow: 0 0 8px rgba(229, 168, 75, 0.25);
  }
  .fill-danger {
    background: linear-gradient(90deg, #b85a54, var(--danger));
    box-shadow: 0 0 8px rgba(212, 115, 109, 0.25);
  }

  /* ── Stat grid ──────────────── */
  .stat-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 6px 14px;
  }
  .stat-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 3px 0;
  }
  .stat-label {
    font-size: 10px;
    color: var(--text-dim);
  }
  .stat-value {
    font-size: 10.5px;
    font-weight: 600;
    color: var(--text-secondary);
    font-variant-numeric: tabular-nums;
  }

  /* ── Tokens display ─────────── */
  .token-row {
    display: flex;
    gap: 6px;
    margin-top: 8px;
    padding-top: 8px;
    border-top: 1px solid var(--border);
    flex-wrap: wrap;
  }
  .token-chip {
    font-size: 9.5px;
    padding: 3px 8px;
    border-radius: 6px;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.04);
    white-space: nowrap;
  }
  .token-chip .label { color: var(--text-dim); font-size: 8.5px; }
  .token-chip .val { color: var(--text-secondary); font-weight: 600; }

  /* ── Footer ─────────────────── */
  .footer {
    text-align: center;
    padding-top: 10px;
    margin-top: 2px;
    font-size: 9px;
    color: var(--text-dim);
    letter-spacing: 0.3px;
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
    <div class="status-dot {{if .SessionActive}}active{{else}}inactive{{end}}" title="{{if .SessionActive}}Session active{{else}}No active session{{end}}"></div>
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
