from __future__ import annotations

import math
from typing import Literal

LengthUnit = Literal["mm", "cm", "m"]


def to_cm(value: float, unit: LengthUnit) -> float:
    if unit == "mm":
        return value / 10.0
    if unit == "cm":
        return value
    if unit == "m":
        return value * 100.0
    raise ValueError(f"unsupported units: {unit}")


def from_cm(value_cm: float, unit: LengthUnit) -> float:
    if unit == "mm":
        return value_cm * 10.0
    if unit == "cm":
        return value_cm
    if unit == "m":
        return value_cm / 100.0
    raise ValueError(f"unsupported units: {unit}")


def cm_int(value: float) -> int:
    # Normalize to integer centimeters for py3dbp performance.
    # Round (user preference) and clamp to >= 1.
    return max(1, int(round(value)))


def ceil_cm_int(value: float) -> int:
    return max(1, int(math.ceil(value)))
