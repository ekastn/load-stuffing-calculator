"""Rotation code utilities for mapping py3dbp rotation to visualization."""

from __future__ import annotations


def rotation_code(orig_lwh: tuple[int, int, int], rotated_lwh: tuple[int, int, int]) -> int:
    """
    Determine the rotation code (0-5) based on how dimensions changed.
    
    The rotation code indicates which rotation was applied:
    - 0: No rotation (L,W,H)
    - 1: Rotate 90° around Z axis (W,L,H)
    - 2: Rotate 90° around Y axis (H,W,L)
    - 3: Rotate 90° around X axis (L,H,W)
    - 4: Rotate 180° around Z, then 90° Y (H,L,W)
    - 5: Rotate 90° Z, then 90° Y (W,H,L)
    """
    l, w, h = orig_lwh
    rl, rw, rh = rotated_lwh
    
    # Map all possible permutations to rotation codes
    rotations = {
        (l, w, h): 0,  # Original
        (w, l, h): 1,  # Rotated around Z
        (h, w, l): 2,  # Rotated around Y
        (l, h, w): 3,  # Rotated around X
        (h, l, w): 4,  # Combined rotation
        (w, h, l): 5,  # Combined rotation
    }
    
    return rotations.get((rl, rw, rh), 0)


def permuted_lwh_from_packing_dims(pack_dims: tuple[int, int, int]) -> tuple[int, int, int]:
    """
    Convert py3dbp packing dimensions back to our L,W,H convention.
    
    py3dbp uses WHD (width, height, depth) internally.
    We pack with WHD = (L, H, W), so we need to convert back.
    """
    # pack_dims is in py3dbp's WHD order after rotation
    # Convert back to our L,W,H
    return (pack_dims[0], pack_dims[2], pack_dims[1])
