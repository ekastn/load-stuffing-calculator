import 'package:freezed_annotation/freezed_annotation.dart';

part 'plan_detail_model.freezed.dart';

@freezed
abstract class PlanDetailModel with _$PlanDetailModel {
  const factory PlanDetailModel({
    required String id,
    required String code,
    required String title,
    String? notes,
    required String status,
    required ContainerInfo container,
    required PlanStats stats,
    required List<PlanItem> items,
    CalculationResult? calculation,
    required String createdBy,
    required DateTime createdAt,
    required DateTime updatedAt,
    DateTime? completedAt,
  }) = _PlanDetailModel;
}

@freezed
abstract class ContainerInfo with _$ContainerInfo {
  const factory ContainerInfo({
    String? containerId,
    String? name,
    required double lengthMm,
    required double widthMm,
    required double heightMm,
    required double maxWeightKg,
    required double volumeM3,
  }) = _ContainerInfo;
}

@freezed
abstract class PlanStats with _$PlanStats {
  const factory PlanStats({
    required int totalItems,
    required double totalWeightKg,
    required double totalVolumeM3,
    required double volumeUtilizationPct,
    required double weightUtilizationPct,
  }) = _PlanStats;
}

@freezed
abstract class PlanItem with _$PlanItem {
  const factory PlanItem({
    required String itemId,
    String? productSku,
    String? label,
    required double lengthMm,
    required double widthMm,
    required double heightMm,
    required double weightKg,
    required int quantity,
    required double totalWeightKg,
    required double totalVolumeM3,
    required bool allowRotation,
    required int stackingLimit,
    String? colorHex,
    required DateTime createdAt,
  }) = _PlanItem;
}

@freezed
abstract class CalculationResult with _$CalculationResult {
  const factory CalculationResult({
    required String jobId,
    required String status,
    required String algorithm,
    DateTime? calculatedAt,
    int? durationMs,
    double? efficiencyScore,
    double? volumeUtilizationPct,
    required String visualizationUrl,
    List<PlacementDetail>? placements,
  }) = _CalculationResult;
}

@freezed
abstract class PlacementDetail with _$PlacementDetail {
  const factory PlacementDetail({
    required String placementId,
    required String itemId,
    required double posX,
    required double posY,
    required double posZ,
    required int rotation,
    required int stepNumber,
  }) = _PlacementDetail;
}
