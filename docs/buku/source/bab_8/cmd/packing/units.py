"""Unit conversion utilities for the packing service."""

from __future__ import annotations

from typing import Literal

# Supported length units
LengthUnit = Literal["mm", "cm", "m"]


def to_cm(value: float, unit: LengthUnit) -> float:
    """Convert a value from the given unit to centimeters."""
    if unit == "mm":
        return value / 10.0
    elif unit == "cm":
        return value
    elif unit == "m":
        return value * 100.0
    else:
        raise ValueError(f"unsupported unit: {unit}")


def from_cm(value_cm: float, unit: LengthUnit) -> float:
    """Convert a value from centimeters to the given unit."""
    if unit == "mm":
        return value_cm * 10.0
    elif unit == "cm":
        return value_cm
    elif unit == "m":
        return value_cm / 100.0
    else:
        raise ValueError(f"unsupported unit: {unit}")


def cm_int(value: float) -> int:
    """Round a float cm value to the nearest integer."""
    return int(round(value))
