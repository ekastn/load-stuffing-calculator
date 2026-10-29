export interface ContainerData {
    name: string;
    length_mm: number;
    width_mm: number;
    height_mm: number;
    max_weight_kg: number;
    volume_m3: number;
}

export interface ItemData {
    item_id: string;
    label: string;
    sku?: string;
    length_mm: number;
    width_mm: number;
    height_mm: number;
    weight_kg: number;
    quantity: number;
    total_volume_m3: number;
    total_weight_kg: number;
    color_hex: string;
    allow_rotation: boolean;
    stacking_limit: number;
    created_at: string;
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
    stats: {
        total_items: number;
        total_weight_kg: number;
        total_volume_m3: number;
        volume_utilization_pct: number;
        weight_utilization_pct: number;
    };
    calculation: {
        job_id: string;
        status: string;
        algorithm: string;
        placements: PlacementData[];
        volume_utilization_pct: number;
        efficiency_score: number;
        visualization_url: string;
    };
}

export interface SceneConfig {
    backgroundColor?: string;
    cameraPosition?: [number, number, number];
    cameraNear?: number;
    cameraFar?: number;
    stepDuration?: number;
    companyName?: string;
    gridHelper?: boolean;
    axesHelper?: boolean;
}
