export interface ApiResponse<T> {
    success: boolean;
    data: T;
    error?: {
        code: string;
        message: string;
        details?: any;
    };
}

// --- Container Types ---

export interface Container {
    id: string;
    name: string;
    length_mm: number;
    width_mm: number;
    height_mm: number;
    max_weight_kg: number;
}

export interface CreateContainerRequest {
    name: string;
    length_mm: number;
    width_mm: number;
    height_mm: number;
    max_weight_kg: number;
}

export interface UpdateContainerRequest {
    name: string;
    length_mm: number;
    width_mm: number;
    height_mm: number;
    max_weight_kg: number;
}

// --- Product Types ---

export interface Product {
    id: string;
    label: string;
    sku: string;
    length_mm: number;
    width_mm: number;
    height_mm: number;
    weight_kg: number;
}

export interface CreateProductRequest {
    label: string;
    sku: string;
    length_mm: number;
    width_mm: number;
    height_mm: number;
    weight_kg: number;
}

export interface UpdateProductRequest {
    label: string;
    sku: string;
    length_mm: number;
    width_mm: number;
    height_mm: number;
    weight_kg: number;
}

// --- Plan Types ---

export interface Plan {
    id: string;
    container_id: string;
    container_name?: string;
    status: string; // "draft", "created", "completed", "failed"? Check backend status enum if any
}

export interface PlanItem {
    id: string;
    product_id: string;
    label: string;
    quantity: number;
    length_mm: number;
    width_mm: number;
    height_mm: number;
    weight_kg: number;
}

export interface Placement {
    id: string;
    product_id: string;
    label: string;
    pos_x: number;
    pos_y: number;
    pos_z: number;
    rotation: number;
    step_number: number;
}

export interface PlanDetail extends Plan {
    items: PlanItem[];
    placements?: Placement[];
}

export interface CreatePlanRequest {
    container_id: string;
}

export interface AddPlanItemRequest {
    product_id: string;
    quantity: number;
}

export interface UpdatePlanItemRequest {
    quantity: number;
}
