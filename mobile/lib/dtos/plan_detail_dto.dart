import 'package:json_annotation/json_annotation.dart';

part 'plan_detail_dto.g.dart';

// Main plan detail response
@JsonSerializable()
class PlanDetailDto {
  @JsonKey(name: 'plan_id')
  final String planId;
  @JsonKey(name: 'plan_code')
  final String planCode;
  final String title;
  final String? notes;
  final String status;
  final PlanContainerInfoDto container;
  final PlanStatsDto stats;
  final List<PlanItemDetailDto> items;
  final CalculationResultDto? calculation;
  @JsonKey(name: 'created_by')
  final UserSummaryDto createdBy;
  @JsonKey(name: 'created_at')
  final String createdAt;
  @JsonKey(name: 'updated_at')
  final String updatedAt;
  @JsonKey(name: 'completed_at')
  final String? completedAt;

  PlanDetailDto({
    required this.planId,
    required this.planCode,
    required this.title,
    this.notes,
    required this.status,
    required this.container,
    required this.stats,
    required this.items,
    this.calculation,
    required this.createdBy,
    required this.createdAt,
    required this.updatedAt,
    this.completedAt,
  });

  factory PlanDetailDto.fromJson(Map<String, dynamic> json) =>
      _$PlanDetailDtoFromJson(json);
  Map<String, dynamic> toJson() => _$PlanDetailDtoToJson(this);
}

// Container info
@JsonSerializable()
class PlanContainerInfoDto {
  @JsonKey(name: 'container_id')
  final String? containerId;
  final String? name;
  @JsonKey(name: 'length_mm')
  final double lengthMm;
  @JsonKey(name: 'width_mm')
  final double widthMm;
  @JsonKey(name: 'height_mm')
  final double heightMm;
  @JsonKey(name: 'max_weight_kg')
  final double maxWeightKg;
  @JsonKey(name: 'volume_m3')
  final double volumeM3;

  PlanContainerInfoDto({
    this.containerId,
    this.name,
    required this.lengthMm,
    required this.widthMm,
    required this.heightMm,
    required this.maxWeightKg,
    required this.volumeM3,
  });

  factory PlanContainerInfoDto.fromJson(Map<String, dynamic> json) =>
      _$PlanContainerInfoDtoFromJson(json);
  Map<String, dynamic> toJson() => _$PlanContainerInfoDtoToJson(this);
}

// Plan stats
@JsonSerializable()
class PlanStatsDto {
  @JsonKey(name: 'total_items')
  final int totalItems;
  @JsonKey(name: 'total_weight_kg')
  final double totalWeightKg;
  @JsonKey(name: 'total_volume_m3')
  final double totalVolumeM3;
  @JsonKey(name: 'volume_utilization_pct')
  final double volumeUtilizationPct;
  @JsonKey(name: 'weight_utilization_pct')
  final double weightUtilizationPct;

  PlanStatsDto({
    required this.totalItems,
    required this.totalWeightKg,
    required this.totalVolumeM3,
    required this.volumeUtilizationPct,
    required this.weightUtilizationPct,
  });

  factory PlanStatsDto.fromJson(Map<String, dynamic> json) =>
      _$PlanStatsDtoFromJson(json);
  Map<String, dynamic> toJson() => _$PlanStatsDtoToJson(this);
}

// Plan item detail  
@JsonSerializable()
class PlanItemDetailDto {
  @JsonKey(name: 'item_id')
  final String itemId;
  @JsonKey(name: 'product_sku')
  final String? productSku;
  final String? label;
  @JsonKey(name: 'length_mm')
  final double lengthMm;
  @JsonKey(name: 'width_mm')
  final double widthMm;
  @JsonKey(name: 'height_mm')
  final double heightMm;
  @JsonKey(name: 'weight_kg')
  final double weightKg;
  final int quantity;
  @JsonKey(name: 'total_weight_kg')
  final double totalWeightKg;
  @JsonKey(name: 'total_volume_m3')
  final double totalVolumeM3;
  @JsonKey(name: 'allow_rotation')
  final bool allowRotation;
  @JsonKey(name: 'stacking_limit')
  final int stackingLimit;
  @JsonKey(name: 'color_hex')
  final String? colorHex;
  @JsonKey(name: 'created_at')
  final String createdAt;

  PlanItemDetailDto({
    required this.itemId,
    this.productSku,
    this.label,
    required this.lengthMm,
    required this.widthMm,
    required this.heightMm,
    required this.weightKg,
    required this.quantity,
    required this.totalWeightKg,
    required this.totalVolumeM3,
    required this.allowRotation,
    required this.stackingLimit,
    this.colorHex,
    required this.createdAt,
  });

  factory PlanItemDetailDto.fromJson(Map<String, dynamic> json) =>
      _$PlanItemDetailDtoFromJson(json);
  Map<String, dynamic> toJson() => _$PlanItemDetailDtoToJson(this);
}

// Calculation result
@JsonSerializable()
class CalculationResultDto {
  @JsonKey(name: 'job_id')
  final String jobId;
  final String status; // queued | running | completed | failed
  final String algorithm;
  @JsonKey(name: 'calculated_at')
  final String? calculatedAt;
  @JsonKey(name: 'duration_ms')
  final int? durationMs;
  @JsonKey(name: 'efficiency_score')
  final double? efficiencyScore;
  @JsonKey(name: 'volume_utilization_pct')
  final double? volumeUtilizationPct;
  @JsonKey(name: 'visualization_url')
  final String visualizationUrl;
  final List<PlacementDetailDto>? placements;

  CalculationResultDto({
    required this.jobId,
    required this.status,
    required this.algorithm,
    this.calculatedAt,
    this.durationMs,
    this.efficiencyScore,
    this.volumeUtilizationPct,
    required this.visualizationUrl,
    this.placements,
  });

  factory CalculationResultDto.fromJson(Map<String, dynamic> json) =>
      _$CalculationResultDtoFromJson(json);
  Map<String, dynamic> toJson() => _$CalculationResultDtoToJson(this);
}

// Placement detail
@JsonSerializable()
class PlacementDetailDto {
  @JsonKey(name: 'placement_id')
  final String placementId;
  @JsonKey(name: 'item_id')
  final String itemId;
  @JsonKey(name: 'pos_x')
  final double posX;
  @JsonKey(name: 'pos_y')
  final double posY;
  @JsonKey( name: 'pos_z')
  final double posZ;
  final int rotation;
  @JsonKey(name: 'step_number')
  final int stepNumber;

  PlacementDetailDto({
    required this.placementId,
    required this.itemId,
    required this.posX,
    required this.posY,
    required this.posZ,
    required this.rotation,
    required this.stepNumber,
  });

  factory PlacementDetailDto.fromJson(Map<String, dynamic> json) =>
      _$PlacementDetailDtoFromJson(json);
  Map<String, dynamic> toJson() => _$PlacementDetailDtoToJson(this);
}

// User summary
@JsonSerializable()
class UserSummaryDto {
  @JsonKey(name: 'id')
  final String userId;
  final String username;
  final String? role;

  UserSummaryDto({
    required this.userId,
    required this.username,
    this.role,
  });

  factory UserSummaryDto.fromJson(Map<String, dynamic> json) =>
      _$UserSummaryDtoFromJson(json);
  Map<String, dynamic> toJson() => _$UserSummaryDtoToJson(this);
}
