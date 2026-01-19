/**
 * Utility functions for coordinate conversion
 */

export function mmToMeters(mm: number): number {
    return mm / 1000;
}

/**
 * Converts placement coordinates from container space to Three.js space
 * 
 * Container coordinate system (from API):
 * - X axis: along container length (left to right)
 * - Y axis: along container width (front to back)
 * - Z axis: along container height (bottom to top)
 * - Origin: front-left-bottom corner
 *
 * Three.js coordinate system:
 * - X axis: horizontal (left to right)
 * - Y axis: vertical (bottom to top)
 * - Z axis: depth (front to back)
 * - Origin: center of container
 */
export function containerToThreeCoords(
    x: number,
    y: number,
    z: number,
    containerLength: number,
    containerWidth: number,
    containerHeight: number
): [number, number, number] {
    const xMeters = mmToMeters(x);
    const yMeters = mmToMeters(y);
    const zMeters = mmToMeters(z);

    const containerLengthMeters = mmToMeters(containerLength);
    const containerWidthMeters = mmToMeters(containerWidth);
    const containerHeightMeters = mmToMeters(containerHeight);

    return [
        xMeters - containerLengthMeters / 2,
        zMeters - containerHeightMeters / 2,
        -(yMeters - containerWidthMeters / 2),
    ];
}

/**
 * Calculate the center offset for an item
 * Items are placed by their corner, but Three.js positions by center
 */
export function getItemCenterOffset(
    itemLength: number,
    itemWidth: number,
    itemHeight: number
): [number, number, number] {
    return [
        mmToMeters(itemLength) / 2,
        mmToMeters(itemHeight) / 2,
        -mmToMeters(itemWidth) / 2
    ];
}
