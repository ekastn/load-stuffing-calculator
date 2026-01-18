"""Core packing logic using py3dbp algorithm."""

from __future__ import annotations

import os
import sys
import time
from typing import Any

from .rotation import permuted_lwh_from_packing_dims, rotation_code
from .schema import NormalizedItem, PackSuccessResponse, PlacementOut, UnfittedOut, parse_request
from .units import LengthUnit, cm_int, from_cm, to_cm


def _load_py3dbp() -> None:
    """Add the vendored py3dbp library to the Python path."""
    lib_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "py3dbp"))
    if lib_root not in sys.path:
        sys.path.insert(0, lib_root)


def pack_request(body: dict[str, Any]) -> PackSuccessResponse:
    """
    Process a pack request and return placements.
    
    This function orchestrates the entire packing process:
    1. Parse and validate the request
    2. Convert units to internal representation (integer cm)
    3. Initialize py3dbp Packer with bin and items
    4. Run the packing algorithm
    5. Transform results back to API format
    """
    req = parse_request(body)
    units: LengthUnit = req["units"]

    # Convert container dimensions to integer cm
    container_cm = {
        "length": cm_int(to_cm(req["container"]["length"], units)),
        "width": cm_int(to_cm(req["container"]["width"], units)),
        "height": cm_int(to_cm(req["container"]["height"], units)),
    }

    # py3dbp uses WHD (width, height, depth) convention
    # We map our L,W,H to their W,H,D as (L, H, W)
    bin_whd = (container_cm["length"], container_cm["height"], container_cm["width"])

    # Normalize items to internal representation
    normalized_items: list[NormalizedItem] = []
    for item in req["items"]:
        normalized_items.append(
            NormalizedItem(
                item_id=item["item_id"],
                label=item["label"],
                quantity=int(item["quantity"]),
                l_cm=cm_int(to_cm(item["length"], units)),
                w_cm=cm_int(to_cm(item["width"], units)),
                h_cm=cm_int(to_cm(item["height"], units)),
                weight_kg=float(item["weight"]),
            )
        )

    # Load py3dbp library
    _load_py3dbp()
    from py3dbp import Bin, Item, Packer  # type: ignore

    start_pack_at = time.perf_counter()

    # Initialize packer
    packer = Packer()
    packer.addBin(
        Bin(
            partno="bin",
            WHD=bin_whd,
            max_weight=float(req["container"]["max_weight"]),
            corner=0,
            put_type=1,
        )
    )

    # Expand items by quantity and add to packer
    expanded_total = 0
    labels_by_id: dict[str, str] = {it.item_id: it.label for it in normalized_items}

    for it in normalized_items:
        item_whd = (it.l_cm, it.h_cm, it.w_cm)  # Map L,W,H to py3dbp's W,H,D
        for i in range(it.quantity):
            expanded_total += 1
            packer.addItem(
                Item(
                    partno=f"{it.item_id}:{i + 1}",
                    name=it.item_id,
                    typeof="cube",
                    WHD=item_whd,
                    weight=it.weight_kg,
                    level=1,
                    loadbear=100,
                    updown=True,
                    color="#cccccc",
                )
            )

    # Run packing algorithm
    packer.pack(
        bigger_first=bool(req["options"].get("bigger_first", True)),
        distribute_items=False,
        fix_point=bool(req["options"].get("fix_point", True)),
        check_stable=bool(req["options"].get("check_stable", True)),
        support_surface_ratio=float(req["options"].get("support_surface_ratio", 0.75)),
        binding=[],
        number_of_decimals=0,
    )

    pack_time_ms = int((time.perf_counter() - start_pack_at) * 1000)

    # Get results from first (and only) bin
    b = packer.bins[0]

    # Build lookup for rotation mapping
    item_dims_by_id: dict[str, tuple[int, int, int]] = {
        it.item_id: (it.l_cm, it.w_cm, it.h_cm) for it in normalized_items
    }

    # Process fitted items
    placements: list[PlacementOut] = []
    for step_idx, pit in enumerate(b.items):
        item_id = pit.name
        orig_lwh = item_dims_by_id[item_id]
        rotated_lwh = permuted_lwh_from_packing_dims(pack_dims=pit.getDimension())
        rot_code = rotation_code(orig_lwh, rotated_lwh)

        # Position in packing axes: (x=L, y=H, z=W)
        px, py, pz = pit.position

        # Map to API axes: x=L, y=W, z=H
        api_pos_x_cm = int(px)
        api_pos_y_cm = int(pz)
        api_pos_z_cm = int(py)

        placements.append({
            "item_id": item_id,
            "label": labels_by_id.get(item_id, item_id),
            "pos_x": from_cm(api_pos_x_cm, units),
            "pos_y": from_cm(api_pos_y_cm, units),
            "pos_z": from_cm(api_pos_z_cm, units),
            "rotation": rot_code,
            "step_number": step_idx + 1,
        })

    # Process unfitted items
    unfitted_counts: dict[str, int] = {}
    for uit in getattr(b, "unfitted_items", []):
        unfitted_counts[uit.name] = unfitted_counts.get(uit.name, 0) + 1

    unfitted: list[UnfittedOut] = [
        {
            "item_id": item_id,
            "label": labels_by_id.get(item_id, item_id),
            "count": count,
        }
        for item_id, count in sorted(unfitted_counts.items())
    ]

    return {
        "success": True,
        "data": {
            "units": units,
            "placements": placements,
            "unfitted": unfitted,
            "stats": {
                "expanded_items": expanded_total,
                "fitted_count": len(placements),
                "unfitted_count": sum(unfitted_counts.values()),
                "pack_time_ms": pack_time_ms,
            },
        },
    }
