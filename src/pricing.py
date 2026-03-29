# Model pricing table (USD per 1M tokens)
# Update these values as Anthropic adjusts pricing

PRICING: dict[str, dict[str, float]] = {
    "opus": {
        "input_per_m": 15.00,
        "output_per_m": 75.00,
        "cache_creation_per_m": 18.75,  # 1.25x input
        "cache_read_per_m": 1.50,        # 0.1x input
    },
    "sonnet": {
        "input_per_m": 3.00,
        "output_per_m": 15.00,
        "cache_creation_per_m": 3.75,
        "cache_read_per_m": 0.30,
    },
    "haiku": {
        "input_per_m": 0.80,
        "output_per_m": 4.00,
        "cache_creation_per_m": 1.00,
        "cache_read_per_m": 0.08,
    },
}

_DEFAULT_TIER = "sonnet"


def _get_tier(model: str) -> dict[str, float]:
    model_lower = model.lower()
    for tier in ("opus", "sonnet", "haiku"):
        if tier in model_lower:
            return PRICING[tier]
    return PRICING[_DEFAULT_TIER]


def calculate_cost(usage: dict, model: str = "") -> float:
    """Return estimated USD cost for a single usage record."""
    tier = _get_tier(model)
    m = 1_000_000

    cost = (
        usage.get("input_tokens", 0) * tier["input_per_m"] / m
        + usage.get("output_tokens", 0) * tier["output_per_m"] / m
        + usage.get("cache_creation_input_tokens", 0) * tier["cache_creation_per_m"] / m
        + usage.get("cache_read_input_tokens", 0) * tier["cache_read_per_m"] / m
    )
    return cost
