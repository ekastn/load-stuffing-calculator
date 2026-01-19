"""Schema definitions and request parsing for the packing service."""

from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Literal, NotRequired, TypedDict

try:
    from .units import LengthUnit
except ImportError:  # pragma: no cover
    from units import LengthUnit


# ============================================================================
# Request Types
# ============================================================================

class ContainerIn(TypedDict):
    """Container dimensions from the request."""
    length: float
    width: float
    height: float
    max_weight: float


class ItemIn(TypedDict):
    """Item to be packed from the request."""
    item_id: str
    label: str
    length: float
    width: float
    height: float
    weight: float
    quantity: int


class OptionsIn(TypedDict, total=False):
    """Packing algorithm options."""
    fix_point: bool
    check_stable: bool
    support_surface_ratio: float
    bigger_first: bool


class RequestIn(TypedDict):
    """Parsed and validated request."""
    units: LengthUnit
    container: ContainerIn
    items: list[ItemIn]
    options: OptionsIn


# ============================================================================
# Internal Types
# ============================================================================

@dataclass(frozen=True)
class NormalizedItem:
    """Item normalized to internal units (integer cm)."""
    item_id: str
    label: str
    quantity: int
    l_cm: int
    w_cm: int
    h_cm: int
    weight_kg: float


# ============================================================================
# Response Types
# ============================================================================

class PlacementOut(TypedDict):
    """Placement result for a single item."""
    item_id: str
    label: str
    pos_x: float
    pos_y: float
    pos_z: float
    rotation: int
    step_number: int


class UnfittedOut(TypedDict):
    """Items that could not fit in the container."""
    item_id: str
    label: str
    count: int


class StatsOut(TypedDict):
    """Packing statistics."""
    expanded_items: int
    fitted_count: int
    unfitted_count: int
    pack_time_ms: int
    total_time_ms: NotRequired[int]


class PackDataOut(TypedDict):
    """Data payload for successful pack response."""
    units: LengthUnit
    placements: list[PlacementOut]
    unfitted: list[UnfittedOut]
    stats: StatsOut


class PackSuccessResponse(TypedDict):
    """Successful pack response."""
    success: Literal[True]
    data: PackDataOut


class ErrorInfo(TypedDict):
    """Error details."""
    code: str
    message: str
    details: dict[str, Any]


class ErrorResponse(TypedDict):
    """Error response."""
    success: Literal[False]
    error: ErrorInfo


class HealthDataOut(TypedDict):
    """Health check data."""
    status: Literal["ok"]


class HealthResponse(TypedDict):
    """Health check response."""
    success: Literal[True]
    data: HealthDataOut


# ============================================================================
# Request Parsing
# ============================================================================

def require_key(obj: dict[str, Any], key: str) -> Any:
    """Get a required key from dict or raise ValueError."""
    if key not in obj:
        raise ValueError(f"missing required field: {key}")
    return obj[key]


def as_float(value: Any, field: str) -> float:
    """Parse a value as a positive float."""
    try:
        n = float(value)
    except Exception as e:
        raise ValueError(f"invalid number for {field}") from e
    if n <= 0:
        raise ValueError(f"{field} must be > 0")
    return n


def as_int(value: Any, field: str) -> int:
    """Parse a value as a positive integer."""
    try:
        n = int(value)
    except Exception as e:
        raise ValueError(f"invalid integer for {field}") from e
    if n <= 0:
        raise ValueError(f"{field} must be > 0")
    return n


def parse_request(body: dict[str, Any]) -> RequestIn:
    """Parse and validate a pack request body."""
    # Validate units
    units = require_key(body, "units")
    if units not in {"mm", "cm", "m"}:
        raise ValueError("units must be one of: mm, cm, m")

    # Validate container
    container_raw = require_key(body, "container")
    if not isinstance(container_raw, dict):
        raise ValueError("container must be an object")

    container: ContainerIn = {
        "length": as_float(require_key(container_raw, "length"), "container.length"),
        "width": as_float(require_key(container_raw, "width"), "container.width"),
        "height": as_float(require_key(container_raw, "height"), "container.height"),
        "max_weight": as_float(require_key(container_raw, "max_weight"), "container.max_weight"),
    }

    # Validate items
    items_raw = require_key(body, "items")
    if not isinstance(items_raw, list) or len(items_raw) == 0:
        raise ValueError("items must be a non-empty array")

    items: list[ItemIn] = []
    for idx, item in enumerate(items_raw):
        if not isinstance(item, dict):
            raise ValueError(f"items[{idx}] must be an object")
        items.append({
            "item_id": str(require_key(item, "item_id")),
            "label": str(require_key(item, "label")),
            "length": as_float(require_key(item, "length"), f"items[{idx}].length"),
            "width": as_float(require_key(item, "width"), f"items[{idx}].width"),
            "height": as_float(require_key(item, "height"), f"items[{idx}].height"),
            "weight": as_float(require_key(item, "weight"), f"items[{idx}].weight"),
            "quantity": as_int(require_key(item, "quantity"), f"items[{idx}].quantity"),
        })

    # Parse options with defaults
    options_raw = body.get("options", {})
    if not isinstance(options_raw, dict):
        raise ValueError("options must be an object")

    options: OptionsIn = {
        "fix_point": bool(options_raw.get("fix_point", True)),
        "check_stable": bool(options_raw.get("check_stable", True)),
        "support_surface_ratio": float(options_raw.get("support_surface_ratio", 0.75)),
        "bigger_first": bool(options_raw.get("bigger_first", True)),
    }

    return {
        "units": units,  # type: ignore
        "container": container,
        "items": items,
        "options": options,
    }
