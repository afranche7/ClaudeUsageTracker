import glob
import json
import os
from datetime import datetime, timezone, timedelta
from pathlib import Path

from .pricing import calculate_cost


def _claude_dir() -> Path:
    return Path.home() / ".claude"


def get_all_jsonl_files() -> list[Path]:
    pattern = str(_claude_dir() / "projects" / "**" / "*.jsonl")
    return [Path(p) for p in glob.glob(pattern, recursive=True)]


def get_current_session_id() -> str | None:
    """Return the sessionId of the most recently started Claude Code session."""
    sessions_dir = _claude_dir() / "sessions"
    if not sessions_dir.exists():
        return None

    best: dict | None = None
    for path in sessions_dir.glob("*.json"):
        try:
            data = json.loads(path.read_text(encoding="utf-8"))
            if best is None or data.get("startedAt", 0) > best.get("startedAt", 0):
                best = data
        except Exception:
            continue

    return best.get("sessionId") if best else None


def parse_usage_from_file(path: Path) -> list[dict]:
    """Parse a JSONL file and return usage records from assistant messages."""
    records: list[dict] = []
    try:
        text = path.read_text(encoding="utf-8")
    except Exception:
        return records

    for line in text.splitlines():
        line = line.strip()
        if not line:
            continue
        try:
            obj = json.loads(line)
        except json.JSONDecodeError:
            continue

        msg = obj.get("message", {})
        if not isinstance(msg, dict):
            continue

        # Only assistant messages carry usage data
        if msg.get("role") != "assistant":
            continue

        usage = msg.get("usage")
        if not usage:
            continue

        records.append(
            {
                "timestamp": obj.get("timestamp", ""),
                "sessionId": obj.get("sessionId", ""),
                "model": msg.get("model", ""),
                "usage": usage,
            }
        )
    return records


def _parse_ts(ts: str) -> datetime | None:
    try:
        return datetime.fromisoformat(ts.replace("Z", "+00:00"))
    except Exception:
        return None


def aggregate_session(session_id: str | None) -> dict:
    """Sum tokens + cost for a specific session across all JSONL files."""
    totals = _empty_totals()
    if not session_id:
        return totals

    for path in get_all_jsonl_files():
        for rec in parse_usage_from_file(path):
            if rec["sessionId"] == session_id:
                _add_record(totals, rec)
    return totals


def aggregate_weekly() -> dict:
    """Sum tokens + cost for all sessions in the past 7 days."""
    cutoff = datetime.now(tz=timezone.utc) - timedelta(days=7)
    totals = _empty_totals()

    for path in get_all_jsonl_files():
        for rec in parse_usage_from_file(path):
            ts = _parse_ts(rec["timestamp"])
            if ts and ts >= cutoff:
                _add_record(totals, rec)
    return totals


def _empty_totals() -> dict:
    return {
        "input_tokens": 0,
        "output_tokens": 0,
        "cache_creation_input_tokens": 0,
        "cache_read_input_tokens": 0,
        "total_tokens": 0,
        "cost_usd": 0.0,
    }


def _add_record(totals: dict, rec: dict) -> None:
    u = rec["usage"]
    totals["input_tokens"] += u.get("input_tokens", 0)
    totals["output_tokens"] += u.get("output_tokens", 0)
    totals["cache_creation_input_tokens"] += u.get("cache_creation_input_tokens", 0)
    totals["cache_read_input_tokens"] += u.get("cache_read_input_tokens", 0)
    totals["total_tokens"] = (
        totals["input_tokens"]
        + totals["output_tokens"]
        + totals["cache_creation_input_tokens"]
        + totals["cache_read_input_tokens"]
    )
    totals["cost_usd"] += calculate_cost(u, rec.get("model", ""))


def format_tokens(n: int) -> str:
    if n >= 1_000_000:
        return f"{n / 1_000_000:.1f}M"
    if n >= 1_000:
        return f"{n / 1_000:.1f}k"
    return str(n)


def format_cost(usd: float) -> str:
    if usd < 0.01:
        return f"${usd:.4f}"
    if usd < 10:
        return f"${usd:.3f}"
    return f"${usd:.2f}"
