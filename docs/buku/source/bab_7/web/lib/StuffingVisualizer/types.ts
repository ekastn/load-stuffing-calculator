/**
 * Type definitions for the StuffingVisualizer
 */

export interface ContainerData {
    name: string;
    length_mm: number;
    width_mm: number;
    height_mm: number;
    max_weight_kg: number;
}

export interface ItemData {
    item_id: string;
    label: string;
    length_mm: number;
    width_mm: number;
    height_mm: number;
    weight_kg: number;
    quantity: number;
    color_hex: string;
}

export interface PlacementData {
    placement_id: string;
    item_id: string;
    pos_x: number;
    pos_y: number;
    pos_z: number;
    rotation: number;
    step_number: number;
}

export interface StuffingPlanData {
    plan_id: string;
    plan_code: string;
    container: ContainerData;
    items: ItemData[];
    placements: PlacementData[];
    stats: {
        total_items: number;
        fitted_count: number;
        unfitted_count: number;
        volume_utilization_pct: number;
    };
}

export interface SceneConfig {
    backgroundColor?: string;
    stepDuration?: number;
}
