import type { StuffingPlanData } from "@/lib/StuffingVisualizer"
import type { PlanDetailResponse } from "@/lib/types"

export function toStuffingPlanData(plan: PlanDetailResponse): StuffingPlanData {
  return {
    plan_id: plan.plan_id,
    plan_code: plan.plan_code,
    container: {
      name: plan.container.name || "Container",
      length_mm: plan.container.length_mm,
      width_mm: plan.container.width_mm,
      height_mm: plan.container.height_mm,
      max_weight_kg: plan.container.max_weight_kg,
      volume_m3: plan.container.volume_m3,
    },
    items: plan.items.map((item) => ({
      item_id: item.item_id,
      label: item.label || "Item",
      length_mm: item.length_mm,
      width_mm: item.width_mm,
      height_mm: item.height_mm,
      weight_kg: item.weight_kg,
      quantity: item.quantity,
      total_volume_m3: item.total_volume_m3,
      total_weight_kg: item.total_weight_kg,
      color_hex: item.color_hex || "#3498db",
      allow_rotation: item.allow_rotation,
      stacking_limit: item.stacking_limit,
      created_at: item.created_at,
    })),
    stats: {
      total_items: plan.stats.total_items,
      total_weight_kg: plan.stats.total_weight_kg,
      total_volume_m3: plan.stats.total_volume_m3,
      volume_utilization_pct: plan.stats.volume_utilization_pct,
      weight_utilization_pct: plan.stats.weight_utilization_pct,
    },
    calculation: {
      job_id: plan.calculation?.job_id || "",
      status: plan.calculation?.status || "",
      algorithm: plan.calculation?.algorithm || "",
      placements: plan.calculation?.placements || [],
      volume_utilization_pct: plan.calculation?.volume_utilization_pct || 0,
      efficiency_score: plan.calculation?.efficiency_score || 0,
      visualization_url: plan.calculation?.visualization_url || "",
    },
  }
}
