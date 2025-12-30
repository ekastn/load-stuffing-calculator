from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Literal, NotRequired, TypedDict

try:
    from .units import LengthUnit
except ImportError:  # pragma: no cover
    from units import LengthUnit


class ContainerIn(TypedDict):
    length: float
    width: float
    height: float
    max_weight: float


class ItemIn(TypedDict):
    item_id: str
    label: str
    length: float
    width: float
    height: float
    weight: float
    quantity: int


PutType = Literal[1, 2]


class OptionsIn(TypedDict, total=False):
    fix_point: bool
    check_stable: bool
    support_surface_ratio: float
    bigger_first: bool
    put_type: PutType


class RequestIn(TypedDict):
    units: LengthUnit
    container: ContainerIn
    items: list[ItemIn]
    options: OptionsIn


@dataclass(frozen=True)
class NormalizedItem:
    item_id: str
    label: str
    quantity: int
    l_cm: int
    w_cm: int
    h_cm: int
    weight_kg: float


class PlacementOut(TypedDict):
    item_id: str
    label: str
    pos_x: float
    pos_y: float
    pos_z: float
    rotation: int
    step_number: int


class UnfittedOut(TypedDict):
    item_id: str
    label: str
    count: int


class StatsOut(TypedDict):
    expanded_items: int
    fitted_count: int
    unfitted_count: int
    pack_time_ms: int
    total_time_ms: NotRequired[int]


class PackDataOut(TypedDict):
    units: LengthUnit
    placements: list[PlacementOut]
    unfitted: list[UnfittedOut]
    stats: StatsOut


class PackSuccessResponse(TypedDict):
    success: Literal[True]
    data: PackDataOut


ErrorCode = Literal["INVALID_JSON", "INVALID_REQUEST", "PACKING_FAILED"]


class ErrorInfo(TypedDict):
    code: ErrorCode
    message: str
    details: dict[str, Any]


class ErrorResponse(TypedDict):
    success: Literal[False]
    error: ErrorInfo


PackResponse = PackSuccessResponse | ErrorResponse


class HealthDataOut(TypedDict):
    status: Literal["ok"]


class HealthResponse(TypedDict):
    success: Literal[True]
    data: HealthDataOut


def require_key(obj: dict[str, Any], key: str) -> Any:
    if key not in obj:
        raise ValueError(f"missing required field: {key}")
    return obj[key]


def as_float(value: Any, field: str) -> float:
    try:
        n = float(value)
    except Exception as e:  # noqa: BLE001
        raise ValueError(f"invalid number for {field}") from e
    if n <= 0:
        raise ValueError(f"{field} must be > 0")
    return n


def as_int(value: Any, field: str) -> int:
    try:
        n = int(value)
    except Exception as e:  # noqa: BLE001
        raise ValueError(f"invalid integer for {field}") from e
    if n <= 0:
        raise ValueError(f"{field} must be > 0")
    return n


def parse_request(body: dict[str, Any]) -> RequestIn:
    units = require_key(body, "units")
    if units not in {"mm", "cm", "m"}:
        raise ValueError("units must be one of: mm, cm, m")

    container_raw = require_key(body, "container")
    if not isinstance(container_raw, dict):
        raise ValueError("container must be an object")

    container: ContainerIn = {
        "length": as_float(require_key(container_raw, "length"), "container.length"),
        "width": as_float(require_key(container_raw, "width"), "container.width"),
        "height": as_float(require_key(container_raw, "height"), "container.height"),
        "max_weight": as_float(require_key(container_raw, "max_weight"), "container.max_weight"),
    }

    items_raw = require_key(body, "items")
    if not isinstance(items_raw, list) or len(items_raw) == 0:
        raise ValueError("items must be a non-empty array")

    items: list[ItemIn] = []
    for idx, item in enumerate(items_raw):
        if not isinstance(item, dict):
            raise ValueError(f"items[{idx}] must be an object")
        items.append(
            {
                "item_id": str(require_key(item, "item_id")),
                "label": str(require_key(item, "label")),
                "length": as_float(require_key(item, "length"), f"items[{idx}].length"),
                "width": as_float(require_key(item, "width"), f"items[{idx}].width"),
                "height": as_float(require_key(item, "height"), f"items[{idx}].height"),
                "weight": as_float(require_key(item, "weight"), f"items[{idx}].weight"),
                "quantity": as_int(require_key(item, "quantity"), f"items[{idx}].quantity"),
            }
        )

    options_raw = body.get("options")
    if options_raw is None:
        options_raw = {}
    if not isinstance(options_raw, dict):
        raise ValueError("options must be an object")

    put_type = int(options_raw.get("put_type", 1))
    if put_type not in (1, 2):
        raise ValueError("options.put_type must be 1 (general) or 2 (open top)")

    support_surface_ratio = float(options_raw.get("support_surface_ratio", 0.75))
    if support_surface_ratio <= 0 or support_surface_ratio > 1:
        raise ValueError("options.support_surface_ratio must be in (0, 1]")

    options: OptionsIn = {
        "fix_point": bool(options_raw.get("fix_point", True)),
        "check_stable": bool(options_raw.get("check_stable", True)),
        "support_surface_ratio": support_surface_ratio,
        "bigger_first": bool(options_raw.get("bigger_first", True)),
        "put_type": put_type,  # type: ignore[typeddict-item]
    }

    return {
        "units": units,  # type: ignore[typeddict-item]
        "container": container,
        "items": items,
        "options": options,
    }
