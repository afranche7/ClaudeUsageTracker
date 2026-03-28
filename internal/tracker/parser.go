package tracker

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// SessionStats holds usage data for a single session
type SessionStats struct {
	InputTokens     int64     `json:"input_tokens"`
	OutputTokens    int64     `json:"output_tokens"`
	CacheReadTokens int64     `json:"cache_read_tokens"`
	CacheWriteTokens int64    `json:"cache_write_tokens"`
	TotalCostUSD    float64   `json:"total_cost_usd"`
	RequestCount    int       `json:"request_count"`
	Model           string    `json:"model"`
	StartTime       time.Time `json:"start_time"`
	LastActivity    time.Time `json:"last_activity"`
	SessionID       string    `json:"session_id"`
}

// AggregateStats holds combined usage over a period
type AggregateStats struct {
	InputTokens      int64              `json:"input_tokens"`
	OutputTokens     int64              `json:"output_tokens"`
	CacheReadTokens  int64              `json:"cache_read_tokens"`
	CacheWriteTokens int64              `json:"cache_write_tokens"`
	TotalCostUSD     float64            `json:"total_cost_usd"`
	RequestCount     int                `json:"request_count"`
	SessionCount     int                `json:"session_count"`
	ModelsUsed       map[string]int     `json:"models_used"`
}

// CalculateCost computes USD cost for given token counts
func CalculateCost(inputTokens, outputTokens int64, model string, cacheRead, cacheWrite int64) float64 {
	pricing, ok := PricingTable[model]
	if !ok {
		pricing = PricingTable["default"]
	}
	inputCost := float64(inputTokens) / 1_000_000 * pricing.Input
	outputCost := float64(outputTokens) / 1_000_000 * pricing.Output
	cacheReadCost := float64(cacheRead) / 1_000_000 * pricing.Input * 0.1
	cacheWriteCost := float64(cacheWrite) / 1_000_000 * pricing.Input * 1.25
	return inputCost + outputCost + cacheReadCost + cacheWriteCost
}

// findSessionFiles locates all JSONL session files under Claude's projects dir
func findSessionFiles() []string {
	claudeHome := ClaudeHomePath()
	if claudeHome == "" {
		return nil
	}
	projectsDir := filepath.Join(claudeHome, "projects")
	var files []string
	_ = filepath.Walk(projectsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip errors
		}
		if !info.IsDir() && filepath.Ext(path) == ".jsonl" {
			files = append(files, path)
		}
		return nil
	})
	return files
}

// jsonEntry represents a single line from a JSONL session file
type jsonEntry struct {
	Timestamp interface{}            `json:"timestamp"`
	CreatedAt interface{}            `json:"created_at"`
	Ts        interface{}            `json:"ts"`
	Usage     *usageData             `json:"usage"`
	TokenUsage *usageData            `json:"token_usage"`
	Model     string                 `json:"model"`
	Message   json.RawMessage        `json:"message"`
}

type usageData struct {
	InputTokens              int64 `json:"input_tokens"`
	OutputTokens             int64 `json:"output_tokens"`
	CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
	CacheReadTokens          int64 `json:"cache_read_tokens"`
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
	CacheWriteTokens         int64 `json:"cache_write_tokens"`
}

type messageWithUsage struct {
	Usage *usageData `json:"usage"`
	Model string     `json:"model"`
}

func parseTimestamp(v interface{}) (time.Time, bool) {
	if v == nil {
		return time.Time{}, false
	}
	switch t := v.(type) {
	case float64:
		return time.Unix(int64(t), 0), true
	case string:
		if ts, err := time.Parse(time.RFC3339, t); err == nil {
			return ts, true
		}
		if ts, err := time.Parse(time.RFC3339Nano, t); err == nil {
			return ts, true
		}
		if ts, err := time.Parse("2006-01-02T15:04:05", t); err == nil {
			return ts, true
		}
	}
	return time.Time{}, false
}

func parseSessionFile(path string, since *time.Time) *SessionStats {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	stats := &SessionStats{
		SessionID: filepath.Base(path),
	}
	foundUsage := false

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 256*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var entry jsonEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			continue
		}

		// Extract timestamp
		for _, tsVal := range []interface{}{entry.Timestamp, entry.CreatedAt, entry.Ts} {
			if ts, ok := parseTimestamp(tsVal); ok {
				if stats.StartTime.IsZero() || ts.Before(stats.StartTime) {
					stats.StartTime = ts
				}
				if stats.LastActivity.IsZero() || ts.After(stats.LastActivity) {
					stats.LastActivity = ts
				}
				break
			}
		}

		// Find usage data
		usage := entry.Usage
		if usage == nil {
			usage = entry.TokenUsage
		}
		if usage == nil && len(entry.Message) > 0 {
			var msg messageWithUsage
			if json.Unmarshal(entry.Message, &msg) == nil {
				usage = msg.Usage
				if msg.Model != "" && entry.Model == "" {
					entry.Model = msg.Model
				}
			}
		}

		if usage != nil {
			stats.InputTokens += usage.InputTokens
			stats.OutputTokens += usage.OutputTokens
			cr := usage.CacheReadInputTokens
			if cr == 0 {
				cr = usage.CacheReadTokens
			}
			stats.CacheReadTokens += cr
			cw := usage.CacheCreationInputTokens
			if cw == 0 {
				cw = usage.CacheWriteTokens
			}
			stats.CacheWriteTokens += cw
			stats.RequestCount++
			foundUsage = true
		}

		if entry.Model != "" {
			stats.Model = entry.Model
		}
	}

	if !foundUsage {
		return nil
	}
	if since != nil && !stats.LastActivity.IsZero() && stats.LastActivity.Before(*since) {
		return nil
	}

	stats.TotalCostUSD = CalculateCost(
		stats.InputTokens, stats.OutputTokens, stats.Model,
		stats.CacheReadTokens, stats.CacheWriteTokens,
	)
	return stats
}

// GetCurrentSessionStats returns stats for the most recently active session
func GetCurrentSessionStats() SessionStats {
	files := findSessionFiles()
	if len(files) == 0 {
		return SessionStats{}
	}

	// Sort by modification time, newest first
	type fileWithMtime struct {
		path  string
		mtime time.Time
	}
	var withTimes []fileWithMtime
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			continue
		}
		withTimes = append(withTimes, fileWithMtime{f, info.ModTime()})
	}
	sort.Slice(withTimes, func(i, j int) bool {
		return withTimes[i].mtime.After(withTimes[j].mtime)
	})

	if len(withTimes) == 0 {
		return SessionStats{}
	}

	stats := parseSessionFile(withTimes[0].path, nil)
	if stats == nil {
		return SessionStats{}
	}
	return *stats
}

// GetWeeklyStats returns aggregated stats for the past 7 days
func GetWeeklyStats() AggregateStats {
	since := time.Now().AddDate(0, 0, -7)
	return aggregateStats(&since)
}

// GetTodayStats returns aggregated stats for today
func GetTodayStats() AggregateStats {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return aggregateStats(&today)
}

func aggregateStats(since *time.Time) AggregateStats {
	agg := AggregateStats{ModelsUsed: make(map[string]int)}
	files := findSessionFiles()

	for _, f := range files {
		if since != nil {
			info, err := os.Stat(f)
			if err != nil || info.ModTime().Before(*since) {
				continue
			}
		}

		stats := parseSessionFile(f, since)
		if stats == nil {
			continue
		}

		agg.InputTokens += stats.InputTokens
		agg.OutputTokens += stats.OutputTokens
		agg.CacheReadTokens += stats.CacheReadTokens
		agg.CacheWriteTokens += stats.CacheWriteTokens
		agg.TotalCostUSD += stats.TotalCostUSD
		agg.RequestCount += stats.RequestCount
		agg.SessionCount++

		if stats.Model != "" && stats.Model != "unknown" {
			agg.ModelsUsed[stats.Model]++
		}
	}
	return agg
}
