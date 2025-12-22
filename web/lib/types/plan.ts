import { UserSummary } from "./auth"

export interface CreatePlanContainer {
  container_id?: string
  length_mm?: number
  width_mm?: number
  height_mm?: number
  max_weight_kg?: number
}

export interface CreatePlanItem {
  product_sku?: string
  label?: string
  length_mm: number
  width_mm: number
  height_mm: number
  weight_kg: number
  quantity: number
  allow_rotation?: boolean
  color_hex?: string
}

export interface CreatePlanRequest {
  title: string
  notes?: string
  auto_calculate?: boolean
  container: CreatePlanContainer
  items: CreatePlanItem[]
}

export interface PlacementDetail {
  placement_id: string
  item_id: string
  pos_x: number
  pos_y: number
  pos_z: number
  rotation: number
  step_number: number
}

export interface CalculationResult {
  job_id: string
  status: string
  algorithm: string
  calculated_at?: string
  duration_ms: number
  efficiency_score: number
  volume_utilization_pct: number
  visualization_url: string
  placements?: PlacementDetail[]
}

export interface CreatePlanResponse {
  plan_id: string
  plan_code: string
  title: string
  status: string
  total_items: number
  total_weight_kg: number
  total_volume_m3: number
  calculation_job_id?: string
  calculation?: CalculationResult
  created_at: string
}

export interface PlanContainerInfo {
  container_id?: string
  name?: string
  length_mm: number
  width_mm: number
  height_mm: number
  max_weight_kg: number
  volume_m3: number
}

export interface PlanStats {
  total_items: number
  total_weight_kg: number
  total_volume_m3: number
  volume_utilization_pct: number
  weight_utilization_pct: number
}

export interface PlanItemDetail {
  item_id: string
  product_sku?: string
  label?: string
  length_mm: number
  width_mm: number
  height_mm: number
  weight_kg: number
  quantity: number
  total_weight_kg: number
  total_volume_m3: number
  allow_rotation: boolean
  stacking_limit: number
  color_hex?: string
  created_at: string
}

export interface PlanDetailResponse {
  plan_id: string
  plan_code: string
  title: string
  notes?: string
  status: string
  container: PlanContainerInfo
  stats: PlanStats
  items: PlanItemDetail[]
  calculation?: CalculationResult
  created_by: UserSummary
  created_at: string
  updated_at: string
  completed_at?: string
}

export interface PlanListItem {
  plan_id: string
  plan_code: string
  title: string
  status: string
  total_items: number
  total_weight_kg: number
  volume_utilization_pct?: number
  created_by: string
  created_at: string
}

export interface UpdatePlanRequest {
  status?: string
  container?: CreatePlanContainer
}

export interface AddPlanItemRequest extends CreatePlanItem {}

export interface UpdatePlanItemRequest {
  label?: string
  length_mm?: number
  width_mm?: number
  height_mm?: number
  weight_kg?: number
  quantity?: number
  allow_rotation?: boolean
  color_hex?: string
}

export interface CalculatePlanRequest {
  strategy?: string
  goal?: string
  gravity?: boolean
}
