import type { PlanDetail, Container } from "./api/types";
import type { StuffingPlanData } from "./StuffingVisualizer/types";

export function transformToVisualizerData(plan: PlanDetail, container: Container): StuffingPlanData {
    // Generate random colors for items since backend doesn't store them yet
    const getColor = (id: string) => {
        const colors = ["#3b82f6", "#22c55e", "#f59e0b", "#ef4444", "#8b5cf6", "#ec4899", "#06b6d4"];
        let hash = 0;
        for (let i = 0; i < id.length; i++) {
            hash = id.charCodeAt(i) + ((hash << 5) - hash);
        }
        return colors[Math.abs(hash) % colors.length];
    };

    return {
        plan_id: plan.id,
        plan_code: plan.id.substring(0, 8).toUpperCase(),
        container: {
            name: container.name,
            length_mm: container.length_mm,
            width_mm: container.width_mm,
            height_mm: container.height_mm,
            max_weight_kg: container.max_weight_kg,
        },
        items: plan.items.map((item) => ({
            item_id: item.product_id, // Group by product type
            label: item.label,
            length_mm: item.length_mm,
            width_mm: item.width_mm,
            height_mm: item.height_mm,
            weight_kg: item.weight_kg,
            quantity: item.quantity,
            color_hex: getColor(item.product_id),
        })),
        placements: (plan.placements || []).map((p) => ({
            placement_id: p.id,
            item_id: p.product_id,
            pos_x: p.pos_x,
            pos_y: p.pos_y,
            pos_z: p.pos_z,
            rotation: p.rotation,
            step_number: p.step_number,
        })),
        stats: {
            total_items: plan.items.reduce((sum, item) => sum + item.quantity, 0),
            fitted_count: plan.placements ? plan.placements.length : 0,
            unfitted_count: 0, // Should be calculated if backend provides unfitted list
            volume_utilization_pct: 0, // Could calculate if needed
        },
    };
}
