from __future__ import annotations

from typing import Iterable


def rotation_code(original_lwh: tuple[int, int, int], rotated_lwh: tuple[int, int, int]) -> int:
    l, w, h = original_lwh

    candidates: list[tuple[int, int, int]] = [
        (l, w, h),
        (w, l, h),
        (w, h, l),
        (h, w, l),
        (h, l, w),
        (l, h, w),
    ]

    try:
        return candidates.index(rotated_lwh)
    except ValueError:
        # Equal-sided items can be ambiguous; fall back to a best effort.
        for idx, dims in enumerate(candidates):
            if set(dims) == set(rotated_lwh):
                return idx
        return 0


def permuted_lwh_from_packing_dims(
    *,
    pack_dims: Iterable[int],
) -> tuple[int, int, int]:
    # We pack py3dbp with WHD=(L,H,W). Item.getDimension() returns (x,y,z) in that space => (Lr,Hr,Wr).
    dx, dy, dz = list(pack_dims)
    return (int(dx), int(dz), int(dy))
