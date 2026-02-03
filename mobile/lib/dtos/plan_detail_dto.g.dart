// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'plan_detail_dto.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

PlanDetailDto _$PlanDetailDtoFromJson(Map<String, dynamic> json) =>
    PlanDetailDto(
      planId: json['plan_id'] as String,
      planCode: json['plan_code'] as String,
      title: json['title'] as String,
      notes: json['notes'] as String?,
      status: json['status'] as String,
      container: PlanContainerInfoDto.fromJson(
        json['container'] as Map<String, dynamic>,
      ),
      stats: PlanStatsDto.fromJson(json['stats'] as Map<String, dynamic>),
      items: (json['items'] as List<dynamic>)
          .map((e) => PlanItemDetailDto.fromJson(e as Map<String, dynamic>))
          .toList(),
      calculation: json['calculation'] == null
          ? null
          : CalculationResultDto.fromJson(
              json['calculation'] as Map<String, dynamic>,
            ),
      createdBy: UserSummaryDto.fromJson(
        json['created_by'] as Map<String, dynamic>,
      ),
      createdAt: json['created_at'] as String,
      updatedAt: json['updated_at'] as String,
      completedAt: json['completed_at'] as String?,
    );

Map<String, dynamic> _$PlanDetailDtoToJson(PlanDetailDto instance) =>
    <String, dynamic>{
      'plan_id': instance.planId,
      'plan_code': instance.planCode,
      'title': instance.title,
      'notes': instance.notes,
      'status': instance.status,
      'container': instance.container,
      'stats': instance.stats,
      'items': instance.items,
      'calculation': instance.calculation,
      'created_by': instance.createdBy,
      'created_at': instance.createdAt,
      'updated_at': instance.updatedAt,
      'completed_at': instance.completedAt,
    };

PlanContainerInfoDto _$PlanContainerInfoDtoFromJson(
  Map<String, dynamic> json,
) => PlanContainerInfoDto(
  containerId: json['container_id'] as String?,
  name: json['name'] as String?,
  lengthMm: (json['length_mm'] as num).toDouble(),
  widthMm: (json['width_mm'] as num).toDouble(),
  heightMm: (json['height_mm'] as num).toDouble(),
  maxWeightKg: (json['max_weight_kg'] as num).toDouble(),
  volumeM3: (json['volume_m3'] as num).toDouble(),
);

Map<String, dynamic> _$PlanContainerInfoDtoToJson(
  PlanContainerInfoDto instance,
) => <String, dynamic>{
  'container_id': instance.containerId,
  'name': instance.name,
  'length_mm': instance.lengthMm,
  'width_mm': instance.widthMm,
  'height_mm': instance.heightMm,
  'max_weight_kg': instance.maxWeightKg,
  'volume_m3': instance.volumeM3,
};

PlanStatsDto _$PlanStatsDtoFromJson(Map<String, dynamic> json) => PlanStatsDto(
  totalItems: (json['total_items'] as num).toInt(),
  totalWeightKg: (json['total_weight_kg'] as num).toDouble(),
  totalVolumeM3: (json['total_volume_m3'] as num).toDouble(),
  volumeUtilizationPct: (json['volume_utilization_pct'] as num).toDouble(),
  weightUtilizationPct: (json['weight_utilization_pct'] as num).toDouble(),
);

Map<String, dynamic> _$PlanStatsDtoToJson(PlanStatsDto instance) =>
    <String, dynamic>{
      'total_items': instance.totalItems,
      'total_weight_kg': instance.totalWeightKg,
      'total_volume_m3': instance.totalVolumeM3,
      'volume_utilization_pct': instance.volumeUtilizationPct,
      'weight_utilization_pct': instance.weightUtilizationPct,
    };

PlanItemDetailDto _$PlanItemDetailDtoFromJson(Map<String, dynamic> json) =>
    PlanItemDetailDto(
      itemId: json['item_id'] as String,
      productSku: json['product_sku'] as String?,
      label: json['label'] as String?,
      lengthMm: (json['length_mm'] as num).toDouble(),
      widthMm: (json['width_mm'] as num).toDouble(),
      heightMm: (json['height_mm'] as num).toDouble(),
      weightKg: (json['weight_kg'] as num).toDouble(),
      quantity: (json['quantity'] as num).toInt(),
      totalWeightKg: (json['total_weight_kg'] as num).toDouble(),
      totalVolumeM3: (json['total_volume_m3'] as num).toDouble(),
      allowRotation: json['allow_rotation'] as bool,
      stackingLimit: (json['stacking_limit'] as num).toInt(),
      colorHex: json['color_hex'] as String?,
      createdAt: json['created_at'] as String,
    );

Map<String, dynamic> _$PlanItemDetailDtoToJson(PlanItemDetailDto instance) =>
    <String, dynamic>{
      'item_id': instance.itemId,
      'product_sku': instance.productSku,
      'label': instance.label,
      'length_mm': instance.lengthMm,
      'width_mm': instance.widthMm,
      'height_mm': instance.heightMm,
      'weight_kg': instance.weightKg,
      'quantity': instance.quantity,
      'total_weight_kg': instance.totalWeightKg,
      'total_volume_m3': instance.totalVolumeM3,
      'allow_rotation': instance.allowRotation,
      'stacking_limit': instance.stackingLimit,
      'color_hex': instance.colorHex,
      'created_at': instance.createdAt,
    };

CalculationResultDto _$CalculationResultDtoFromJson(
  Map<String, dynamic> json,
) => CalculationResultDto(
  jobId: json['job_id'] as String,
  status: json['status'] as String,
  algorithm: json['algorithm'] as String,
  calculatedAt: json['calculated_at'] as String?,
  durationMs: (json['duration_ms'] as num?)?.toInt(),
  efficiencyScore: (json['efficiency_score'] as num?)?.toDouble(),
  volumeUtilizationPct: (json['volume_utilization_pct'] as num?)?.toDouble(),
  visualizationUrl: json['visualization_url'] as String,
  placements: (json['placements'] as List<dynamic>?)
      ?.map((e) => PlacementDetailDto.fromJson(e as Map<String, dynamic>))
      .toList(),
);

Map<String, dynamic> _$CalculationResultDtoToJson(
  CalculationResultDto instance,
) => <String, dynamic>{
  'job_id': instance.jobId,
  'status': instance.status,
  'algorithm': instance.algorithm,
  'calculated_at': instance.calculatedAt,
  'duration_ms': instance.durationMs,
  'efficiency_score': instance.efficiencyScore,
  'volume_utilization_pct': instance.volumeUtilizationPct,
  'visualization_url': instance.visualizationUrl,
  'placements': instance.placements,
};

PlacementDetailDto _$PlacementDetailDtoFromJson(Map<String, dynamic> json) =>
    PlacementDetailDto(
      placementId: json['placement_id'] as String,
      itemId: json['item_id'] as String,
      posX: (json['pos_x'] as num).toDouble(),
      posY: (json['pos_y'] as num).toDouble(),
      posZ: (json['pos_z'] as num).toDouble(),
      rotation: (json['rotation'] as num).toInt(),
      stepNumber: (json['step_number'] as num).toInt(),
    );

Map<String, dynamic> _$PlacementDetailDtoToJson(PlacementDetailDto instance) =>
    <String, dynamic>{
      'placement_id': instance.placementId,
      'item_id': instance.itemId,
      'pos_x': instance.posX,
      'pos_y': instance.posY,
      'pos_z': instance.posZ,
      'rotation': instance.rotation,
      'step_number': instance.stepNumber,
    };

UserSummaryDto _$UserSummaryDtoFromJson(Map<String, dynamic> json) =>
    UserSummaryDto(
      userId: json['id'] as String,
      username: json['username'] as String,
      role: json['role'] as String?,
    );

Map<String, dynamic> _$UserSummaryDtoToJson(UserSummaryDto instance) =>
    <String, dynamic>{
      'id': instance.userId,
      'username': instance.username,
      'role': instance.role,
    };
