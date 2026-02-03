import '../dtos/plan_dto.dart';
import '../dtos/plan_detail_dto.dart';
import '../models/plan_model.dart';
import '../models/plan_detail_model.dart';

class PlanMapper {
  static PlanModel fromDto(PlanListItemDto dto) {
    return PlanModel(
      id: dto.planId,
      code: dto.planCode,
      title: dto.title,
      status: dto.status,
      totalItems: dto.totalItems,
      totalWeightKg: dto.totalWeightKg,
      volumeUtilizationPct: dto.volumeUtilizationPct,
      createdBy: dto.createdBy,
      createdAt: DateTime.parse(dto.createdAt),
    );
  }

  // Helper to safely parse dates from backend
  static DateTime _parseDate(String dateStr) {
    try {
      return DateTime.parse(dateStr);
    } catch (e) {
      // Fallback to epoch if parsing fails
      return DateTime.fromMillisecondsSinceEpoch(0);
    }
  }

  static DateTime? _parseDateNullable(String? dateStr) {
    if (dateStr == null) return null;
    try {
      return DateTime.parse(dateStr);
    } catch (e) {
      return null;
    }
  }

  // Map PlanDetailDto to PlanDetailModel
  static PlanDetailModel fromDetailDto(PlanDetailDto dto) {
    return PlanDetailModel(
      id: dto.planId,
      code: dto.planCode,
      title: dto.title,
      notes: dto.notes,
      status: dto.status,
      container: _mapContainer(dto.container),
      stats: _mapStats(dto.stats),
      items: dto.items.map(_mapItem).toList(),
      calculation: dto.calculation != null
          ? _mapCalculation(dto.calculation!)
          : null,
      createdBy: dto.createdBy.username ?? 'Unknown',
      createdAt: _parseDate(dto.createdAt),
      updatedAt: _parseDate(dto.updatedAt),
      completedAt: _parseDateNullable(dto.completedAt),
    );
  }

  static ContainerInfo _mapContainer(PlanContainerInfoDto dto) {
    return ContainerInfo(
      containerId: dto.containerId,
      name: dto.name,
      lengthMm: dto.lengthMm,
      widthMm: dto.widthMm,
      heightMm: dto.heightMm,
      maxWeightKg: dto.maxWeightKg,
      volumeM3: dto.volumeM3,
    );
  }

  static PlanStats _mapStats(PlanStatsDto dto) {
    return PlanStats(
      totalItems: dto.totalItems,
      totalWeightKg: dto.totalWeightKg,
      totalVolumeM3: dto.totalVolumeM3,
      volumeUtilizationPct: dto.volumeUtilizationPct,
      weightUtilizationPct: dto.weightUtilizationPct,
    );
  }

  static PlanItem _mapItem(PlanItemDetailDto dto) {
    return PlanItem(
      itemId: dto.itemId,
      productSku: dto.productSku,
      label: dto.label,
      lengthMm: dto.lengthMm,
      widthMm: dto.widthMm,
      heightMm: dto.heightMm,
      weightKg: dto.weightKg,
      quantity: dto.quantity,
      totalWeightKg: dto.totalWeightKg,
      totalVolumeM3: dto.totalVolumeM3,
      allowRotation: dto.allowRotation,
      stackingLimit: dto.stackingLimit,
      colorHex: dto.colorHex,
      createdAt: _parseDate(dto.createdAt),
    );
  }

  static CalculationResult _mapCalculation(CalculationResultDto dto) {
    return CalculationResult(
      jobId: dto.jobId,
      status: dto.status,
      algorithm: dto.algorithm,
      calculatedAt: _parseDateNullable(dto.calculatedAt),
      durationMs: dto.durationMs,
      efficiencyScore: dto.efficiencyScore,
      volumeUtilizationPct: dto.volumeUtilizationPct,
      visualizationUrl: dto.visualizationUrl,
      placements:
          dto.placements?.map(_mapPlacement).toList(),
    );
  }

  static PlacementDetail _mapPlacement(PlacementDetailDto dto) {
    return PlacementDetail(
      placementId: dto.placementId,
      itemId: dto.itemId,
      posX: dto.posX,
      posY: dto.posY,
      posZ: dto.posZ,
      rotation: dto.rotation,
      stepNumber: dto.stepNumber,
    );
  }
}
