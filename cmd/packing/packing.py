from __future__ import annotations

import os
import sys
import time
from typing import Any

try:
    from .rotation import permuted_lwh_from_packing_dims, rotation_code
    from .schema import NormalizedItem, PackSuccessResponse, PlacementOut, UnfittedOut, parse_request
    from .units import LengthUnit, cm_int, from_cm, to_cm
except ImportError:  # pragma: no cover
    # Allow running as a script.
    from rotation import permuted_lwh_from_packing_dims, rotation_code
    from schema import NormalizedItem, PackSuccessResponse, PlacementOut, UnfittedOut, parse_request
    from units import LengthUnit, cm_int, from_cm, to_cm


def _load_py3dbp() -> None:
    # `py3dbp` is vendored under this service directory.
    lib_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "3D-bin-packing"))
    if lib_root not in sys.path:
        sys.path.insert(0, lib_root)



def pack_request(body: dict[str, Any]) -> PackSuccessResponse:
    req = parse_request(body)
    units: LengthUnit = req["units"]

    # Convert request-length units -> integer cm.
    container_cm = {
        "length": cm_int(to_cm(req["container"]["length"], units)),
        "width": cm_int(to_cm(req["container"]["width"], units)),
        "height": cm_int(to_cm(req["container"]["height"], units)),
    }

    # Pack axes (physics correct): py3dbp WHD maps to (x,y,z). We want y=height.
    # We pack WHD=(L,H,W).
    bin_whd = (container_cm["length"], container_cm["height"], container_cm["width"])

    normalized_items: list[NormalizedItem] = []
    for item in req["items"]:
        length_cm = cm_int(to_cm(item["length"], units))
        width_cm = cm_int(to_cm(item["width"], units))
        height_cm = cm_int(to_cm(item["height"], units))
        normalized_items.append(
            NormalizedItem(
                item_id=item["item_id"],
                label=item["label"],
                quantity=int(item["quantity"]),
                l_cm=length_cm,
                w_cm=width_cm,
                h_cm=height_cm,
                weight_kg=float(item["weight"]),
            )
        )

    _load_py3dbp()
    from py3dbp import Bin, Item, Packer  # type: ignore

    start_pack_at = time.perf_counter()

    packer = Packer()
    packer.addBin(
        Bin(
            partno="bin",
            WHD=bin_whd,
            max_weight=float(req["container"]["max_weight"]),
            corner=0,
            put_type=int(req["options"].get("put_type", 1)),
        )
    )

    # Expand quantities into individual items. We keep a reverse mapping so we can
    # re-aggregate unfitted counts by item_id.
    expanded_total = 0
    labels_by_id: dict[str, str] = {it.item_id: it.label for it in normalized_items}

    for it in normalized_items:
        # Feed py3dbp WHD in our packing axes: (L,H,W)
        item_whd = (it.l_cm, it.h_cm, it.w_cm)
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

    b = packer.bins[0]

    # Build lookup for original lwh in cm for rotation mapping.
    item_dims_by_id: dict[str, tuple[int, int, int]] = {
        it.item_id: (it.l_cm, it.w_cm, it.h_cm) for it in normalized_items
    }

    # Fitted items are already ordered by py3dbp's putOrder() at end of pack().
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

        placements.append(
            {
                "item_id": item_id,
                "label": labels_by_id.get(item_id, item_id),
                "pos_x": from_cm(api_pos_x_cm, units),
                "pos_y": from_cm(api_pos_y_cm, units),
                "pos_z": from_cm(api_pos_z_cm, units),
                "rotation": rot_code,
                "step_number": step_idx + 1,
            }
        )

    # Unfitted items: aggregate counts by item_id.
    unfitted_counts: dict[str, int] = {}
    for uit in getattr(b, "unfitted_items", []):
        unfitted_counts[uit.name] = unfitted_counts.get(uit.name, 0) + 1

    unfitted: list[UnfittedOut] = [
        {
            "item_id": item_id,
            "label": labels_by_id.get(item_id, item_id),
            "count": count,
        }
        for item_id, count in sorted(unfitted_counts.items(), key=lambda kv: kv[0])
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
